package operation

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/go-to-k/delstack/client"
)

var _ IOperator = (*BackupVaultOperator)(nil)

type BackupVaultOperator struct {
	client    *client.Backup
	resources []types.StackResourceSummary
}

func NewBackupVaultOperator(config aws.Config) *BackupVaultOperator {
	client := client.NewBackup(config)
	return &BackupVaultOperator{
		client:    client,
		resources: []types.StackResourceSummary{},
	}
}

func (operator *BackupVaultOperator) AddResources(resource types.StackResourceSummary) {
	operator.resources = append(operator.resources, resource)
}

func (operator *BackupVaultOperator) GetResourcesLength() int {
	return len(operator.resources)
}

func (operator *BackupVaultOperator) DeleteResources() error {
	// TODO: Concurrency Delete
	for _, backupVault := range operator.resources {
		err := operator.DeleteBackupVault(backupVault.PhysicalResourceId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (operator *BackupVaultOperator) DeleteBackupVault(backupVaultName *string) error {
	recoveryPoints, err := operator.client.ListRecoveryPointsByBackupVault(backupVaultName)
	if err != nil {
		return err
	}

	if err := operator.client.DeleteRecoveryPoints(backupVaultName, recoveryPoints); err != nil {
		return err
	}

	if err := operator.client.DeleteBackupVault(backupVaultName); err != nil {
		return err
	}

	return nil
}
