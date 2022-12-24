package repositories

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/tkanzakic/cellar/internal/core/domain"
)

// AWS DynamoDB Repository

type userDynamoDBRepository struct {
	db *dynamodb.DynamoDB
}

func NewDynamoDBUserRepository() *userDynamoDBRepository {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(sess)
	return &userDynamoDBRepository{db: db}
}

func (r *userDynamoDBRepository) GetByEmail(family, email string) (*domain.User, error) {
	result, err := r.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Users"),
		Key: map[string]*dynamodb.AttributeValue{
			"Family": {
				S: aws.String(family),
			},
			"Email": {
				S: aws.String(email),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, errors.New("User does not exists")
	}

	user := domain.User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userDynamoDBRepository) Create(user *domain.User) (*domain.User, error) {
	marshalled, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      marshalled,
		TableName: aws.String("Users"),
	}
	_, err = r.db.PutItem(input)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// In Memory Repository

type userInMemoryRepository struct {
	kvs map[string][]byte
}

func NewInMemoryUserRepository() *userInMemoryRepository {
	return &userInMemoryRepository{kvs: map[string][]byte{}}
}

func getInMemoryKey(family, email string) string {
	return family + "-" + email
}

func (r *userInMemoryRepository) GetByEmail(family, email string) (*domain.User, error) {
	if value, ok := r.kvs[getInMemoryKey(family, email)]; ok {
		user, err := unmarshal(value)
		return &user, err
	}
	return nil, errors.New("User does not exists")
}

func (r *userInMemoryRepository) Create(user *domain.User) (*domain.User, error) {
	value, err := json.Marshal(user)
	if err != nil {
		return nil, errors.New("Error marshalling user")
	}
	r.kvs[getInMemoryKey(user.Family, user.Email)] = value
	return user, nil
}

func unmarshal(value []byte) (domain.User, error) {
	user := domain.User{}
	err := json.Unmarshal(value, &user)
	if err != nil {
		return domain.User{}, errors.New("Could not unmarshal user")
	}
	return user, nil
}
