package main

import (
	"auth"
	"cfg"
	"context"
	"db"
	"encoding/json"
	"fmt"
	"net/http"
	"src"
	"utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ListApiKeysDeps src.Deps

func main() {
	d := ListApiKeysDeps{
		DbClient:  db.GetDynamoClient(context.Background()),
		TableName: cfg.Config.ApiKeyTable,
	}
	lambda.Start(d.handler)
}

func (d *ListApiKeysDeps) handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Verify the request authentication
	workspaceId, err := auth.VerifyAuthHeader(request, d.DbClient)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusUnauthorized, fmt.Sprintf("Unauthorized: %s", err.Error())), nil
	}

	// Parse API ID and key ID from request parameters
	apiId := request.QueryStringParameters["apiId"]

	if apiId == "" {
		return utils.HttpErrorResponse(http.StatusBadRequest, "Missing required query parameters: apiId and apiKeyId"), nil
	}

	// Query the database for all API keys for the workspace
	keys, err := db.ListApiKeys(workspaceId, apiId, d.DbClient)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error listing API keys: %v", err)), nil
	}

	// Construct the response
	keysJson, err := json.Marshal(keys)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error marshalling keys: %v", err)), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(keysJson),
	}, nil
}
