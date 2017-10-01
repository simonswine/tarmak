package amazon

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/spf13/viper"

	tarmakv1alpha1 "github.com/jetstack/tarmak/pkg/apis/tarmak/v1alpha1"
)

type Amazon struct {
	flags *viper.Viper
	log   *logrus.Entry

	queueURL string
	session  *session.Session
	sqs      SQS
}

type SQS interface {
	SendMessage(*sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
}

func New(flags *viper.Viper, log logrus.Entry) *Amazon {
	return &Amazon{
		log:   log.WithField("provider", "amazon"),
		flags: flags,
	}
}

func (a *Amazon) Session() (*session.Session, error) {

	// return cached session
	if a.session != nil {
		return a.session, nil
	}

	return session.NewSession()
}

func (a *Amazon) SQS() (SQS, error) {
	if a.sqs == nil {
		sess, err := a.Session()
		if err != nil {
			return nil, fmt.Errorf("error getting Amazon session: %s", err)
		}
		a.sqs = sqs.New(sess)
	}
	return a.sqs, nil
}

func (a *Amazon) SendReport(report *tarmakv1alpha1.InstanceState) error {

	svc, err := a.SQS()
	if err != nil {
		return fmt.Errorf("error getting SQS service: %s", err)
	}

	result, err := svc.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"type": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(tarmakv1alpha1.InstanceStateTypeReport),
			},
		},
		MessageBody: aws.String("Information about current NY Times fiction bestseller for week of 12/11/2016."),
		QueueUrl:    &a.queueURL,
	})

	if err != nil {
		return fmt.Errorf("error sending report: %s", err)
	}

	a.log.Debugf("report sent correctly message_id=%s", *result.MessageId)
	return nil
}

func (a *Amazon) ReceiveUpdate(chan *tarmakv1alpha1.InstanceState) error {
	return nil
}
