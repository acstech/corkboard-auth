package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	corkboardauth "github.com/acstech/corkboard-auth"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	auth, err := corkboardauth.New(&corkboardauth.Config{
		// CBConnection:   os.Getenv("CB_CONNECTION"),
		// CBBucket:       os.Getenv("CB_BUCKET"),
		// CBBucketPass:   os.Getenv("CB_BUCKET_PASS"),
		PrivateRSAFile: os.Getenv("CB_PRIVATE_RSA"),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", auth.Router()))
}
