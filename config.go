package corsns

// PublisherEncodingType ...
type PublisherEncodingType string

const (
	// PublisherE ...
	PublisherEncodingRaw PublisherEncodingType = "raw"
	// PublisherEncodingJSON ...
	PublisherEncodingJSON PublisherEncodingType = "json"
)

// PublisherConfig ...
type PublisherConfig struct {
	// Encoding ...
	Encoding PublisherEncodingType
	// TopicARN for the publisher
	TopicARN string
}
