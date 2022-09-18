package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/go-to-k/delstack/logger"
)

var _ IEcrSDKClient = (*mockEcrSDKClient)(nil)
var _ IEcrSDKClient = (*errorMockEcrSDKClient)(nil)

/*
	Mocks for SDK Client
*/
type mockEcrSDKClient struct{}

func NewMockEcrSDKClient() *mockEcrSDKClient {
	return &mockEcrSDKClient{}
}

func (m *mockEcrSDKClient) DeleteRepository(ctx context.Context, params *ecr.DeleteRepositoryInput, optFns ...func(*ecr.Options)) (*ecr.DeleteRepositoryOutput, error) {
	return nil, nil
}

type errorMockEcrSDKClient struct{}

func NewErrorMockEcrSDKClient() *errorMockEcrSDKClient {
	return &errorMockEcrSDKClient{}
}

func (m *errorMockEcrSDKClient) DeleteRepository(ctx context.Context, params *ecr.DeleteRepositoryInput, optFns ...func(*ecr.Options)) (*ecr.DeleteRepositoryOutput, error) {
	return nil, fmt.Errorf("DeleteRepositoryError")
}

/*
	Test Cases
*/

func TestDeleteRepository(t *testing.T) {
	logger.NewLogger()
	ctx := context.TODO()
	mock := NewMockEcrSDKClient()
	errorMock := NewErrorMockEcrSDKClient()

	type args struct {
		ctx            context.Context
		repositoryName *string
		client         IEcrSDKClient
	}

	cases := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "delete repository successfully",
			args: args{
				ctx:            ctx,
				repositoryName: aws.String("test"),
				client:         mock,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "delete repository failure",
			args: args{
				ctx:            ctx,
				repositoryName: aws.String("test"),
				client:         errorMock,
			},
			want:    fmt.Errorf("DeleteRepositoryError"),
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ecrClient := NewEcr(tt.args.client)

			err := ecrClient.DeleteRepository(tt.args.repositoryName)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.want.Error() {
				t.Errorf("err = %#v, want %#v", err, tt.want)
			}
		})
	}
}