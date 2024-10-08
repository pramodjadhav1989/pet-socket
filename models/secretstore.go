package models

type SecretStore struct {
	JWTKEY          string `json:"JWTKEY"`
	SUPERUSERKEY    string `json:"SUPERUSERKEY"`
	ENCRYPTION      string `json:"ENCRYPTION"`
	BUCKET          string `json:"BUCKET"`
	BUCKETURL       string `json:"BUCKETURL"`
	MYSQLDBSERVER   string `json:"MYSQLDB_SERVER"`
	MYSQLDBUSERID   string `json:"MYSQLDB_USERID"`
	MYSQLDBPASSWORD string `json:"MYSQLDB_PASSWORD"`
	MYSQLDB         string `json:"MYSQLDB"`
}
