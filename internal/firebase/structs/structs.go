package structs

//Registration DB entity for registration.
type Registration struct {
	Platform                  string `json:"platform"`
	PlatformVersion           string `json:"platformVersion"`
	Manufacturer              string `json:"manufacturer"`
	Model                     string `json:"model"`
	Locale                    string `json:"locale"`
	PushRegistrationToken     string `json:"pushRegistrationToken"`
	CreatedAt                 int64  `json:"createdAt"`
	LastNotificationStatus    string `json:"lastNotificationStatus"`
	LastNotificationUpdatedAt int64  `json:"lastNotificationUpdatedAt"`
}

//RegistrationV1 DB entity for registration V1 - the legacy one.
type RegistrationV1 struct {
	FUID                  string `firestore:"fuid" json:"fuid"`
	Platform              string `firestore:"platform" json:"platform"`
	PlatformVersion       string `firestore:"platformVersion" json:"platformVersion"`
	Manufacturer          string `firestore:"manufacturer" json:"manufacturer"`
	Model                 string `firestore:"model" json:"model"`
	Locale                string `firestore:"locale" json:"locale"`
	PushRegistrationToken string `firestore:"pushRegistrationToken" json:"pushRegistrationToken"`
	CreatedAt             int64  `firestore:"createdAt" json:"createdAt"`
}

//NotificationCounter DB entity for notification counter.
type NotificationCounter struct {
	NotificationsCount int `json:"notificationsCount"`
}

//UserCounter DB entity for users counter.
type UserCounter struct {
	UsersCount int `json:"usersCount"`
}

//PublisherCounter DB entity for publishers counter.
type PublisherCounter struct {
	PublishersCount int `json:"publishersCount"`
	KeysCount       int `json:"keysCount"`
}

//VerificationCodeMetadata DB entity for verification code metadata.
type VerificationCodeMetadata struct {
	VsMetadata map[string]interface{} `json:"vsMetadata"`
	IssuedAt   int64                  `json:"issuedAt"`
}

// IntegerValue represents integer (as string) in firestore events
type IntegerValue struct {
	IntegerValue string `json:"integerValue"`
}

// StringValue represents strings in firestore events
type StringValue struct {
	StringValue string `json:"stringValue"`
}
