package models

type LoginResponse struct {
	Token   string `json:"token"`
	TokenId string `json:"tokenId"`
}

type Login struct {
	ID         int    `json:"id"`
	UserId     string `json:"userid"`
	Isparent   int    `json:"isparent"`
	Isprovider int    `json:"isprovider"`
	Deviceid   string `json:"deviceid"`
	Createdon  string `json:"createdon"`
	Modifiedon string `json:"modifiedon"`
	Isactive   int    `json:"isactive"`
}

type RefreshTokenRequest struct {
	Token    string `json:"token"`
	Deviceid string `json:"deviceid"`
}
