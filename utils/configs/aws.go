package configs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	log "github.com/smartpet/websocket/utils/logger"
)

var awscfg aws.Config

const DefaultRegion = "ap-south-1"

func InitAWS(ctx context.Context) error {
	var err error
	// var region = os.Getenv(constant.AWSRegionKey)
	// if region == "" {
	region := DefaultRegion
	// }
	log.ApplicationInfo(ctx).Msg("AWS LoadDefaultConfig started")
	awscfg, err = awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(region))
	if err != nil {
		log.ApplicationError(ctx).Msg(err.Error())
		return err
	}

	// perform loading of aws configs only during release mode

	// initialize AWS secrets - passwords are stored in secrets
	InitSecrets(ctx)

	//Note: Except secrets - do not add any other AWS service initialization here; rest all should be part of lazy load
	return err
}

// Returns the AWS config to be used by server initialization
func GetConfig() aws.Config {
	return awscfg
}
