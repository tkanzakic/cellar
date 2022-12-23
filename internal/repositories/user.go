package repositories

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/tkanzakic/cellar/internal/core/domain"
)

// AWS DynamoDB Repository

type userDynamoDBRepository struct {
	db *dynamodb.DynamoDB
}

type repositoryUser struct {
	ID             string
	Email          string
	Name           string
	HashedPassword string
}

func domainToRepositoryUser(user *domain.User) *repositoryUser {
	return &repositoryUser{
		ID:             user.ID,
		Email:          user.Email,
		Name:           user.Name,
		HashedPassword: user.Password,
	}
}

func repositoryUserToDomain(user *repositoryUser) *domain.User {
	return &domain.User{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		Password: user.HashedPassword,
	}
}

func NewDynamoDBUserRepository() *userDynamoDBRepository {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(sess)
	return &userDynamoDBRepository{db: db}
}

func (r *userDynamoDBRepository) Get(id string) (*domain.User, error) {
	result, err := r.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Users"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, errors.New("User does not exists")
	}

	repoUser := repositoryUser{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &repoUser)
	if err != nil {
		return nil, err
	}

	user := repositoryUserToDomain(&repoUser)
	return user, nil
}

func (r *userDynamoDBRepository) GetByEmail(email string) (*domain.User, error) {
	user, err := getUser(r, "Email", email)
	return user, err
}

func getUser(r *userDynamoDBRepository, credential, value string) (*domain.User, error) {
	filter := expression.Name(credential).Equal(expression.Value(value))
	projection := expression.NamesList(expression.Name("ID"), expression.Name("Email"), expression.Name("Name"), expression.Name("HashedPassword"))
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		return nil, err
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("Users"),
	}

	result, err := r.db.Scan(params)
	if err != nil {
		return nil, err
	}

	for _, item := range result.Items {
		repoUser := repositoryUser{}
		err = dynamodbattribute.UnmarshalMap(item, &repoUser)
		if err != nil {
			return nil, err
		}

		user := repositoryUserToDomain(&repoUser)
		return user, nil
	}

	return nil, errors.New("No user found")
}

func (r *userDynamoDBRepository) Create(user *domain.User) (*domain.User, error) {
	repoUser := domainToRepositoryUser(user)
	marshalled, err := dynamodbattribute.MarshalMap(repoUser)
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

func (r *userInMemoryRepository) Get(id string) (*domain.User, error) {
	if value, ok := r.kvs[id]; ok {
		user, err := unmarshal(value)
		return &user, err
	}
	return nil, errors.New("User does not exists")
}

func (r *userInMemoryRepository) GetByEmail(email string) (*domain.User, error) {
	for _, value := range r.kvs {
		user, err := unmarshal(value)
		if err == nil && user.Email == email {
			return &user, nil
		}
	}
	return nil, errors.New("User does not exists")
}

func (r *userInMemoryRepository) Create(user *domain.User) (*domain.User, error) {
	value, err := json.Marshal(user)
	if err != nil {
		return nil, errors.New("Error marshalling user")
	}
	r.kvs[user.ID] = value
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
