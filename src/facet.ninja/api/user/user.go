package user

import (
	"sync"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

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

var lock = &sync.Mutex{}
var (
	identityMetadata IdentityMetadata
)

func GetIdentityMetadata() IdentityMetadata {
	var once sync.Once
	once.Do(func() {
		fmt.Println("ETREKSA")
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		identityMetadata = IdentityMetadata{
			CognitoClient:   cognito.New(sess),
			UserPoolID:      "us-west-2_vnM0aVcxD",
			AppClientID:     "not-added-yet",
			AppClientSecret: "not-added-yet",
		}
	})
	return identityMetadata
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
	if error == nil {
		return user.addUserToUserPool()
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

	_, err := cognitoClient.AdminCreateUser(newUserData)
	return err
}

func (user *User) login() error {
	identityMetadata = GetIdentityMetadata()
	fmt.Println("CHECKAREME", identityMetadata)
	const flowUsernamePassword = "USER_PASSWORD_AUTH"

	username := user.Email
	password := user.Password

	flow := aws.String(flowUsernamePassword)
	params := map[string]*string{
		"USERNAME": aws.String(username),
		"PASSWORD": aws.String(password),
	}

	authTry := &cognito.InitiateAuthInput{
		AuthFlow:       flow,
		AuthParameters: params,
	}

	res, err := identityMetadata.CognitoClient.InitiateAuth(authTry)
	if err != nil {
		fmt.Println(err)
		//http.Redirect(w, r, fmt.Sprintf("/login?message=%s", err.Error()), http.StatusSeeOther)

	}
	fmt.Println(res)
	return err
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
