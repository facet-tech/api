package backend

import (
	"errors"
	"facet/api/db"
	"facet/api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Annotation struct {
	ClassName  string            `json:"className"`
	Parameters []Parameter       `json:"parameters,omitempty"`
	Visibility map[string]string `json:"visibility,omitempty"`
}

type CircuitBreaker struct {
	Precedence      int      `json:"precedence"`
	MethodsToCreate []Method `json:"methodsToCreate"`
	Toggle          Toggle   `json:"Toggle"`
	ReturnType      string   `json:"returnType"`
}

type Exception struct {
	ClassName string `json:"className"`
}

type Framework struct {
	Name            string           `json:"name"`
	CircuitBreakers []CircuitBreaker `json:"circuitBreakers"`
	Sensors         []Sensor         `json:"sensors"`
	Version         string           `json:"version"`
}

type Language struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Method struct {
	Annotations []Annotation `json:"annotations"`
	Body        string       `json:"body"`
	Exceptions  []Exception  `json:"exceptions"`
	Modifier    string       `json:"modifier"`
	Name        string       `json:"name"`
	Parameters  []Parameter  `json:"parameters"`
	ReturnType  string       `json:"returnType"`
}

type Parameter struct {
	ClassName string      `json:"className"`
	Name      string      `json:"name"`
	Values    interface{} `json:"values"`
	Value     string      `json:"value`
	Type      interface{} `json:"type"`
	Position  int         `json:"position"`
}

type Toggle struct {
	Method           []Method          `json:"method"`
	ParameterMapping map[string]string `json:"parameterMapping"`
}

type Sensor struct {
	Annotations []Annotation `json:"annotations"`
	ReturnType  string       `json:"returnType"`
}

type Signature struct {
	Enabled    bool         `json:"enabled"`
	Name       string       `json:"name"`
	Parameter  []Parameter  `json:"parameter"`
	ReturnType string       `json:"returnType"`
	Signature  string       `json:"signature"`
	Annotation []Annotation `json:"annotation,omitempty"`
}

type DTO struct {
	AppId              string                 `json:"appId"`
	FullyQualifiedName string                 `json:"fullyQualifiedName"`
	Signature          []Signature            `json:"signature"`
	Version            string                 `json:"version"`
	Attribute          map[string]interface{} `json:"attribute,omitempty"`
	Language           Language               `json:"language"`
	Framework          []Framework            `json:"framework,omitempty"`
	InterfaceSignature []string               `json:"interfaceSignature,omitempty"`
	Type               string                 `json:"type,omitempty"`
	ParentSignature    string                 `json:"parentSignature,omitempty"`
	Annotation         []Annotation           `json:"annotation,omitempty"`
}

func FetchAll(AppId string) (*[]DTO, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(db.BackendTableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"appId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(AppId),
					},
				},
			},
		},
	}
	result, error := db.Database.Query(input)
	facets := new([]DTO)
	if error == nil && result.Items != nil {
		error = dynamodbattribute.UnmarshalListOfMaps(result.Items, facets)
	}
	return facets, error
}

func (facet *DTO) fetch() error {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.BackendTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"appId": {
				S: aws.String(facet.AppId),
			},
			"fullyQualifiedName": {
				S: aws.String(facet.FullyQualifiedName),
			},
		},
	}
	result, error := db.Database.GetItem(input)
	if error == nil && result != nil {
		if len(result.Item) == 0 {
			error = errors.New(util.NOT_FOUND)
		} else {
			error = dynamodbattribute.UnmarshalMap(result.Item, facet)
		}
	}
	return error
}

func (facet *DTO) create() error {
	item, error := dynamodbattribute.MarshalMap(facet)
	if error == nil {
		input := &dynamodb.PutItemInput{
			TableName: aws.String(db.BackendTableName),
			Item:      item,
		}
		_, error = db.Database.PutItem(input)
	}
	return error
}

func (facet *DTO) delete() error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.BackendTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"appId": {
				S: aws.String(facet.AppId),
			},
			"fullyQualifiedName": {
				S: aws.String(facet.FullyQualifiedName),
			},
		},
	}
	_, error := db.Database.DeleteItem(input)
	return error
}

func (facet *DTO) deleteAll() error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.BackendTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"appId": {
				S: aws.String(facet.AppId),
			},
		},
	}
	_, error := db.Database.DeleteItem(input)
	return error
}
