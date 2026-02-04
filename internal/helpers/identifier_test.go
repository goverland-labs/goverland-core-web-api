package helpers

import (
	"context"
	"testing"

	"github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockEnsClient struct {
	getAddressesByEnsNamesFunc func(ctx context.Context, in *storagepb.AddressesByEnsNamesRequest, opts ...grpc.CallOption) (*storagepb.AddressesByEnsNamesResponse, error)
}

func (m *mockEnsClient) GetEnsByAddresses(_ context.Context, _ *storagepb.EnsByAddressesRequest, _ ...grpc.CallOption) (*storagepb.EnsByAddressesResponse, error) {
	return nil, nil
}

func (m *mockEnsClient) GetAddressesByEnsNames(ctx context.Context, in *storagepb.AddressesByEnsNamesRequest, opts ...grpc.CallOption) (*storagepb.AddressesByEnsNamesResponse, error) {
	if m.getAddressesByEnsNamesFunc != nil {
		return m.getAddressesByEnsNamesFunc(ctx, in, opts...)
	}
	return &storagepb.AddressesByEnsNamesResponse{}, nil
}

func TestIdentifierResolver_Resolve_HexAddress(t *testing.T) {
	resolver := NewIdentifierResolver(&mockEnsClient{})

	tests := []struct {
		name    string
		input   string
		wantAddr string
		wantENS bool
	}{
		{"lowercase hex", "0x329c54289ff5d6b7b7dae13592c6b1eda1543ed4", "0x329c54289ff5d6b7b7dae13592c6b1eda1543ed4", false},
		{"checksum hex", "0x329c54289Ff5D6B7b7daE13592C6B1EDA1543eD4", "0x329c54289Ff5D6B7b7daE13592C6B1EDA1543eD4", false},
		{"uppercase prefix", "0X329c54289ff5d6b7b7dae13592c6b1eda1543ed4", "0X329c54289ff5d6b7b7dae13592c6b1eda1543ed4", false},
		{"with whitespace", "  0x329c54289ff5d6b7b7dae13592c6b1eda1543ed4  ", "0x329c54289ff5d6b7b7dae13592c6b1eda1543ed4", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolver.Resolve(context.Background(), tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Address != tt.wantAddr {
				t.Errorf("address = %q, want %q", result.Address, tt.wantAddr)
			}
			if result.WasENS != tt.wantENS {
				t.Errorf("wasENS = %v, want %v", result.WasENS, tt.wantENS)
			}
			if result.ENSName != "" {
				t.Errorf("ensName = %q, want empty", result.ENSName)
			}
		})
	}
}

func TestIdentifierResolver_Resolve_ENSName(t *testing.T) {
	mock := &mockEnsClient{
		getAddressesByEnsNamesFunc: func(_ context.Context, in *storagepb.AddressesByEnsNamesRequest, _ ...grpc.CallOption) (*storagepb.AddressesByEnsNamesResponse, error) {
			if len(in.GetNames()) > 0 && in.GetNames()[0] == "aci.eth" {
				return &storagepb.AddressesByEnsNamesResponse{
					EnsNames: []*storagepb.EnsName{
						{Address: "0x329c54289Ff5D6B7b7daE13592C6B1EDA1543eD4", Name: "aci.eth"},
					},
				}, nil
			}
			return &storagepb.AddressesByEnsNamesResponse{}, nil
		},
	}

	resolver := NewIdentifierResolver(mock)

	result, err := resolver.Resolve(context.Background(), "aci.eth")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Address != "0x329c54289Ff5D6B7b7daE13592C6B1EDA1543eD4" {
		t.Errorf("address = %q, want %q", result.Address, "0x329c54289Ff5D6B7b7daE13592C6B1EDA1543eD4")
	}
	if result.ENSName != "aci.eth" {
		t.Errorf("ensName = %q, want %q", result.ENSName, "aci.eth")
	}
	if !result.WasENS {
		t.Error("wasENS = false, want true")
	}
}

