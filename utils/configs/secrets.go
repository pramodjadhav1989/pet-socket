package configs

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/smartpet/websocket/constant"
	log "github.com/smartpet/websocket/utils/logger"
)

func InitSecrets(ctx context.Context) {
	awsConfig := GetConfig()

	// initialize all secrets store one-by-one
	initSecrets(ctx, awsConfig)
}

// initialize master
func initSecrets(ctx context.Context, awsConfig aws.Config) {

	if err := initStore(ctx, awsConfig, constant.AWSAccessKeyID, constant.AWSSecretAccessKey); err != nil {
		log.Warn(ctx).Msg("secrets manager failed to load.. will skip this and fall back to ENV mode")
		return
	}

	// if no erros , then override environment variables
	//Get().ReloadEnvironment()
}

// performs secrets manager initialization
// if any errors, it is logged and the functionality is suppressed by falling back to previous mode of execution
func initStore(ctx context.Context, awsConfig aws.Config, sn string, uv string) error {

	svc := secretsmanager.NewFromConfig(awsConfig)
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(os.Getenv(constant.SecrateName)),
		VersionStage: aws.String("AWSCURRENT"),
	}
	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		log.Error(ctx).Err(err).Stack().Msgf("Failed to read SecretManager values for: '%s' ", sn)
		return err
	}
	//data := models.SecretStore{}

	data := make(map[string]string)
	if err := json.Unmarshal([]byte(*result.SecretString), &data); err != nil {
		log.Error(ctx).Err(err).Stack().Msgf("Failed to parse SecretManager string for: '%s", sn)
		return err
	}

	// GetClient().setSecrets(data)

	for key, val := range data {
		os.Setenv(key, val)
	}
	// log.ApplicationInfo(ctx).Msgf("MYSQLDBSERVER %s", data.MYSQLDBSERVER)
	// log.ApplicationInfo(ctx).Msgf("BUCKETURL %s", data.BUCKETURL)

	// os.Setenv("JWTKEY", data.JWTKEY)
	// os.Setenv("SUPERUSERKEY", data.SUPERUSERKEY)
	// os.Setenv("ENCRYPTION", data.ENCRYPTION)
	// os.Setenv("BUCKET", data.BUCKET)
	// os.Setenv("BUCKETURL", data.BUCKETURL)
	// os.Setenv("MYSQLDB_SERVER", data.MYSQLDBSERVER)
	// os.Setenv("MYSQLDB_USERID", data.MYSQLDBUSERID)
	// os.Setenv("MYSQLDB_PASSWORD", data.MYSQLDBPASSWORD)
	// os.Setenv("MYSQLDB", data.MYSQLDB)

	// os.Setenv(constants.ProfileServiceTokenKey, data.ProfileServiceApiToken)

	return nil
}
