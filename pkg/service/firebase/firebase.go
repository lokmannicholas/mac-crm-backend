package firebase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"cloud.google.com/go/firestore"
	"dmglab.com/mac-crm/pkg/config"
	f "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/db"
	"firebase.google.com/go/storage"
	"google.golang.org/api/option"
)

type FirebaseApp struct {
	DB        *db.Client
	FireStore *firestore.Client
	Auth      *auth.Client
	Context   context.Context
	Storage   *storage.Client
	TenantID  string
}

var firebaseApp *FirebaseApp

func GetFirebaseService(ctx context.Context) *FirebaseApp {
	if firebaseApp != nil {
		return firebaseApp
	}
	var app *f.App
	// ctx := context.Background()
	configFile := filepath.Join("asset", "firebase.json")

	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		app, err = f.NewApp(ctx, nil)
		if err != nil {
			log.Fatalf("error initializing app : %v", err)
		}
	} else {

		opt := option.WithCredentialsFile(configFile)
		b, err := ioutil.ReadFile(configFile)
		if err != nil {
			panic(err)
		}
		firebaseConfig := make(map[string]string)
		err = json.Unmarshal(b, &firebaseConfig)
		if err != nil {
			panic(err)
		}
		conf := &f.Config{
			ServiceAccountID: firebaseConfig["client_id"],
			ProjectID:        firebaseConfig["project_id"],
			StorageBucket:    fmt.Sprintf("%s.appspot.com", firebaseConfig["project_id"]),
			DatabaseURL:      fmt.Sprintf("https://%s.firebaseio.com/", firebaseConfig["project_id"]),
		}
		app, err = f.NewApp(ctx, conf, opt)
		if err != nil {
			log.Fatalf("error initializing app : %v", err)
		}
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("app.Auth: %v", err)
	}

	dbClient, err := app.Database(ctx)
	if err != nil {
		log.Fatalf("app.database: %v", err)
	}

	fireStoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}
	storageClient, err := app.Storage(ctx)
	if err != nil {
		log.Fatalf("app.Storage: %v", err)
	}

	// fireStoreClient.Close()
	firebaseApp = &FirebaseApp{
		DB:        dbClient,
		FireStore: fireStoreClient,
		Auth:      authClient,
		Storage:   storageClient,
		Context:   ctx,
		TenantID:  config.GetConfig().CompanyID,
	}
	return firebaseApp
}
func (app *FirebaseApp) UpdateContext(ctx context.Context) {
	app.Context = ctx
}
