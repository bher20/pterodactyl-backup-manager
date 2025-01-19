package azure

import (
	"bufio"
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

const (
	DefaultBlobSuffix = ".blob.core.windows.net"
)

type AzureAuthentication struct {
	TenantId     string `json:"tenantId"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type AzureStorageAccountDestination struct {
	Name                string              `json:"name"`
	ContainerName       string              `json:"containerName"`
	BlobPrefix          *string             `json:"blobPrefix,omitempty"`
	BlobUrlSuffix       *string             `json:"blobUrlSuffix,omitempty"`
	AzureAuthentication AzureAuthentication `json:"azureAuthentication"`
}

func getClient(azureStorageAccount *AzureStorageAccountDestination) (*azblob.Client, error) {
	url := fmt.Sprintf("https://%s", azureStorageAccount.Name)
	blobSuffix := DefaultBlobSuffix

	if azureStorageAccount.BlobUrlSuffix != nil {
		blobSuffix = *azureStorageAccount.BlobUrlSuffix
	}
	url = fmt.Sprintf("%s%s/", url, blobSuffix)

	credential, err := azidentity.NewClientSecretCredential(azureStorageAccount.AzureAuthentication.TenantId, azureStorageAccount.AzureAuthentication.ClientId, azureStorageAccount.AzureAuthentication.ClientSecret, nil)
	if err != nil {
		return nil, err
	}

	client, err := azblob.NewClient(url, credential, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Backup(azureStorageAccount *AzureStorageAccountDestination, fileName string, filePath string) error {
	blobClient, err := getClient(azureStorageAccount)
	if err != nil {
		return err
	}

	// Upload the file to the specified container with the specified blob name
	blobName := fileName
	if azureStorageAccount.BlobPrefix != nil {
		blobName = fmt.Sprintf("%s/%s", *azureStorageAccount.BlobPrefix, blobName)
	}
	log.Info(fmt.Sprintf("Uploading a blob '%s:/%s' to Azure Storage Account '%s'", azureStorageAccount.ContainerName, blobName, azureStorageAccount.Name))

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	_, err = blobClient.UploadStream(context.TODO(), azureStorageAccount.ContainerName, blobName, bufio.NewReader(file), nil)

	return err
}
