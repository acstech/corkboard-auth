package corkboardauth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"github.com/julienschmidt/httprouter"
)

//CorkboardAuth is an instance of the Corkboard Auth server
type CorkboardAuth struct {
	//bucket     *gocb.Bucket
	privateKey *rsa.PrivateKey
}

//Config is the configuration values for creating a new CorkboardAuth
type Config struct {
	//CBConnection   string
	//CBBucket       string
	//CBBucketPass   string
	PrivateRSAFile string
}

//TODO: Definitely need to change "New" so it points to the correct
//database

//Can I put environment variable or an RSA key into a Heroku app?

//New creates a new CorkboardAuth
func New(config *Config) (*CorkboardAuth, error) {
	// cluster, err := gocb.Connect(config.CBConnection)
	// if err != nil {
	// 	return nil, err
	// }
	// bucket, err := cluster.OpenBucket(config.CBBucket, config.CBBucketPass)
	// if err != nil {
	// 	return nil, err
	// }

	//TODO: connect to the database

	privateKey, err := getPrivate(config.PrivateRSAFile)
	if err != nil {
		return nil, err
	}
	return &CorkboardAuth{
		//bucket:     bucket,
		privateKey: privateKey,
	}, nil
}

//Router returns the basic router router for the auth endpoints
func (cba *CorkboardAuth) Router() *httprouter.Router {
	router := httprouter.New()
	router.POST("/api/register", cba.RegisterUser())
	router.POST("/api/authenticate", cba.AuthUser())
	router.GET("/api/public_key", cba.PublicKey())
	return router
}

func getPrivate(location string) (*rsa.PrivateKey, error) {
	privFile, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}
	privPem, _ := pem.Decode(privFile)
	if privPem == nil {
		return nil, errors.New("Could not get the private key from file")
	}
	if privPem.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("Not sure what kind of block this is")
	}
	return x509.ParsePKCS1PrivateKey(privPem.Bytes)
}
