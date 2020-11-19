package publishkeys

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/covid19cz/erouska-backend/internal/functions/efgs"
	efgsapi "github.com/covid19cz/erouska-backend/internal/functions/efgs/api"
	efgsdatabase "github.com/covid19cz/erouska-backend/internal/functions/efgs/database"
	efgsutils "github.com/covid19cz/erouska-backend/internal/functions/efgs/utils"
	"github.com/covid19cz/erouska-backend/internal/logging"
	"github.com/covid19cz/erouska-backend/internal/utils"
	"github.com/covid19cz/erouska-backend/pkg/api/v1"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	countryOfOrigin              = "CZ"
	defaultTransmissionRiskLevel = 2 // see docs for ExposureKey - "CONFIRMED will lead to TR 2"
)

var defaultVisitedCountries = []string{"AT", "DE", "DK", "ES", "IE", "NL", "PL"} // this could be a constant but we're in fckn Go

//PublishKeys Handler
func PublishKeys(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	logger := logging.FromContext(ctx).Named("PublishKeys")

	var request v1.PublishKeysRequestDevice

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		logger.Errorf("Could not deserialize request from device: %v", err)
		http.Error(w, "Could not deserialize", http.StatusBadRequest)
		return
	}

	if efgsutils.EfgsExtendedLogging {
		logger.Debugf("Handling PublishKeys request: %+v", request)
	}

	var serverRequest = toServerRequest(&request)

	serverResponse, err := passToKeyServer(ctx, serverRequest)
	if err != nil {
		logger.Errorf("Could not obtain response from Key server: %v", err)
		return
	}

	if efgsutils.EfgsExtendedLogging {
		logger.Debugf("Received response from Key server: %+v", serverResponse)
	}

	if serverResponse.Code == "" && serverResponse.ErrorMessage == "" {
		logger.Infof("Successfully uploaded %v keys to Key server (%v keys sent)", serverResponse.InsertedExposures, len(serverRequest.Keys))

		if request.ConsentToFederation {
			logger.Debug("Going to save uploaded keys to EFGS database")

			if err = handleKeysUpload(request); err != nil {
				logger.Errorf("Error while processing keys persistence: %v", err)
			} else {
				logger.Info("Saved uploaded keys to efgs database")
			}
		} else {
			logger.Info("Federation is disabled for this request")
		}
	} else {
		// error has occurred!
		logger.Errorf("Key server has refused the keys; code %v, message '%v'", serverResponse.Code, serverResponse.ErrorMessage)
	}

	sendResponseToClient(logger, w, toDeviceResponse(serverResponse))
}

func handleKeysUpload(request v1.PublishKeysRequestDevice) error {
	visitedCountries := request.VisitedCountries
	if len(visitedCountries) == 0 {
		visitedCountries = defaultVisitedCountries
	}

	dos := extractDSOS(request)

	if dos <= 0 { // one would use MAX function if Go has some...
		dos = 3
	}

	var keys []*efgsapi.DiagnosisKey
	for _, k := range request.Keys {
		diagnosisKey := efgs.ToDiagnosisKey(&k, countryOfOrigin, visitedCountries, dos)
		if diagnosisKey.TransmissionRiskLevel == 0 {
			diagnosisKey.TransmissionRiskLevel = defaultTransmissionRiskLevel
		}
		keys = append(keys, diagnosisKey)
	}

	return efgsdatabase.Database.PersistDiagnosisKeys(keys)
}

func passToKeyServer(ctx context.Context, request *v1.PublishKeysRequestServer) (*v1.PublishKeysResponseServer, error) {
	logger := logging.FromContext(ctx).Named("passToKeyServer")

	blob, err := json.Marshal(request)
	if err != nil {
		logger.Debugf("Could not serialize request for Key server: %v", err)
		return nil, err
	}

	keyServerConfig, err := utils.LoadKeyServerConfig(ctx)
	if err != nil {
		logger.Fatalf("Could not load key server config: %v", err)
		return nil, err
	}

	response, err := http.Post(keyServerConfig.GetURL("v1/publish"), "application/json", bytes.NewBuffer(blob))
	if err != nil {
		logger.Debugf("Could not obtain response from Key server: %v", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if err := response.Body.Close(); err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %v: %v", response.StatusCode, string(body))
	}

	var serverResponse v1.PublishKeysResponseServer

	if err = json.Unmarshal(body, &serverResponse); err != nil {
		logger.Debugf("Could not deserialize response from Key server: %v", err)
		return nil, err
	}

	return &serverResponse, nil
}

func sendResponseToClient(logger *zap.SugaredLogger, w http.ResponseWriter, response *v1.PublishKeysResponseDevice) {
	blob, err := json.Marshal(response)
	if err != nil {
		logger.Warnf("Could not serialize response for device: %v", err)
		return
	}

	if efgsutils.EfgsExtendedLogging {
		logger.Debugf("Sending response to client: %+v", response)
	}

	_, err = w.Write(blob)
	if err != nil {
		logger.Warnf("Could not send response to device: %v", err)
		return
	}
}

func extractDSOS(request v1.PublishKeysRequestDevice) int {
	// We parse the token but we don't care about signature validation.
	token, _ := jwt.Parse(request.VerificationPayload, func(token *jwt.Token) (interface{}, error) {
		return []byte("hello-world"), nil
	})

	// Here we certainly got validation error but we don't care, the validation was already done by Key server.
	// If we got the token too, it's just enough.

	if token == nil {
		return -1
	}

	// Extract DSOS.
	soi := int64(token.Claims.(jwt.MapClaims)["symptomOnsetInterval"].(float64))
	return int((time.Now().Unix() - soi*600) / 86400)
}

func toServerRequest(request *v1.PublishKeysRequestDevice) *v1.PublishKeysRequestServer {
	return &v1.PublishKeysRequestServer{
		Keys:                 request.Keys,
		HealthAuthorityID:    request.HealthAuthorityID,
		VerificationPayload:  request.VerificationPayload,
		HMACKey:              request.HMACKey,
		SymptomOnsetInterval: request.SymptomOnsetInterval,
		Traveler:             request.Traveler,
		RevisionToken:        request.RevisionToken,
		Padding:              request.Padding,
	}
}

func toDeviceResponse(response *v1.PublishKeysResponseServer) *v1.PublishKeysResponseDevice {
	return &v1.PublishKeysResponseDevice{
		RevisionToken:     response.RevisionToken,
		InsertedExposures: response.InsertedExposures,
		ErrorMessage:      response.ErrorMessage,
		Code:              response.Code,
		Padding:           response.Padding,
	}
}
