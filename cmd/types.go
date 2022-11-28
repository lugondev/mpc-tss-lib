package cmd

import "context"

type NetworkingOperations interface {
	Reshare(ctx context.Context, partyIDsOld, partyIDsNew []string, thresholdNew int) error
	Keygen(ctx context.Context, party []string) error
	Sign(ctx context.Context, party []string, msg []byte) ([]byte, error)
}
