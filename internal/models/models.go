package models

import (
	"context"
	"html/template"
	"net/http"
	"time"
)

type Middleware func(handler http.HandlerFunc) http.HandlerFunc

type Session struct {
	UserID         int       `json:"user_id"`
	ConnectionID   int       `json:"connection_id"`
	Pseudo         string    `json:"Pseudo"`
	IpAddress      string    `json:"ip_address"`
	ExpirationTime time.Time `json:"expiration_time"`
}

var (
	Ctx context.Context
	//Db  *sql.DB
)
var Tmpl = make(map[string]*template.Template)

type TConf_Mysql struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

var Conf_Mysql TConf_Mysql

const (
	Port = ":8080"

// NomBase = "basemw8"
)

type Credentials struct {
	Username string
	Password string
}
type User struct {
	Id        int    `json:"id"`
	Pseudo    string `json:"Pseudo"`
	Droits    string `json:"Droits"`
	HashedPwd string `json:"hash"`
	Salt      string `json:"salt"`
	Email     string `json:"email"`
}

type TempUser struct {
	ConfirmID    string
	CreationTime time.Time
	User         User
}

type MailConfig struct {
	Email    string `json:"email_addr"`
	Auth     string `json:"email_auth"`
	Hostname string `json:"host"`
	Port     int    `json:"port"`
}
