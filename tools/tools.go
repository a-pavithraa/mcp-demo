package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func ListDynamoDbTables() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list-dynamodb-tables",
			mcp.WithDescription("List all DynamoDB tables"),
		), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			//cfg, err := config.LoadDefaultConfig(ctx)
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
			if err != nil {
				return nil, fmt.Errorf("failed to load AWS configuration: %v", err)
			}

			dynamoClient := dynamodb.NewFromConfig(cfg)
			output, err := dynamoClient.ListTables(ctx, &dynamodb.ListTablesInput{})
			if err != nil {
				return nil, fmt.Errorf("failed to list DynamoDB tables: %v", err)
			}

			jsonResult, err := json.Marshal(output.TableNames)
			if err != nil {
				return nil, fmt.Errorf("error marshalling result to JSON: %v", err)
			}

			return mcp.NewToolResultText(string(jsonResult)), nil
		}
}

func GetDynamoDbTableMetadata() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get-dynamodb-table-metadata",
			mcp.WithDescription("Get metadata of DynamoDB tables like created date, size, and pricing model"),
		), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			//	cfg, err := config.LoadDefaultConfig(ctx)
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
			if err != nil {
				return nil, fmt.Errorf("failed to load AWS configuration: %v", err)
			}

			dynamoClient := dynamodb.NewFromConfig(cfg)
			output, err := dynamoClient.ListTables(ctx, &dynamodb.ListTablesInput{})
			if err != nil {
				return nil, fmt.Errorf("failed to list DynamoDB tables: %v", err)
			}

			metadata := make(map[string]interface{})
			for _, tableName := range output.TableNames {
				descOutput, err := dynamoClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{
					TableName: &tableName,
				})
				if err != nil {
					return nil, fmt.Errorf("failed to describe table %s: %v", tableName, err)
				}

				metadata[tableName] = map[string]interface{}{
					"CreatedDate":  descOutput.Table.CreationDateTime,
					"SizeBytes":    descOutput.Table.TableSizeBytes,
					"PricingModel": descOutput.Table.BillingModeSummary,
				}
			}

			jsonResult, err := json.Marshal(metadata)
			if err != nil {
				return nil, fmt.Errorf("error marshalling result to JSON: %v", err)
			}

			return mcp.NewToolResultText(string(jsonResult)), nil
		}
}

func ListKmsKeysWithMetadata() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list-kms-keys",
			mcp.WithDescription("List all KMS keys along with metadata like created date, key description, and ID"),
		), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
			if err != nil {
				return nil, fmt.Errorf("failed to load AWS configuration: %v", err)
			}

			kmsClient := kms.NewFromConfig(cfg)

			// List all keys
			output, err := kmsClient.ListKeys(ctx, &kms.ListKeysInput{})
			if err != nil {
				return nil, fmt.Errorf("failed to list KMS keys: %v", err)
			}

			keysMetadata := make([]map[string]interface{}, 0, len(output.Keys))

			// Fetch metadata for each key
			for _, keyInfo := range output.Keys {
				keyId := *keyInfo.KeyId

				// Get key details
				descOutput, err := kmsClient.DescribeKey(ctx, &kms.DescribeKeyInput{
					KeyId: &keyId,
				})
				if err != nil {
					return nil, fmt.Errorf("failed to describe key %s: %v", keyId, err)
				}

				// Extract key metadata
				keyMetadata := map[string]interface{}{
					"KeyId":        keyId,
					"Arn":          *descOutput.KeyMetadata.Arn,
					"CreationDate": descOutput.KeyMetadata.CreationDate,
					"Description":  descOutput.KeyMetadata.Description,
					"Enabled":      descOutput.KeyMetadata.Enabled,
					"KeyState":     descOutput.KeyMetadata.KeyState,
					"KeyManager":   descOutput.KeyMetadata.KeyManager,
					"KeyUsage":     descOutput.KeyMetadata.KeyUsage,
				}

				keysMetadata = append(keysMetadata, keyMetadata)
			}

			jsonResult, err := json.Marshal(keysMetadata)
			if err != nil {
				return nil, fmt.Errorf("error marshalling result to JSON: %v", err)
			}

			return mcp.NewToolResultText(string(jsonResult)), nil
		}
}

func ListS3BucketsWithMetadata() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list-s3-buckets",
			mcp.WithDescription("List all S3 buckets along with commonly used metadata like creation date, region, versioning status, etc."),
		), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
			if err != nil {
				return nil, fmt.Errorf("failed to load AWS configuration: %v", err)
			}

			s3Client := s3.NewFromConfig(cfg)

			// List all S3 buckets
			listOutput, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
			if err != nil {
				return nil, fmt.Errorf("failed to list S3 buckets: %v", err)
			}

			bucketsMetadata := make([]map[string]interface{}, 0, len(listOutput.Buckets))

			// Fetch metadata for each bucket
			for _, bucket := range listOutput.Buckets {
				bucketName := *bucket.Name

				metadata := map[string]interface{}{
					"Name":         bucketName,
					"CreationDate": bucket.CreationDate,
				}

				// Get bucket location (region)
				locationOutput, err := s3Client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
					Bucket: &bucketName,
				})
				if err == nil {
					metadata["Region"] = locationOutput.LocationConstraint
				}

				// Get bucket versioning status
				versioningOutput, err := s3Client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
					Bucket: &bucketName,
				})
				if err == nil {
					metadata["Versioning"] = versioningOutput.Status
				}

				// Get bucket encryption settings
				encryptionOutput, err := s3Client.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{
					Bucket: &bucketName,
				})
				if err == nil && encryptionOutput.ServerSideEncryptionConfiguration != nil {
					metadata["Encryption"] = encryptionOutput.ServerSideEncryptionConfiguration.Rules
				}

				// Get public access block settings
				publicAccessOutput, err := s3Client.GetPublicAccessBlock(ctx, &s3.GetPublicAccessBlockInput{
					Bucket: &bucketName,
				})
				if err == nil && publicAccessOutput.PublicAccessBlockConfiguration != nil {
					metadata["PublicAccessBlock"] = map[string]interface{}{
						"BlockPublicAcls":       publicAccessOutput.PublicAccessBlockConfiguration.BlockPublicAcls,
						"BlockPublicPolicy":     publicAccessOutput.PublicAccessBlockConfiguration.BlockPublicPolicy,
						"IgnorePublicAcls":      publicAccessOutput.PublicAccessBlockConfiguration.IgnorePublicAcls,
						"RestrictPublicBuckets": publicAccessOutput.PublicAccessBlockConfiguration.RestrictPublicBuckets,
					}
				}

				// Get bucket tagging
				taggingOutput, err := s3Client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
					Bucket: &bucketName,
				})
				if err == nil {
					tags := make(map[string]string)
					for _, tag := range taggingOutput.TagSet {
						tags[*tag.Key] = *tag.Value
					}
					metadata["Tags"] = tags
				}

				bucketsMetadata = append(bucketsMetadata, metadata)
			}

			jsonResult, err := json.Marshal(bucketsMetadata)
			if err != nil {
				return nil, fmt.Errorf("error marshalling result to JSON: %v", err)
			}

			return mcp.NewToolResultText(string(jsonResult)), nil
		}
}
