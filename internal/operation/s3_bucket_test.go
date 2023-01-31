package operation

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfnTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/go-to-k/delstack/internal/io"
	"github.com/go-to-k/delstack/pkg/client"
)

/*
	Test Cases
*/

func TestS3BucketOperator_DeleteS3Bucket(t *testing.T) {
	io.NewLogger(false)
	mock := client.NewMockS3()
	allErrorMock := client.NewAllErrorMockS3()
	deleteBucketErrorMock := client.NewDeleteBucketErrorMockS3()
	deleteObjectsErrorMock := client.NewDeleteObjectsErrorMockS3()
	deleteObjectsErrorAfterZeroLengthMock := client.NewDeleteObjectsErrorAfterZeroLengthMockS3()
	deleteObjectsOutputErrorMock := client.NewDeleteObjectsOutputErrorMockS3()
	deleteObjectsOutputErrorAfterZeroLengthMock := client.NewDeleteObjectsOutputErrorAfterZeroLengthMockS3()
	listObjectVersionsErrorMock := client.NewListObjectVersionsErrorMockS3()
	checkBucketExistsErrorMock := client.NewCheckBucketExistsErrorMockS3()
	checkBucketNotExistsMock := client.NewCheckBucketNotExistsMockS3()

	type args struct {
		ctx        context.Context
		bucketName *string
		client     client.IS3
	}

	cases := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "delete bucket successfully",
			args: args{
				ctx:        context.Background(),
				bucketName: aws.String("test"),
				client:     mock,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "delete bucket failure for all errors",
			args: args{
				ctx:        context.Background(),
				bucketName: aws.String("test"),
				client:     allErrorMock,
			},
			want:    fmt.Errorf("ListBucketsError"),
			wantErr: true,
		},
		{
			name: "delete bucket failure for check bucket exists errors",
			args: args{
				ctx:        context.Background(),
				bucketName: aws.String("test"),
				client:     checkBucketExistsErrorMock,
			},
			want:    fmt.Errorf("ListBucketsError"),
			wantErr: true,
		},
		{
			name: "delete bucket successfully for bucket not exists",
			args: args{
				ctx:        context.Background(),
				bucketName: aws.String("test"),
				client:     checkBucketNotExistsMock,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "delete bucket failure for list object versions errors",
			args: args{
				ctx:        context.Background(),
				bucketName: aws.String("test"),
				client:     listObjectVersionsErrorMock,
			},
			want:    fmt.Errorf("ListObjectVersionsError"),
			wantErr: true,
		},
		{
			name: "delete bucket failure for delete objects errors",
			args: args{
				ctx:        context.Background(),
				bucketName: aws.String("test"),
				client:     deleteObjectsErrorMock,
			},
			want:    fmt.Errorf("DeleteObjectsError"),
			wantErr: true,
		},
		{
			name: "delete bucket successfully for delete objects errors after zero length",
			args: args{
				ctx:        context.Background(),
				bucketName: aws.String("test"),
				client:     deleteObjectsErrorAfterZeroLengthMock,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "delete bucket failure for delete objects output errors",
			args: args{
				ctx:        context.Background(),
				bucketName: aws.String("test"),
				client:     deleteObjectsOutputErrorMock,
			},
			want:    fmt.Errorf("DeleteObjectsError: followings \nCode: Code\nKey: Key\nVersionId: VersionId\nMessage: Message\n"),
			wantErr: true,
		},
		{
			name: "delete bucket successfully for delete objects output errors after zero length",
			args: args{
				ctx:        context.Background(),
				bucketName: aws.String("test"),
				client:     deleteObjectsOutputErrorAfterZeroLengthMock,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "delete bucket failure for delete bucket errors",
			args: args{
				ctx:        context.Background(),
				bucketName: aws.String("test"),
				client:     deleteBucketErrorMock,
			},
			want:    fmt.Errorf("DeleteBucketError"),
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s3BucketOperator := NewS3BucketOperator(tt.args.client)

			err := s3BucketOperator.DeleteS3Bucket(tt.args.ctx, tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %#v, wantErr %#v", err.Error(), tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.want.Error() {
				t.Errorf("err = %#v, want %#v", err.Error(), tt.want.Error())
				return
			}
		})
	}
}

func TestS3BucketOperator_DeleteResourcesForS3Bucket(t *testing.T) {
	io.NewLogger(false)
	mock := client.NewMockS3()
	allErrorMock := client.NewAllErrorMockS3()

	type args struct {
		ctx    context.Context
		client client.IS3
	}

	cases := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "delete resources successfully",
			args: args{
				ctx:    context.Background(),
				client: mock,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "delete resources failure",
			args: args{
				ctx:    context.Background(),
				client: allErrorMock,
			},
			want:    fmt.Errorf("ListBucketsError"),
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s3BucketOperator := NewS3BucketOperator(tt.args.client)
			s3BucketOperator.AddResource(&cfnTypes.StackResourceSummary{
				LogicalResourceId:  aws.String("LogicalResourceId1"),
				ResourceStatus:     "DELETE_FAILED",
				ResourceType:       aws.String("AWS::S3::Bucket"),
				PhysicalResourceId: aws.String("PhysicalResourceId1"),
			})

			err := s3BucketOperator.DeleteResources(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %#v, wantErr %#v", err.Error(), tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.want.Error() {
				t.Errorf("err = %#v, want %#v", err.Error(), tt.want.Error())
				return
			}
		})
	}
}