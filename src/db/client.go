package db

import (
	"context"

	"cfg"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func GetDynamoClient(ctx context.Context) *dynamodb.Client {
	// Configure AWS SDK client
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Config.AwsRegion), // Replace with your region
	)
	if err != nil {
		panic("Failed to load aws sdk config, " + err.Error())
	}

	client := dynamodb.NewFromConfig(cfg)

	return client
}
