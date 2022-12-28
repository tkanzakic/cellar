package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/tkanzakic/cellar/internal/core/domain"
)

// AWS DynamoDB Repository

type userDynamoDBRepository struct {
	client *dynamodb.Client
}

func NewDynamoDBUserRepository() *userDynamoDBRepository {
	host := os.Getenv("HOST")
	var client *dynamodb.Client
	if host == "localhost" {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			log.Fatalf("Error loading config %v", err)
		}
		client = dynamodb.NewFromConfig(cfg)
	} else {
		client = dynamodb.NewFromConfig(aws.Config{
			Credentials: ec2rolecreds.New(),
			Region:      "us-east-1",
		})
	}
	return &userDynamoDBRepository{client: client}
}

func getUserKey(family, email string) map[string]types.AttributeValue {
	f, err := attributevalue.Marshal(family)
	if err != nil {
		log.Printf("Error marshalling family. %v", err)
		return nil
	}
	e, err := attributevalue.Marshal(email)
	if err != nil {
		log.Printf("Error marshaling email. %v", err)
		return nil
	}
	return map[string]types.AttributeValue{"Family": f, "Email": e}
}

func (r *userDynamoDBRepository) GetByEmail(family, email string) (*domain.User, error) {
	result, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       getUserKey(family, email),
		TableName: aws.String("Users"),
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, errors.New("User does not exists")
	}

	user := domain.User{}
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userDynamoDBRepository) Create(user *domain.User) (*domain.User, error) {
	marshalled, err := attributevalue.MarshalMap(user)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      marshalled,
		TableName: aws.String("Users"),
	}
	_, err = r.client.PutItem(context.TODO(), input)
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
