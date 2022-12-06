package grpc_client

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/lugondev/mpc-tss-lib/db"
	sqlc "github.com/lugondev/mpc-tss-lib/db/sqlc"
	"github.com/lugondev/mpc-tss-lib/pb"
	"github.com/lugondev/mpc-tss-lib/pkg/mpc/networking/client"
	"github.com/lugondev/mpc-tss-lib/pkg/mpc/tss_wrap"
	zerolog "github.com/rs/zerolog/log"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc"
	"net"
	"time"
)

var logger = zerolog.Logger

const ContextTimeout = 2 * time.Minute

type grpcClient struct {
	pb.MpcPartyServer

	hostUrl  string
	sqlStore *db.SQLStore
}

func (s *grpcClient) RequestParty(_ context.Context, _ *pb.EmptyParams) (*pb.RequestPartyResponse, error) {
	logger.Info().Msg("RequestParty called")
	return &pb.RequestPartyResponse{Id: uuid.New().String()}, nil
}

func (s *grpcClient) Ping(_ context.Context, _ *pb.EmptyParams) (*pb.Pong, error) {
	logger.Info().Msg("Ping called")
	return &pb.Pong{Message: fmt.Sprintf("Pong: %s", time.Now().String())}, nil
}

func (s *grpcClient) KeygenGenerator(_ context.Context, keygenRequest *pb.KeygenGeneratorParams) (*pb.KeygenGeneratorResponse, error) {
	logger.Info().Msg("KeygenGenerator called")
	pubkey, err := keygen(s.hostUrl, "", keygenRequest.Id, keygenRequest.Ids, len(keygenRequest.Ids), s.sqlStore)
	if err != nil {
		return nil, err
	}
	return &pb.KeygenGeneratorResponse{
		Pubkey: pubkey,
	}, nil
}

func (s *grpcClient) GetParty(ctx context.Context, getPartiesParams *pb.GetPartyParams) (*pb.GetPartyResponse, error) {
	logger.Info().Msg("GetParty called")
	share, err := s.sqlStore.GetPartyIdByPubkey(ctx, getPartiesParams.Pubkey)
	if err != nil {
		logger.Error().Err(err).Msgf("can't get party id by pubkey: %s", getPartiesParams.Pubkey)
		return nil, err
	}
	return &pb.GetPartyResponse{
		Id:     share.PartyID,
		Pubkey: share.Pubkey,
	}, nil
}

func (s *grpcClient) GetParties(ctx context.Context, _ *pb.GetPartiesParams) (*pb.GetPartiesResponse, error) {
	logger.Info().Msg("GetParties called")
	shares, err := s.sqlStore.ListShare(ctx)
	if err != nil {
		logger.Error().Err(err).Msgf("can't get parties")
		return nil, err
	}

	return &pb.GetPartiesResponse{
		Shares: funk.Map(shares, func(share sqlc.ListShareRow) *pb.PartyShare {
			return &pb.PartyShare{
				Pubkey:  share.Pubkey,
				Address: share.Address,
				PartyId: share.PartyID,
			}
		}).([]*pb.PartyShare),
	}, nil
}

func (s *grpcClient) Sign(_ context.Context, signParams *pb.SignParams) (*pb.SignResponse, error) {
	logger.Info().Msg("Sign called")
	signature, err := sign(s.hostUrl, "", signParams, s.sqlStore)
	if err != nil {
		return nil, err
	}

	return &pb.SignResponse{
		Signature: common.Bytes2Hex(signature),
		Message:   common.Bytes2Hex(signParams.Message),
		Id:        signParams.Id,
	}, nil
}

func StartGrpcClient(port int64, hostUrl string, sqlStore *db.SQLStore) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to listen")
	}
	s := grpc.NewServer()
	pb.RegisterMpcPartyServer(s, &grpcClient{
		hostUrl:  hostUrl,
		sqlStore: sqlStore,
	})

	logger.Info().Msgf("server listening: %s", lis.Addr())
	if err := s.Serve(lis); err != nil {
		logger.Fatal().Err(err).Msg("failed to serve")
	}
}

func keygen(hostURL, accessToken string, id string, partyIDs []string, threshold int, dbStore *db.SQLStore) (string, error) {
	mpcc := tss_wrap.NewMpc(id, threshold, &logger)
	netOperation := client.NewClient(mpcc, hostURL, accessToken, &logger)
	ctx, _ := context.WithTimeout(context.Background(), ContextTimeout)
	if err := netOperation.Keygen(ctx, partyIDs); err != nil {
		return "", err
	}
	err := saveShare(id, mpcc, dbStore)
	if err != nil {
		return "", err
	}

	if publicKey, err := mpcc.GetPublicKey(); err != nil {
		return "", err
	} else {
		return common.Bytes2Hex(crypto.CompressPubkey(publicKey)), nil
	}
}

func saveShare(id string, tss *tss_wrap.Mpc, dbStore *db.SQLStore) error {
	share, err := tss.Share()
	if err != nil {
		logger.Error().Err(err).Msgf("can't save share: %s", id)
		return err
	}
	sharePubkey, err := tss.GetPublicKey()
	addressFromPubKey, _ := tss.GetAddress()
	_, err = dbStore.CreateShare(context.Background(), sqlc.CreateShareParams{
		Pubkey:  common.Bytes2Hex(crypto.CompressPubkey(sharePubkey)),
		Data:    share,
		Address: addressFromPubKey.String(),
		PartyID: id,
	})
	if err != nil {
		logger.Error().Err(err).Msg("can't save share")
		return err
	}
	return nil
}

func sign(hostURL, accessToken string, signParam *pb.SignParams, dbStore *db.SQLStore) ([]byte, error) {
	share, err := getShareData(signParam.Id, dbStore)
	if err != nil {
		logger.Error().Err(err).Msgf("can't get share: %s", signParam.Id)
		return nil, err
	}
	mpcc, err := tss_wrap.NewMpcWithShare(signParam.Id, len(signParam.Parties), share, &logger)
	if err != nil {
		logger.Error().Err(err).Msg("can't create mpc")
		return nil, err
	}
	publicKey, _ := mpcc.GetPublicKey()
	fmt.Println("address: ", crypto.PubkeyToAddress(*publicKey).Hex())

	netOperation := client.NewClient(mpcc, hostURL, accessToken, &logger)
	ctx, _ := context.WithTimeout(context.Background(), ContextTimeout)
	fmt.Println("msg:", common.Bytes2Hex(signParam.Message))

	return netOperation.Sign(ctx, signParam.Parties, signParam.Message)
}

func getShareData(id string, dbStore *db.SQLStore) ([]byte, error) {
	if share, err := dbStore.GetShareByID(context.Background(), id); err != nil {
		return nil, err
	} else {
		return share.Data, nil
	}
}