func TestIdentifierResolver_Resolve_ENSNameLowercase(t *testing.T) {
	var receivedName string
	mock := &mockEnsClient{
		getAddressesByEnsNamesFunc: func(_ context.Context, in *storagepb.AddressesByEnsNamesRequest, _ ...grpc.CallOption) (*storagepb.AddressesByEnsNamesResponse, error) {
			if len(in.GetNames()) > 0 {
				receivedName = in.GetNames()[0]
			}
			return &storagepb.AddressesByEnsNamesResponse{
				EnsNames: []*storagepb.EnsName{
					{Address: "0x329c54289Ff5D6B7b7daE13592C6B1EDA1543eD4", Name: "aci.eth"},
				},
			}, nil
		},
	}

	resolver := NewIdentifierResolver(mock)

	_, err := resolver.Resolve(context.Background(), "ACI.ETH")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedName != "aci.eth" {
		t.Errorf("name sent to ENS client = %q, want %q", receivedName, "aci.eth")
	}
}

func TestIdentifierResolver_Resolve_UnknownENSName(t *testing.T) {
	mock := &mockEnsClient{
		getAddressesByEnsNamesFunc: func(_ context.Context, _ *storagepb.AddressesByEnsNamesRequest, _ ...grpc.CallOption) (*storagepb.AddressesByEnsNamesResponse, error) {
			return &storagepb.AddressesByEnsNamesResponse{}, nil
		},
	}

	resolver := NewIdentifierResolver(mock)

	_, err := resolver.Resolve(context.Background(), "unknown.eth")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got %T: %v", err, err)
	}
	if st.Code() != codes.NotFound {
		t.Errorf("error code = %v, want %v", st.Code(), codes.NotFound)
	}
}

func TestIdentifierResolver_Resolve_EmptyAddressResult(t *testing.T) {
	mock := &mockEnsClient{
		getAddressesByEnsNamesFunc: func(_ context.Context, _ *storagepb.AddressesByEnsNamesRequest, _ ...grpc.CallOption) (*storagepb.AddressesByEnsNamesResponse, error) {
			return &storagepb.AddressesByEnsNamesResponse{
				EnsNames: []*storagepb.EnsName{
					{Address: "", Name: "stale.eth"},
				},
			}, nil
		},
	}

	resolver := NewIdentifierResolver(mock)

	_, err := resolver.Resolve(context.Background(), "stale.eth")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got %T: %v", err, err)
	}
	if st.Code() != codes.NotFound {
		t.Errorf("error code = %v, want %v", st.Code(), codes.NotFound)
	}
}

func TestIdentifierResolver_Resolve_GRPCError(t *testing.T) {
	mock := &mockEnsClient{
		getAddressesByEnsNamesFunc: func(_ context.Context, _ *storagepb.AddressesByEnsNamesRequest, _ ...grpc.CallOption) (*storagepb.AddressesByEnsNamesResponse, error) {
			return nil, status.Error(codes.Internal, "db down")
		},
	}

	resolver := NewIdentifierResolver(mock)

	_, err := resolver.Resolve(context.Background(), "aci.eth")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got %T: %v", err, err)
	}
	if st.Code() != codes.Internal {
		t.Errorf("error code = %v, want %v", st.Code(), codes.Internal)
	}
}

func TestIsHexAddress(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"0x329c54289ff5d6b7b7dae13592c6b1eda1543ed4", true},
		{"0X329c54289ff5d6b7b7dae13592c6b1eda1543ed4", true},
		{"0x329C54289Ff5D6B7b7daE13592C6B1EDA1543eD4", true},
		{"0x0000000000000000000000000000000000000000", true},
		{"aci.eth", false},
		{"0x", false},
		{"0xZZZc54289ff5d6b7b7dae13592c6b1eda1543ed4", false},
		{"329c54289ff5d6b7b7dae13592c6b1eda1543ed4", false},
		{"0x329c54289ff5d6b7b7dae13592c6b1eda1543ed", false},  // 41 chars
		{"0x329c54289ff5d6b7b7dae13592c6b1eda1543ed44", false}, // 43 chars
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := isHexAddress(tt.input)
			if got != tt.want {
				t.Errorf("isHexAddress(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
