package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/COOLizh/sqs-sns-examples/internal/event"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/velmie/broker"
	snsBroker "github.com/velmie/broker/sns"
	sqsBroker "github.com/velmie/broker/sqs"
)

const (
	longPollingDuration = 20
	awsRegion           = "us-east-1"
	topicName           = "topic_card_created"
	queueName           = "queue_card_created"
)

func main() {
	logger := logrus.New()

	httpClient := http.Client{Timeout: longPollingDuration * time.Second}

	// aws config creation
	awsConfig := &aws.Config{
		HTTPClient: &httpClient,
		Region:     aws.String(awsRegion),
	}

	// aws session init
	awsSession, err := session.NewSession(awsConfig)
	if err != nil {
		logger.Errorf("cannot create AWS session: %s", err)
		os.Exit(1)
	}

	// getting aws account id
	accountID, err := getAWSAccountID(awsSession)
	if err != nil {
		logger.Errorf("cannot get AWS account ID: %s", err)
		os.Exit(1)
	}

	// configure publisher
	snsService := sns.New(awsSession)
	snsPublisher := snsBroker.NewPublisher(snsService, "", accountID)
	publisher := event.NewBrokerPublisher(snsPublisher)

	// configure subscriber
	sqsService := sqs.New(awsSession)
	sqsSubscriber := sqsBroker.NewSubscriber(
		sqsService,
		sqsBroker.LongPollingDuration(longPollingDuration),
		sqsBroker.RequestMultipleMessage(1),
	)
	subscriber := event.NewSQSSubscriber(logger, sqsSubscriber)

	mockExit := false

	// event handler for incoming messages from queue
	subscriber.AddSubscription(queueName, func(e broker.Event) error {
		logger.Infof("got new message. Body %s", string(e.Message().Body))
		mockExit = true
		return nil
	})
	err = subscriber.SubscribeAll()
	if err != nil {
		logger.Errorf("cannot subscribe to queue: %s", err)
		os.Exit(1)
	}
	defer func() {
		err := subscriber.UnsubscribeAll()
		if err != nil {
			logger.Errorf("cannot unsubscribe from queue: %s", err)
			os.Exit(1)
		}
	}()

	// message publishing
	err = publisher.OnCardCreated(context.Background(), topicName, &event.CardCreated{
		CardID: uuid.New().String(),
		UserID: uuid.New().String(),
	})
	if err != nil {
		logger.Errorf("cannot publish message: %s", err)
		os.Exit(1)
	}

	for !mockExit {
	}
}

func getAWSAccountID(awsSession *session.Session) (string, error) {
	stsService := sts.New(awsSession)
	out, err := stsService.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return *out.Account, nil
}
