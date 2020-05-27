package auth

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"log"
)

var storage *sql.DB

func Init(db *sql.DB) {
	storage = db
}

type Credentials struct {

	// User Id
	Uid int `json:"id"`

	// Organisation Id
	Oid int `json:"org_id"`

	jwt.StandardClaims
}

var secret = []byte("Secret key")

type Authenticator struct {
	db *sql.DB
}

func NewAuthenticator(db *sql.DB) Authenticator {
	return Authenticator{
		db: db,
	}
}

func (a Authenticator) IssueToken(email string, password string) (token string, err error) {
	row := a.db.QueryRow("SELECT id, organisation_id from users where email=$1 and password=$2", email, password)

	var id int
	var orgId int
	err = row.Scan(&id, &orgId)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     id,
		"org_id": orgId,
	})

	return t.SignedString(secret)
}

func ParseToken(token string) *Credentials {
	if token == "" {
		return &(Credentials{})
	}

	t, err := jwt.ParseWithClaims(token, &Credentials{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		log.Printf("Error while jwt token parsing: %v", err)
		return &(Credentials{})
	}

	claims, ok := t.Claims.(*Credentials)

	if !ok || !t.Valid {
		log.Printf("Error while Claims to Credentials transforming, or token is invalid: %v", err)
		return &(Credentials{})
	}

	return claims
}
