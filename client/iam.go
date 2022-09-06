package client

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type IAM struct {
	client *iam.Client
}

func NewIAM(config aws.Config) *IAM {
	client := iam.NewFromConfig(config)
	return &IAM{
		client,
	}
}

func (iamClient *IAM) DeleteRole(roleName *string) error {
	input := &iam.DeleteRoleInput{
		RoleName: roleName,
	}

	_, err := iamClient.client.DeleteRole(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed delete the IAM Role, %v", err)
		return err
	}

	return nil
}

func (iamClient *IAM) ListAttachedRolePolicies(roleName *string) ([]types.AttachedPolicy, error) {
	var marker *string
	policies := []types.AttachedPolicy{}

	for {
		input := &iam.ListAttachedRolePoliciesInput{
			RoleName: roleName,
			Marker:   marker,
		}

		output, err := iamClient.client.ListAttachedRolePolicies(context.TODO(), input)
		if err != nil {
			log.Fatalf("failed list attached Role Policies, %v", err)
			return nil, err
		}

		policies = append(policies, output.AttachedPolicies...)

		marker = output.Marker
		if marker == nil {
			break
		}
	}

	return policies, nil
}

func (iamClient *IAM) DetachRolePolicies(roleName *string, policies []types.AttachedPolicy) error {
	for _, policy := range policies {
		if err := iamClient.DetachRolePolicy(roleName, policy.PolicyArn); err != nil {
			return err
		}
	}

	return nil
}

func (iamClient *IAM) DetachRolePolicy(roleName *string, PolicyArn *string) error {
	input := &iam.DetachRolePolicyInput{
		PolicyArn: PolicyArn,
		RoleName:  roleName,
	}

	_, err := iamClient.client.DetachRolePolicy(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed detach the Role Policy, %v", err)
		return err
	}

	return nil
}