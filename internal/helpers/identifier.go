package helpers

import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResolvedIdentifier struct {
	Address string
	ENSName string
	WasENS  bool
}

type IdentifierResolver struct {
	ensClient storagepb.EnsClient
}

func NewIdentifierResolver(ensClient storagepb.EnsClient) *IdentifierResolver {
	return &IdentifierResolver{
		ensClient: ensClient,
	}
}

func (r *IdentifierResolver) Resolve(ctx context.Context, identifier string) (*ResolvedIdentifier, error) {
	identifier = strings.TrimSpace(identifier)

	if isHexAddress(identifier) {
		return &ResolvedIdentifier{
			Address: identifier,
			WasENS:  false,
		}, nil
	}

	resp, err := r.ensClient.GetAddressesByEnsNames(ctx, &storagepb.AddressesByEnsNamesRequest{
		Names: []string{strings.ToLower(identifier)},
	})
	if err != nil {
		return nil, err
	}

	if len(resp.GetEnsNames()) == 0 || resp.GetEnsNames()[0].GetAddress() == "" {
		return nil, status.Error(codes.NotFound, "identifier not found")
	}

	return &ResolvedIdentifier{
		Address: resp.GetEnsNames()[0].GetAddress(),
		ENSName: identifier,
		WasENS:  true,
	}, nil
}

func isHexAddress(s string) bool {
	if len(s) != 42 {
		return false
	}
	if !strings.HasPrefix(s, "0x") && !strings.HasPrefix(s, "0X") {
		return false
	}
	_, err := hex.DecodeString(s[2:])
	return err == nil
}
