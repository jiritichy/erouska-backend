package functions

import (
	"context"

	"github.com/covid19cz/erouska-backend/internal/functions/checkattemptsthresholds"
	"github.com/covid19cz/erouska-backend/internal/functions/coviddata"
	"github.com/covid19cz/erouska-backend/internal/functions/increaseehridattemptscount"
	"github.com/covid19cz/erouska-backend/internal/functions/increasenotificationscounter"
	"github.com/covid19cz/erouska-backend/internal/functions/isehridactive"
	"github.com/covid19cz/erouska-backend/internal/functions/provideverificationcode"
	"github.com/covid19cz/erouska-backend/internal/functions/registerehrid"
	"github.com/covid19cz/erouska-backend/internal/functions/registernotification"

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

// IncreaseEhridAttemptsCount IncreaseEhridAttemptsCount handler.
func IncreaseEhridAttemptsCount(w http.ResponseWriter, r *http.Request) {
	increaseehridattemptscount.IncreaseEhridAttemptsCount(w, r)
}

// CheckAttemptsThresholds CheckAttemptsThresholds handler.
func CheckAttemptsThresholds(w http.ResponseWriter, r *http.Request) {
	checkattemptsthresholds.CheckAttemptsThresholds(w, r)
}

// IncreaseNotificationsCounter IncreaseNotificationsCounter handler.
func IncreaseNotificationsCounter(w http.ResponseWriter, r *http.Request) {
	increasenotificationscounter.IncreaseNotificationsCounter(w, r)
}

// RegisterNotification RegisterNotification handler.
func RegisterNotification(w http.ResponseWriter, r *http.Request) {
	registernotification.RegisterNotification(w, r)
}

// ProvideVerificationCode ProvideVerificationCode handler.
func ProvideVerificationCode(w http.ResponseWriter, r *http.Request) {
	provideverificationcode.ProvideVerificationCode(w, r)
}

// DownloadCovidDataTotal handler.
func DownloadCovidDataTotal(w http.ResponseWriter, r *http.Request) {
	coviddata.DownloadCovidDataTotal(w, r)
}

// CalculateCovidDataIncrease handler.
func CalculateCovidDataIncrease(ctx context.Context, e coviddata.FirestoreEvent) error {

	return coviddata.CalculateCovidDataIncrease(ctx, e)
}
