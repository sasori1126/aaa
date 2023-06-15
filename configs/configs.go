package configs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/go-redis/redis/v8"
	sendinblue "github.com/sendinblue/APIv3-go-library/lib"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Config struct {
	AppEnv            string
	AppUrl            string
	ApiUrl            string
	AppKey            string
	AppName           string
	SendInBlueKey     string
	AccessTokenAlive  string
	RefreshTokenAlive string
}

type SlackConfig struct {
	Webhook, Title, Event, Message, Channel string
}

type AwsConfig struct {
	AwsKeyId       string
	AwsSecretKeyId string
	AwsRegion      string
}

type AxisAddress struct {
	Name     string `json:"name"`
	Company  string `json:"company"`
	Street1  string `json:"street1"`
	Street2  string `json:"street2"`
	City     string `json:"city"`
	Zip      string `json:"zip"`
	State    string `json:"state"`
	Country  string `json:"country"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	StreetNo string `json:"street_no"`
}

type GormModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type BamboraConfig struct {
	MerchantId      string
	ProfilePasscode string
	PaymentPasscode string
}

var (
	AppConfig *Config
	Logger    *zap.SugaredLogger
	Sib       *sendinblue.APIClient
	Trans     ut.Translator
)

func InitConfig() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config := new(Config)
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "local" || appEnv == "prod" || appEnv == "production" {
		config.AppEnv = appEnv
	} else {
		return errors.New("failed to get application env")
	}

	appKey := os.Getenv("APP_KEY")
	if appKey != "" {
		config.AppKey = appKey
	} else {
		return errors.New("failed to get application key")
	}

	aToken := os.Getenv("ACCESS_TOKEN_ALIVE")
	if aToken != "" {
		config.AccessTokenAlive = aToken
	} else {
		return errors.New("failed to get token alive duration")
	}

	rToken := os.Getenv("REFRESH_TOKEN_ALIVE")
	if rToken != "" {
		config.RefreshTokenAlive = rToken
	} else {
		return errors.New("failed to get token alive duration")
	}

	sibApiKey := os.Getenv("SENDINBLUE_API_KEY")
	if sibApiKey != "" {
		config.SendInBlueKey = sibApiKey
	} else {
		return errors.New("failed to get send in blue key name")
	}

	appName := os.Getenv("APP_NAME")
	if appKey != "" {
		config.AppName = appName
	} else {
		return errors.New("failed to get application name")
	}

	if err := transInit(); err != nil {
		return errors.New("failed to init locale trans")
	}

	appUrl, ok := os.LookupEnv("APP_URL")
	if !ok {
		return errors.New("app url not set")
	}
	config.AppUrl = appUrl

	apiUrl, ok := os.LookupEnv("API_URL")
	if !ok {
		return errors.New("api url not set")
	}
	config.ApiUrl = apiUrl

	if AppConfig == nil {
		AppConfig = config
	}

	return nil
}

func GetAxisDeliveryAddress(country string, currency *string) AxisAddress {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	usAddress := AxisAddress{
		Name:    "Axis Forestry Inc. USA Distribution Centre",
		Street1: "1355 Pacific Pl",
		Street2: "Unit 109",
		City:    "Ferndale",
		State:   "WA",
		Zip:     "98248",
		Country: "US",
		Phone:   "+1 888 678 2947",
		Email:   "parts@axisforestry.com",
	}

	caAddress := AxisAddress{
		Name:    "Axis Forestry Inc.",
		Street1: "65 Vicars Road",
		City:    "Kamloops",
		State:   "BC",
		Zip:     "V2C 0G3",
		Country: "CA",
		Phone:   "+1 888 678 2947",
		Email:   "parts@axisforestry.com",
	}
	if currency == nil {
		if country == "CA" || country == "Canada" {
			return caAddress
		}
		return usAddress
	}

	cur := *currency
	if cur == "CAD" {
		return caAddress
	}
	return usAddress
}

func GetBamboraConfig(cur string) (*BamboraConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	currencies := map[string]bool{
		"USD": true,
		"CAD": true,
	}

	if !currencies[cur] {
		return nil, errors.New("invalid currency")
	}

	bc := &BamboraConfig{}
	merchantId, ok := os.LookupEnv("BOMBORA_MERCHANT_ID_" + cur)
	if !ok || merchantId == "" {
		return nil, errors.New("merchant id not found")
	}
	bc.MerchantId = merchantId

	profilePassCode, ok := os.LookupEnv("BOMBORA_PROFILE_PASSCODE_" + cur)
	if !ok || profilePassCode == "" {
		return nil, errors.New("passcode id not found")
	}
	bc.ProfilePasscode = profilePassCode

	paymentPassCode, ok := os.LookupEnv("BOMBORA_PAYMENT_PASSCODE_" + cur)
	if !ok || paymentPassCode == "" {
		return nil, errors.New("payment passcode id not found")
	}
	bc.PaymentPasscode = paymentPassCode

	return bc, nil
}

func NewAwsConfig() (*AwsConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	keyId, ok := os.LookupEnv("AWS_ACCESS_KEY_ID")
	if !ok {
		return nil, errors.New("aws key id must be set")
	}

	secretId, ok := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	if !ok {
		return nil, errors.New("aws secret key id must be set")
	}

	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		return nil, errors.New("region must be set")
	}

	return &AwsConfig{
		AwsKeyId:       keyId,
		AwsSecretKeyId: secretId,
		AwsRegion:      region,
	}, nil
}

func NewSib() (*sendinblue.APIClient, error) {
	cfg := sendinblue.NewConfiguration()
	sibApiKey := os.Getenv("SENDINBLUE_API_KEY")
	if sibApiKey == "" {
		return nil, errors.New("failed to get application name")
	}
	cfg.AddDefaultHeader("api-key", sibApiKey)
	cfg.AddDefaultHeader("partner-key", sibApiKey)
	sib := sendinblue.NewAPIClient(cfg)
	return sib, nil
}

func transInit() (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		enT := en.New() // english
		uni := ut.New(enT, enT)

		var o bool
		Trans, o = uni.GetTranslator(Locale)
		if !o {
			return fmt.Errorf("uni.GetTranslator(%s) failed", Locale)
		}
		// register translate
		// Register translator
		switch Locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, Trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, Trans)
		}
	}
	return
}

func PostgresDns() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("POSTGRES_HOST")
	password := os.Getenv("POSTGRES_PASSWORD")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	db := os.Getenv("POSTGRES_DB")
	timeZone := os.Getenv("TIME_ZONE")

	sslMode := "require"
	if os.Getenv("APP_ENV") == "local" {
		sslMode = "disable"
	}

	return "host=" + host + " user=" + user + " password=" + password + " dbname=" + db + " port=" + port + " sslmode=" + sslMode + " TimeZone=" + timeZone
}

func GetSlackConfig(title, event, message string, isBug bool) *SlackConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	webhook, ok := os.LookupEnv("SLACK_WEBHOOK_URL")
	if !ok {
		return nil
	}

	channel := "bug"
	if !isBug {
		channel = "team-chat"
	}
	return &SlackConfig{
		Channel: channel,
		Title:   title,
		Message: message,
		Event:   event,
		Webhook: webhook,
	}
}

func GetShippoToken() (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token, ok := os.LookupEnv("SHIPPO_TOKEN")
	if !ok {
		return "", errors.New("shippo token not found")
	}

	return token, nil
}

func RedisOptions() *redis.Options {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var opts *redis.Options
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "local" {
		log.Println("local setup")
		host := os.Getenv("REDIS_HOST")
		port := os.Getenv("REDIS_PORT")
		password := os.Getenv("REDIS_PASSWORD")
		opts = &redis.Options{
			Addr:     host + ":" + port,
			Password: password,
			DB:       0,
		}
	} else {
		redisUrl := os.Getenv("REDIS_URL")
		parsedOpts, err := redis.ParseURL(redisUrl)
		if err != nil {
			Logger.Error(err)
			panic(err)
		}

		opts = parsedOpts
	}

	return opts
}

type ConvergePayConfig struct {
	APIBaseURL  string
	MerchantID  string
	MerchantPIN string
	UserID      string
}

func GetConvergePayConfig() (*ConvergePayConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file, error: %+v", err)
	}

	apiBaseURL, ok := os.LookupEnv("CONVERGE_PAY_API_BASE_URL")
	if !ok || apiBaseURL == "" {
		return nil, errors.New("failed to find converge pay API base url")
	}
	merchantID, ok := os.LookupEnv("CONVERGE_PAY_MERCHANT_ID")
	if !ok || merchantID == "" {
		return nil, errors.New("failed to find converge pay merchant id")
	}
	merchantPIN, ok := os.LookupEnv("CONVERGE_PAY_MERCHANT_PIN")
	if !ok || merchantPIN == "" {
		return nil, errors.New("failed to find converge pay merchant pin")
	}
	userID, ok := os.LookupEnv("CONVERGE_PAY_USER_ID")
	if !ok || userID == "" {
		return nil, errors.New("failed to find converge pay user id")
	}
	return &ConvergePayConfig{
		APIBaseURL:  apiBaseURL,
		MerchantID:  merchantID,
		MerchantPIN: merchantPIN,
		UserID:      userID,
	}, nil
}
