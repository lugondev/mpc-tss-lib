package grpc_client

import (
	"errors"
	"github.com/lugondev/mpc-tss-lib/pb"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

func CreateClientConnections(clients []string) []*grpc.ClientConn {
	connections := funk.Map(clients, func(clientHost string) *grpc.ClientConn {
		conn, err := grpc.Dial(clientHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Error().Err(err).Msgf("failed to connect: %s", clientHost)
			return nil
		}
		return conn
	})

	return funk.Filter(connections, func(conn *grpc.ClientConn) bool {
		return conn != nil
	}).([]*grpc.ClientConn)
}

func GetPartyIDs(parties []*pb.GetPartiesResponse) []string {
	return funk.Map(parties, func(party *pb.GetPartiesResponse) string {
		return party.Id
	}).([]string)
}

func CallPartiesGRPCs[K comparable](clients []string, parties []string, callToClient func(pb.MpcPartyClient, []string, int) (K, error)) ([]K, error) {
	connectionGRPCs := CreateClientConnections(clients)
	defer func() {
		for _, c := range connectionGRPCs {
			_ = c.Close()
		}
	}()
	if len(connectionGRPCs) != len(parties) {
		return nil, errors.New("failed to connect to all clients")
	}

	var wg sync.WaitGroup
	var results []K
	for i, conn := range connectionGRPCs {
		wg.Add(1)
		go func(conn *grpc.ClientConn, index int) {
			client := pb.NewMpcPartyClient(conn)
			r, err := callToClient(client, parties, index)
			if err == nil {
				logger.Debug().Msgf("result: %v", r)
				results = append(results, r)
			} else {
				logger.Error().Err(err).Msgf("failed to call party %d", index)
			}
			wg.Done()
		}(conn, i)
	}
	wg.Wait()
	return results, nil
}

func CallClientGRPCs[K comparable](clients []string, callToClient func(pb.MpcPartyClient, int) (K, error)) ([]K, error) {
	connectionGRPCs := CreateClientConnections(clients)
	defer func() {
		for _, c := range connectionGRPCs {
			_ = c.Close()
		}
	}()
	if len(connectionGRPCs) != len(clients) {
		return nil, errors.New("failed to connect to all clients")
	}

	var wg sync.WaitGroup
	results := make([]K, len(clients))
	for i, conn := range connectionGRPCs {
		wg.Add(1)
		go func(conn *grpc.ClientConn, index int) {
			client := pb.NewMpcPartyClient(conn)
			r, err := callToClient(client, index)
			if err == nil {
				logger.Debug().Msgf("result: %v", r)
				results[index] = r
			} else {
				logger.Error().Err(err).Msgf("failed to call client %d", clients[index])
			}
			wg.Done()
		}(conn, i)
	}
	wg.Wait()
	return results, nil
}
