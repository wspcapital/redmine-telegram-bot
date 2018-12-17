package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"os"
	"strconv"
	"time"
	"redmine-telegram-bot/converse_state"
)

// Format of datetime in database
const dateTimeLayout = "2006-01-02T15:04:05.000Z"

type StateStore struct {
	DynamoDB *dynamodb.DynamoDB
}

// Constructor of StateStore
func New() *StateStore {

	sess := session.Must(session.NewSession(&aws.Config{
		MaxRetries: aws.Int(3),
	}))

	ss := new(StateStore)

	ss.DynamoDB = dynamodb.New(sess, &aws.Config{
		Region: aws.String(os.Getenv("STATE_DYNAMODB_REGION")),
	})

	return ss
}

// Find State in datastore and return it
func (ss *StateStore) GetStateById(id int64) *converse_state.State {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(strconv.FormatInt(id, 10)),
			},
		},
		TableName: aws.String("staffbot_users"),
	}

	result, err := ss.DynamoDB.GetItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				log.Panic(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				log.Panic(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				log.Panic(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				log.Panic(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				log.Panic(aerr.Error())
			}
		}
		log.Panic(err.Error())
	}

	if len(result.Item) == 0 {

		return nil
	}

	state := converse_state.New(ss)
	ss.parseState(state, result.Item)

	return state
}

// Create new state in datastore and return it
func (ss *StateStore) CreateState(id int64) *converse_state.State {
	t := time.Now()
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#IL":  aws.String("IsLogged"),
			"#UAT": aws.String("UpdatedAt"),
			"#CAT": aws.String("CreatedAt"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":il": {
				BOOL: aws.Bool(false),
			},
			":uat": {
				S: aws.String(t.Format(dateTimeLayout)),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(strconv.FormatInt(id, 10)),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		TableName:        aws.String("staffbot_users"),
		UpdateExpression: aws.String("SET #IL = :il, #UAT = :uat, #CAT = :uat"),
	}

	state := converse_state.New(ss)
	ss.sendUpdateStateRequest(state, input)
	state.SetJustCreated(true)

	return state
}

func (ss *StateStore) UpdateState(state *converse_state.State) {
	t := time.Now()
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#IL":  aws.String("IsLogged"),
			"#UAT": aws.String("UpdatedAt"),
			"#UN":  aws.String("UserName"),
			"#CQ":  aws.String("CurrentQuestion"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":il": {
				BOOL: aws.Bool(state.IsLogged),
			},
			":uat": {
				S: aws.String(t.Format(dateTimeLayout)),
			},
			":un": {
				S: aws.String(state.UserName),
			},
			":cq": {
				N: aws.String(strconv.FormatInt(state.CurrentQuestion, 10)),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(strconv.FormatInt(state.ID, 10)),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		TableName:        aws.String("staffbot_users"),
		UpdateExpression: aws.String("SET #IL = :il, #UAT = :uat, #UN = :un, #CQ = :cq"),
	}

	ss.sendUpdateStateRequest(state, input)
}

func (ss *StateStore) parseState(state *converse_state.State, values map[string]*dynamodb.AttributeValue) {

	state.ID, _ = strconv.ParseInt(*values["ID"].S, 10, 64)
	if values["CurrentQuestion"] != nil {
		state.CurrentQuestion, _ = strconv.ParseInt(*values["CurrentQuestion"].N, 10, 64)
	}
	if values["IsLogged"] != nil {
		state.IsLogged = *values["IsLogged"].BOOL
	}
	if values["UserName"] != nil {
		state.UserName = *values["UserName"].S
	}
	if values["CreatedAt"] != nil {
		state.CreatedAt, _ = time.Parse(dateTimeLayout, *values["CreatedAt"].S)
	}
	if values["UpdatedAt"] != nil {
		state.UpdatedAt, _ = time.Parse(dateTimeLayout, *values["UpdatedAt"].S)
	}

}

func (ss *StateStore) sendUpdateStateRequest(state *converse_state.State, input *dynamodb.UpdateItemInput) {
	result, err := ss.DynamoDB.UpdateItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				log.Panic(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				log.Panic(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				log.Panic(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				log.Panic(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeTransactionConflictException:
				log.Panic(dynamodb.ErrCodeTransactionConflictException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				log.Panic(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				log.Panic(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				log.Panic(aerr.Error())
			}
		}
		log.Panic(err.Error())
	}

	ss.parseState(state, result.Attributes)
}
