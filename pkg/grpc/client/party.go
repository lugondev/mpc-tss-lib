package grpc_client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/lugondev/mpc-tss-lib/pb"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"strings"
	"sync"
)

func CreateClientConnections(clients []string) []*grpc.ClientConn {
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		logger.Error().Err(err).Msg("failed to load system roots")
		return []*grpc.ClientConn{}
	}
	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(cred))
	connections := funk.Map(clients, func(clientHost string) *grpc.ClientConn {
		grpcOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		if !strings.Contains(clientHost, "localhost") &&
			!strings.Contains(clientHost, "0.0.0.0") &&
			!strings.Contains(clientHost, "127.0.0.1") {
			grpcOpts = opts
		}
		conn, err := grpc.Dial(clientHost, grpcOpts...)
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

func CreateClientConnection(clientHost string) *grpc.ClientConn {
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		logger.Error().Err(err).Msg("failed to load system roots")
		return nil
	}
	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(cred))
	grpcOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if !strings.Contains(clientHost, "localhost") &&
		!strings.Contains(clientHost, "0.0.0.0") &&
		!strings.Contains(clientHost, "127.0.0.1") {
		grpcOpts = opts
	}
	conn, err := grpc.Dial(clientHost, grpcOpts...)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to connect: %s", clientHost)
		return nil
	}
	return conn
}

func GetPartyIDs(parties []*pb.GetPartyResponse) []string {
	return funk.Map(parties, func(party *pb.GetPartyResponse) string {
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
				logger.Debug().Msgf("party: %v", r)
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
				logger.Error().Err(err).Msgf("failed to call client %s", clients[index])
			}
			wg.Done()
		}(conn, i)
	}
	wg.Wait()
	return results, nil
}

func CallClientGRPC[K comparable](clientHost string, callToClient func(pb.MpcPartyClient) (K, error)) (K, error) {
	connection := CreateClientConnection(clientHost)
	if connection == nil {
		return *new(K), errors.New("failed to connect to client")
	}
	defer func() {
		_ = connection.Close()
	}()

	client := pb.NewMpcPartyClient(connection)
	return callToClient(client)
}
