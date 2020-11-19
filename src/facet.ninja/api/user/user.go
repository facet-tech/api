package user

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"sync"

	"facet.ninja/api/db"
	"facet.ninja/api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type User struct {
	Id          string                 `json:"id"`
	WorkspaceId string                 `json:"workspaceId"`
	Email       string                 `json:"email"`
	Attribute   map[string]interface{} `json:"attribute,omitempty"`
	Password    string                 `json:"password,omitempty"`
	Username    string                 `json:"username,omitempty"`
}

const (
	KEY_USER    = "USER"
	EMAIL_INDEX = "email-index"
)

type IdentityMetadata struct {
	CognitoClient   *cognito.CognitoIdentityProvider
	UserPoolID      string
	AppClientID     string
	AppClientSecret string
}

var (
	identityMetadata IdentityMetadata
)

func GetIdentityMetadata() IdentityMetadata {
	var once sync.Once
	once.Do(func() { //
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		identityMetadata = IdentityMetadata{
			CognitoClient:   cognito.New(sess),
			UserPoolID:      "us-west-2_vnM0aVcxD",
			AppClientID:     "1j7ufkr3pfj82o265ll6iqr6hp",
			AppClientSecret: "lak12dk2c6knbqjnhln637l23041lfae6hthoh2fqdoj5vn6cjv",
		}
	})
	return identityMetadata
}

func computeSecretHash(clientSecret string, username string, clientId string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (user *User) fetch() error {
	input := &dynamodb.QueryInput{
		TableName: aws.String(db.WorkspaceTableName),
		IndexName: aws.String(EMAIL_INDEX),
		KeyConditions: map[string]*dynamodb.Condition{
			"email": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(user.Email),
					},
				},
			},
		},
	}
	result, error := db.Database.Query(input)
	if error == nil && result != nil {
		if len(result.Items) == 0 {
			error = errors.New(util.NOT_FOUND)
		} else {
			error = dynamodbattribute.UnmarshalMap(result.Items[0], user)
		}
	}
	return error
}

func (user *User) create() error {
	user.Id = db.CreateId(KEY_USER)
	user.Password = "" //not storing passwords
	item, error := dynamodbattribute.MarshalMap(user)
	if error == nil {
		input := &dynamodb.PutItemInput{
			TableName: aws.String(db.WorkspaceTableName),
			Item:      item,
		}
		_, error = db.Database.PutItem(input)
	}
	return error
}

func (user *User) addUserToUserPool() error {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	cognitoClient := cognitoidentityprovider.New(sess)

	newUserData := &cognitoidentityprovider.AdminCreateUserInput{
		DesiredDeliveryMediums: []*string{
			aws.String("EMAIL"),
		},
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(user.Email),
			},
		},
	}

	// TODO read from env variable
	newUserData.SetUserPoolId("us-west-2_vnM0aVcxD")
	newUserData.SetUsername(user.Email)

	res, err := cognitoClient.AdminCreateUser(newUserData)
	fmt.Println("CHECK", res, "err", err)
	return err
}

func (user *User) verify() error {
	identityMetadata = GetIdentityMetadata()

	const flowUsernamePassword = "USER_PASSWORD_AUTH"

	username := user.Email
	password := user.Password

	fmt.Println("CHECKAREME", username, password)

	flow := aws.String(flowUsernamePassword)
	params := map[string]*string{
		"USERNAME": aws.String(username),
		"PASSWORD": aws.String(password),
	}

	if identityMetadata.AppClientSecret != "" {
		secretHash := computeSecretHash(identityMetadata.AppClientSecret, username, identityMetadata.AppClientID)

		params["SECRET_HASH"] = aws.String(secretHash)
	}

	authTry := &cognito.InitiateAuthInput{
		AuthFlow:       flow,
		AuthParameters: params,
		ClientId:       aws.String(identityMetadata.AppClientID),
	}

	res, err := identityMetadata.CognitoClient.InitiateAuth(authTry)

	if err != nil {
		fmt.Println(err)
	}

	params["NEW_PASSWORD"] = aws.String(user.Password)

	authChallenge := &cognito.RespondToAuthChallengeInput{
		ChallengeName:      aws.String(cognito.ChallengeNameTypeNewPasswordRequired),
		ChallengeResponses: params,
		ClientId:           aws.String(identityMetadata.AppClientID),
		Session:            res.Session,
	}

	fmt.Println("authChallenge", authChallenge)

	req, output := identityMetadata.CognitoClient.RespondToAuthChallenge(authChallenge)

	fmt.Println("req", req, "output: ", output)

	fmt.Println("Final res", res)

	return err
}

func (user *User) login() *cognito.InitiateAuthOutput {
	identityMetadata = GetIdentityMetadata()

	const flowUsernamePassword = "USER_PASSWORD_AUTH"

	username := user.Email
	password := user.Password

	flow := aws.String(flowUsernamePassword)
	params := map[string]*string{
		"USERNAME": aws.String(username),
		"PASSWORD": aws.String(password),
	}

	if identityMetadata.AppClientSecret != "" {
		secretHash := computeSecretHash(identityMetadata.AppClientSecret, username, identityMetadata.AppClientID)

		params["SECRET_HASH"] = aws.String(secretHash)
	}

	authTry := &cognito.InitiateAuthInput{
		AuthFlow:       flow,
		AuthParameters: params,
		ClientId:       aws.String(identityMetadata.AppClientID),
	}

	res, err := identityMetadata.CognitoClient.InitiateAuth(authTry)
	fmt.Println("ELA RE res", res, "err", err)
	return res
}

func (user *User) delete() error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.WorkspaceTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"workspaceId": {
				S: aws.String(user.WorkspaceId),
			},
			"id": {
				S: aws.String(user.Id),
			},
		},
	}
	_, err := db.Database.DeleteItem(input)
	return err
}
