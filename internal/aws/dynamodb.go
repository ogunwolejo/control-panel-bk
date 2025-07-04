package aws

import (
	"context"
	"control-panel-bk/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var DynamoClient *dynamodb.Client

func ConnectDynamoDbs() {
	client := dynamodb.NewFromConfig(*config.AwsConfig)
	DynamoClient = client
}

func AddDoc(items map[string]types.AttributeValue, tableName string, err chan error, output chan *dynamodb.PutItemOutput) {
	defer close(err)
	defer close(output)

	ctx := context.TODO()
	itm := &dynamodb.PutItemInput{
		TableName:    aws.String(tableName),
		Item:         items,
		ReturnValues: types.ReturnValueAllNew,
	}

	opt, e := DynamoClient.PutItem(ctx, itm)

	if e != nil {
		err <- e
		return
	}

	output <- opt
}
