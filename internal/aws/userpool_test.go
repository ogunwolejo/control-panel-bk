package aws

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/stretchr/testify/assert"
)

type mockCognitoClient struct {
	CreateGroupFunc         func(ctx context.Context, input *cognitoidentityprovider.CreateGroupInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateGroupOutput, error)
	AdminAddUserToGroupFunc func(ctx context.Context, input *cognitoidentityprovider.AdminAddUserToGroupInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error)
	AdminCreateUserFunc     func(ctx context.Context, input *cognitoidentityprovider.AdminCreateUserInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error)
}

func (m *mockCognitoClient) CreateGroup(ctx context.Context, input *cognitoidentityprovider.CreateGroupInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateGroupOutput, error) {
	return m.CreateGroupFunc(ctx, input, opts...)
}

func (m *mockCognitoClient) AdminAddUserToGroup(ctx context.Context, input *cognitoidentityprovider.AdminAddUserToGroupInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
	return m.AdminAddUserToGroupFunc(ctx, input, opts...)
}

func (m *mockCognitoClient) AdminCreateUser(ctx context.Context, input *cognitoidentityprovider.AdminCreateUserInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
	return m.AdminCreateUserFunc(ctx, input, opts...)
}

var mockClient mockCognitoClient

func TestCreateUserPoolGroup(t *testing.T) {
	mockClient = mockCognitoClient{
		CreateGroupFunc: func(ctx context.Context, input *cognitoidentityprovider.CreateGroupInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateGroupOutput, error) {
			if input.GroupName == nil || *input.GroupName == "" {
				return nil, errors.New("group name is required")
			}
			return &cognitoidentityprovider.CreateGroupOutput{}, nil
		},
	}

	cfg := aws.Config{}
	os.Setenv("us-east-1_kNKCRvql2", "test-user-pool-id")
	group := Group{Name: "Admins", Description: "Admin group"}

	output, err := CreateUserPoolGroup(&cfg, group)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func TestAddUsersToUserPoolGroup(t *testing.T) {
	mockClient = mockCognitoClient{
		AdminAddUserToGroupFunc: func(ctx context.Context, input *cognitoidentityprovider.AdminAddUserToGroupInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
			if input.Username == nil || *input.Username == "" {
				return nil, errors.New("username is required")
			}
			return &cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil
		},
	}

	cfg := aws.Config{}
	os.Setenv("us-east-1_kNKCRvql2", "test-user-pool-id")
	output, err := AddUsersToUserPoolGroup(&cfg, "Admins", "testuser")

	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func TestCreateNewUser(t *testing.T) {
	mockClient = mockCognitoClient{
		AdminCreateUserFunc: func(ctx context.Context, input *cognitoidentityprovider.AdminCreateUserInput, opts ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
			if input.Username == nil || *input.Username == "" {
				return nil, errors.New("username is required")
			}
			return &cognitoidentityprovider.AdminCreateUserOutput{}, nil
		},
	}

	cfg := aws.Config{}
	os.Setenv("AWS_USER_POOL_ID", "test-user-pool-id")
	output, err := CreateNewUser(&cfg, "testuser")

	assert.NoError(t, err)
	assert.NotNil(t, output)
}
