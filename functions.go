package functions

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/covid19cz/erouska-backend/internal/functions/changepushtoken"
	"github.com/covid19cz/erouska-backend/internal/functions/coviddata"
	"github.com/covid19cz/erouska-backend/internal/functions/efgs"
	"github.com/covid19cz/erouska-backend/internal/functions/isehridactive"
	"github.com/covid19cz/erouska-backend/internal/functions/metricsapi"
	"github.com/covid19cz/erouska-backend/internal/functions/publishkeys"
	"github.com/covid19cz/erouska-backend/internal/functions/registerehrid"
	"github.com/covid19cz/erouska-backend/internal/functions/registernotification"
	"github.com/covid19cz/erouska-backend/internal/logging"
	"github.com/covid19cz/erouska-backend/internal/pubsub"
	"time"

	"net/http"
)

// RegisterEhrid Registration handler.
func RegisterEhrid(w http.ResponseWriter, r *http.Request) {
	registerehrid.RegisterEhrid(w, r)
}

// IsEhridActive IsEhridActive handler.
func IsEhridActive(w http.ResponseWriter, r *http.Request) {
	isehridactive.IsEhridActive(w, r)
}

// ChangePushToken ChangePushToken handler.
func ChangePushToken(w http.ResponseWriter, r *http.Request) {
	changepushtoken.ChangePushToken(w, r)
}

// RegisterNotification RegisterNotification handler.
func RegisterNotification(w http.ResponseWriter, r *http.Request) {
	registernotification.RegisterNotification(w, r)
}

// RegisterNotificationAfterMath RegisterNotificationAfterMath handler.
func RegisterNotificationAfterMath(ctx context.Context, m pubsub.Message) error {
	return registernotification.AfterMath(ctx, m)
}

// DownloadCovidDataTotal handler.
func DownloadCovidDataTotal(w http.ResponseWriter, r *http.Request) {
	coviddata.DownloadCovidDataTotal(w, r)
}

// GetCovidData handler.
func GetCovidData(w http.ResponseWriter, r *http.Request) {
	coviddata.GetCovidData(w, r)
}

//PrepareNewMetricsVersion handler.
func PrepareNewMetricsVersion(w http.ResponseWriter, r *http.Request) {
	metricsapi.PrepareNewVersion(w, r)
}

//DownloadMetrics handler.
func DownloadMetrics(w http.ResponseWriter, r *http.Request) {
	metricsapi.DownloadMetrics(w, r)
}

//RegisterEhridAfterMath handler.
func RegisterEhridAfterMath(ctx context.Context, m pubsub.Message) error {
	return registerehrid.AfterMath(ctx, m)
}

// ***************
// EFGS functions:
// ***************

// PublishKeys handler.
func PublishKeys(w http.ResponseWriter, r *http.Request) {
	publishkeys.PublishKeys(w, r)
}

//EfgsUploadKeys handler.
func EfgsUploadKeys(w http.ResponseWriter, r *http.Request) {
	efgs.UploadBatch(w, r)
}

// EfgsDownloadKeys downloads EFGS keys batch
func EfgsDownloadKeys(ctx context.Context, m pubsub.Message) error {
	return efgs.DownloadAndSaveKeys(ctx, m)
}

// EfgsDownloadYesterdaysKeys downloads EFGS keys batch from whole yesterday
func EfgsDownloadYesterdaysKeys(w http.ResponseWriter, r *http.Request) {
	efgs.DownloadAndSaveYesterdaysKeys(w, r)
}

func Budicek(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.FromContext(ctx)

	conf := &firebase.Config{
		DatabaseURL: "firebaseDbURL",
	}

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		logger.Fatalf("firebase.NewApp: %v", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		logger.Errorf("error getting Messaging client: %v\n", err)
	}

	// This registration token comes from the client FCM SDKs.
	//registrationToken := "1fb99f06bf7f791bce2ae93f8b26048201501e73a973700528f3c9083c385d52"

	duration, _ := time.ParseDuration("1d")

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"downloadKeyExport": "true",
		},
		Topic: "budicek",
		Android: &messaging.AndroidConfig{
			Priority: "high",
			TTL:      &duration,
		},
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		panic(err)
	}
	// Response is a message ID string.

	logger.Debugf("Response: %+v", response)

	w.Write([]byte("Successfully sent"))
}
