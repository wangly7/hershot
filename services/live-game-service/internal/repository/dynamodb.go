package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBRepository(
	ctx context.Context,
	endpoint string,
	region string,
	tableName string,
) (*DynamoDBRepository, error) {
	region = strings.TrimSpace(region)
	tableName = strings.TrimSpace(tableName)
	endpoint = strings.TrimSpace(endpoint)

	if region == "" {
		return nil, errors.New("AWS region cannot be empty")
	}

	if tableName == "" {
		return nil, errors.New("DynamoDB game events cannot be empty")
	}

	awsConfig, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"load AWS config: %w",
			err,
		)
	}

	client := dynamodb.NewFromConfig(
		awsConfig,
		func(o *dynamodb.Options) {
			if endpoint != "" {
				o.BaseEndpoint = aws.String(endpoint)
			}
		},
	)

	return &DynamoDBRepository{
		client:    client,
		tableName: tableName,
	}, nil
}

func (r *DynamoDBRepository) EnsureTable(
	ctx context.Context,
) error {
	describeTable := &dynamodb.DescribeTableInput{
		TableName: aws.String(r.tableName),
	}

	_, err := r.client.DescribeTable(
		ctx,
		describeTable,
	)

	if err == nil {
		return nil
	}

	var notFound *types.ResourceNotFoundException

	if !errors.As(err, &notFound) {
		return fmt.Errorf(
			"describe DynamoDB table %q: %w",
			r.tableName,
			err,
		)
	}

	_, err = r.client.CreateTable(
		ctx,
		&dynamodb.CreateTableInput{
			TableName: aws.String(r.tableName),

			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("gameId"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("sequence"),
					AttributeType: types.ScalarAttributeTypeN,
				},
			},

			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("gameId"),
					KeyType:       types.KeyTypeHash,
				},
				{
					AttributeName: aws.String("sequence"),
					KeyType:       types.KeyTypeRange,
				},
			},

			BillingMode: types.BillingModePayPerRequest,
		},
	)

	if err != nil {
		var alreadyCreating *types.ResourceInUseException

		if !errors.As(err, &alreadyCreating) {
			return fmt.Errorf(
				"create DynamoDB table %q: %w",
				r.tableName,
				err,
			)
		}
	}

	waiter := dynamodb.NewTableExistsWaiter(r.client)

	err = waiter.Wait(
		ctx,
		describeTable,
		30*time.Second,
	)

	if err != nil {
		return fmt.Errorf(
			"wait for DynamoDB table %q: %w",
			r.tableName,
			err,
		)
	}

	return nil
}
