package models

type RedisClusterConfig struct {
	Addresses []string `json:"addresses"`
	Name      string   `json:"-"`
}

type FTPClientConnConfig struct {
	Host         string
	User         string
	Password     string
	TimeoutInSec int
}

type FtpFilePathAndConfig struct {
	FilePath           string
	ConfigName         string
	RefactoredFilePath string
}

type WebScraperDataConfig struct {
	ConfigName string
}

type FileTransferUtilityDataConfig struct {
	Ftp                   *FtpFilePathAndConfig
	S3                    *S3ConnConfig
	WebScraper            *WebScraperDataConfig
	PlaceHolder           map[string]*DateFormat
	RunJobTill            *TimeFormat
	PollingTime           int
	DelayTime             int
	TransferType          string
	DownloadWaitTime      int
	WebsitePath           string
	RefactoredWebsitePath string
	CopyFromPast          bool
	LookBackLimit         int
}

type DateFormat struct {
	Format  string
	DaysAdd int
}

type TimeFormat struct {
	Hr  int
	Min int
	Sec int
}

type S3ConnConfig struct {
	BucketName         string
	FilePath           string
	FileName           string
	ConfigName         string
	RefactoredFileName string
	RefactoredFilePath string
}

type KafkaProducerConfig struct {
	Name             string   `json:"name"`
	ServerList       []string `json:"server_list"`
	Servers          string   `json:"servers"`
	Compression      int      `json:"compression"`
	CompressionLevel int      `json:"compressionLevel"`
}
type KafkaConsumerConfig struct {
	Version           string
	Name              string
	ConsumerGroup     string
	ConsumerStrategy  string
	Servers           []string
	Topics            []string // need to change in yaml
	ChannelBufferSize int
}

type AWSClientCredentialsConfig struct {
	Region          string
	AccessId        string
	SecretAccessKey string
}

type SNSPublishMessageJSON struct {
	Message  string
	FileName string
	Error    string
}

type WebScraperConfig struct {
	Url      string
	MemCode  string
	UserId   string
	Password string
}
