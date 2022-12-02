package grpc_client_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/lugondev/mpc-tss-lib/pb"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestGRPC(t *testing.T) {
	//clientHost := "mpc-client-02-jhxvtoeu7q-as.a.run.app:443"
	clientHost := "0.0.0.0:8082"
	systemRoots, err := x509.SystemCertPool()
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
	// dial server
	conn, err := grpc.Dial(clientHost, grpcOpts...)
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// create stream
	client := pb.NewMpcPartyClient(conn)
	pongReply, err := client.RequestParty(ctx, &pb.EmptyParams{})
	if err != nil {
		log.Panic(err)
		return
	}
	log.Printf("Pong reply: %s", pongReply.Id)
}
