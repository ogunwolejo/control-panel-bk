package aws

import (
	"context"
	"control-panel-bk/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"log"
	"os"
)

// awsApplication
// arn:aws:resource-groups:us-east-1:528757792684:group/controlPanel/04a7tbp0n21ntgacp58yscw9e0

type Group struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PoolUser struct {
	Username string `json:"username"`
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

func CreateNewUser(cfg *aws.Config, username string, roleId string, tp util.Password) (*string, error) {
	client := getClient(cfg)

	input := cognitoidentityprovider.AdminCreateUserInput{
		Username:   aws.String(username),
		UserPoolId: aws.String(os.Getenv("AWS_USER_POOL_ID")),
		DesiredDeliveryMediums: []types.DeliveryMediumType{
			types.DeliveryMediumTypeEmail,
		},
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: aws.String(username)},
			//{Name: aws.String("password"), Value: aws.String()},
			{Name: aws.String("custom:role"), Value: aws.String(roleId)},
		},
		TemporaryPassword: aws.String(tp.GetPassword()),
	}

	output, err := client.AdminCreateUser(context.TODO(), &input)
	if err != nil {
		log.Println("error: ", err)
		return nil, err
	}

	var userSub string
	for _, attr := range output.User.Attributes {
		if *attr.Name == "sub" {
			userSub = *attr.Value
			break
		}
	}

	return &userSub, nil
}

func DeleteUser(cfg *aws.Config, username string) (*cognitoidentityprovider.AdminDeleteUserOutput, error) {
	client := getClient(cfg)

	input := cognitoidentityprovider.AdminDeleteUserInput{
		Username:   aws.String(username),
		UserPoolId: aws.String(os.Getenv("AWS_USER_POOL_ID")),
	}

	out, err := client.AdminDeleteUser(context.TODO(), &input)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func DisableUser(cfg *aws.Config, username string) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
	client := getClient(cfg)

	input := cognitoidentityprovider.AdminDisableUserInput{
		Username:   aws.String(username),
		UserPoolId: aws.String(os.Getenv("AWS_USER_POOL_ID")),
	}

	output, err := client.AdminDisableUser(context.TODO(), &input)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func ActivateUser(cfg *aws.Config, username string) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
	client := getClient(cfg)

	input := cognitoidentityprovider.AdminEnableUserInput{
		Username:   aws.String(username),
		UserPoolId: aws.String(os.Getenv("AWS_USER_POOL_ID")),
	}

	output, err := client.AdminEnableUser(context.TODO(), &input)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func AuthViaRefreshToken(cfg *aws.Config, clientId, refreshToken string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	client := getClient(cfg)

	input := cognitoidentityprovider.InitiateAuthInput{
		ClientId: aws.String(clientId),
		AuthFlow: types.AuthFlowTypeRefreshTokenAuth,
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
		},
	}

	output, err := client.InitiateAuth(context.TODO(), &input)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func GetUserDetails(cfg *aws.Config, accessToken string) (*cognitoidentityprovider.GetUserOutput, error) {
	client := cognitoidentityprovider.NewFromConfig(*cfg)

	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(accessToken),
	}

	output, err := client.GetUser(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func LogInUser(cfg *aws.Config, clientId, email, password string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	client := getClient(cfg)

	input := cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(clientId),
		AuthParameters: map[string]string{
			"USERNAME": email,
			"PASSWORD": password,
		},
	}

	output, err := client.InitiateAuth(context.TODO(), &input)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func LogOutUser(cfg *aws.Config, token string) error {
	client := getClient(cfg)

	if _, err := client.GlobalSignOut(context.TODO(), &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String(token),
	}); err != nil {
		return err
	}

	return nil
}

func ChangeUserPassword(cfg *aws.Config, token, proposedPassword, oldPassword string) (*cognitoidentityprovider.ChangePasswordOutput, error) {
	client := getClient(cfg)

	input := cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      aws.String(token),
		ProposedPassword: aws.String(proposedPassword),
		PreviousPassword: aws.String(oldPassword),
	}

	output, err := client.ChangePassword(context.TODO(), &input)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func ForgetPasswordOtp(cfg *aws.Config, email string) (*cognitoidentityprovider.ForgotPasswordOutput, error) {
	client := getClient(cfg)

	fgInput := cognitoidentityprovider.ForgotPasswordInput{
		ClientId: aws.String(os.Getenv("AWS_CLIENT_ID")),
		Username: aws.String(email),
	}

	output, err := client.ForgotPassword(context.TODO(), &fgInput)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func ForgetPassword(cfg *aws.Config, email, otp, password string) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
	client := getClient(cfg)

	input := cognitoidentityprovider.ConfirmForgotPasswordInput{
		Username: aws.String(email),
		ClientId: aws.String(os.Getenv("AWS_CLIENT_ID")),
		Password: aws.String(password),
		ConfirmationCode: aws.String(otp),
	}

	return client.ConfirmForgotPassword(context.TODO(), &input)
}