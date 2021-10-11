package corsns

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// Publisher ...
type Publisher interface {
	// PublishInput sends a full input mesage through Sns
	PublishInput(ctx context.Context, input *sns.PublishInput) (string, error)
	// Publish ...
	Publish(ctx context.Context, body interface{}) (string, error)
}

type publisher struct {
	client *sns.SNS
	config PublisherConfig
}

// NewPublisher ...
func NewPublisher(sess *session.Session, config PublisherConfig) Publisher {
	// set defaults
	if len(config.Encoding) == 0 {
		config.Encoding = PublisherEncodingJSON
	}

	return &publisher{
		client: sns.New(sess),
		config: config,
	}
}

// PublishInput ...
func (p *publisher) PublishInput(ctx context.Context, input *sns.PublishInput) (string, error) {
	input.TopicArn = aws.String(p.config.TopicARN)
	output, err := p.client.PublishWithContext(ctx, input)
	if err != nil {
		return "", err
	}

	return *output.MessageId, nil
}

// Publish ...
func (p *publisher) Publish(ctx context.Context, body interface{}) (string, error) {
	encoded, err := p.encodeBody(ctx, body)
	if err != nil {
		return "", err
	}
	return p.PublishInput(ctx, &sns.PublishInput{
		Message: aws.String(encoded),
	})
}

func (p *publisher) encodeBody(ctx context.Context, body interface{}) (string, error) {
	if p.config.Encoding == PublisherEncodingJSON {
		return p.encodeJSON(ctx, body)
	} else if p.config.Encoding == PublisherEncodingRaw {
		return fmt.Sprintf("%v", body), nil
	} else {
		return "", fmt.Errorf("invalid encoding type=%s", p.config.Encoding)
	}
}

func (p *publisher) encodeJSON(ctx context.Context, body interface{}) (string, error) {
	bytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal json: %w", err)
	}
	return string(bytes), nil
}
