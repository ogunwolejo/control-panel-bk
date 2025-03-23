package aws

import (
	"context"
	"control-panel-bk/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"os"
)

// awsApplication
// arn:aws:resource-groups:us-east-1:528757792684:group/controlPanel/04a7tbp0n21ntgacp58yscw9e0

type Group struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func getClient(cfg *aws.Config) *cognitoidentityprovider.Client {
	return cognitoidentityprovider.NewFromConfig(*cfg)
}

func CreateUserPoolGroup(cfg *aws.Config, group Group) (*cognitoidentityprovider.CreateGroupOutput, error) {
	client := getClient(cfg)
	input := cognitoidentityprovider.CreateGroupInput{
		UserPoolId:  aws.String(os.Getenv("us-east-1_kNKCRvql2")),
		Description: aws.String(group.Description),
		GroupName:   aws.String(group.Name),
	}

	output, err := client.CreateGroup(context.TODO(), &input)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func AddUsersToUserPoolGroup(cfg *aws.Config, groupName string, username string) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
	client := getClient(cfg)

	input := cognitoidentityprovider.AdminAddUserToGroupInput{
		UserPoolId: aws.String(os.Getenv("us-east-1_kNKCRvql2")),
		GroupName:  aws.String(groupName),
		Username:   aws.String(username),
	}

	output, err := client.AdminAddUserToGroup(context.Background(), &input)

	if err != nil {
		return nil, err
	}

	return output, nil
}

func CreateNewUser(cfg *aws.Config, username string, roleId string, tp util.Password) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
	client := getClient(cfg)

	input := cognitoidentityprovider.AdminCreateUserInput{
		Username:   aws.String(username),
		UserPoolId: aws.String(os.Getenv("AWS_USER_POOL_ID")),
		DesiredDeliveryMediums: []types.DeliveryMediumType{
			types.DeliveryMediumTypeEmail,
			types.DeliveryMediumTypeSms,
		},
		UserAttributes: []types.AttributeType{
			{Name: aws.String("custom:role"), Value: aws.String(roleId)},
		},
		TemporaryPassword: aws.String(tp.GetPassword()),
	}

	output, err := client.AdminCreateUser(context.TODO(), &input)
	if err != nil {
		return nil, err
	}

	return output, nil
}
