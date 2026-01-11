package driver

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func NewFirestore() *firestore.Client {
	creds := serviceAccountCredentials{
		Type:                    os.Getenv("FIREBASE_SA_TYPE"),
		ProjectID:               os.Getenv("FIREBASE_SA_PROJECT_ID"),
		PrivateKeyID:            os.Getenv("FIREBASE_SA_PRIVATE_KEY_ID"),
		PrivateKey:              os.Getenv("FIREBASE_SA_PRIVATE_KEY"),
		ClientEmail:             os.Getenv("FIREBASE_SA_CLIENT_EMAIL"),
		ClientID:                os.Getenv("FIREBASE_SA_CLIENT_ID"),
		AuthURI:                 os.Getenv("FIREBASE_SA_AUTH_URI"),
		TokenURI:                os.Getenv("FIREBASE_SA_TOKEN_URI"),
		AuthProviderX509CertURL: os.Getenv("FIREBASE_SA_AUTH_PROVIDER_X509_CERT_URL"),
		ClientX509CertURL:       os.Getenv("FIREBASE_SA_CLIENT_X509_CERT_URL"),
		UniverseDomain:          os.Getenv("FIREBASE_SA_UNIVERSE_DOMAIN"),
	}

	if creds.ProjectID == "" || creds.PrivateKey == "" || creds.ClientEmail == "" {
		log.Fatal("firestore: required environment variables (FIREBASE_SA_PROJECT_ID, FIREBASE_SA_PRIVATE_KEY, FIREBASE_SA_CLIENT_EMAIL) are not set")
	}

	credJSON, err := json.Marshal(creds)
	if err != nil {
		log.Fatalf("firestore: failed to marshal credentials: %v", err)
	}

	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON(credJSON))
	if err != nil {
		log.Fatalf("firestore: failed to initialize firebase app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("firestore: failed to get firestore client: %v", err)
	}

	return client
}
