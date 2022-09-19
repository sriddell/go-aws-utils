package awsutils

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var once sync.Once
var cachedConfig *aws.Config
var cachedError error

func GetConfig(region *string, endpoint *string) (*aws.Config, error) {
	once.Do(func() { initConfig(region, endpoint) })

	return cachedConfig, cachedError
}

func initConfig(region *string, endpoint *string) {
	if endpoint == nil {
		config, err := config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(*region))
		if err != nil {
			cachedError = err
		} else {
			cachedConfig = &config
		}
	} else {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           *endpoint,
				SigningRegion: region,
			}, nil
		})
		config, err := config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(*region),
			config.WithEndpointResolverWithOptions(customResolver))
		if err != nil {
			cachedError = err
		} else {
			cachedConfig = &config
		}
	}
}
