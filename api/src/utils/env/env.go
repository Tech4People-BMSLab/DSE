package env

import (
	"dse/src/utils"
	"dse/src/utils/arrays"
	"os"
	"regexp"

	v "github.com/cohesivestack/valgo"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var (
	logger   = utils.NewLogger()
	validate = validator.New()
)

// ------------------------------------------------------------
// : Structs
// ------------------------------------------------------------
type Environment struct {
	ENVIRONMENT string `zog:"ENVIRONMENT"`

	API_HOST    string `zog:"API_HOST"`
	API_PORT    string `zog:"API_PORT"`
	
	DB_NAME string `zog:"DB_NAME"`
	DB_USER string `zog:"DB_USER"`
	DB_PASS string `zog:"DB_PASS"`
	DB_HOST string `zog:"DB_HOST"`
	DB_PORT string `zog:"DB_PORT"`
	DB_SSL  string `zog:"DB_SSL"`
}

// ------------------------------------------------------------
// : Init
// ------------------------------------------------------------
func init() {
	// TODO: TBI Replace the Load func
}

// ------------------------------------------------------------
// : Getters
// ------------------------------------------------------------
func GetEnvironment() string {
	if os.Getenv("ENVIRONMENT") == "" {
		return "production"
	}

	return os.Getenv("ENVIRONMENT")
}

func IsLocal() bool {
	return GetEnvironment() == "local"
}

func IsDevelopment() bool {
	return GetEnvironment() == "development"
}

func IsProduction() bool {
	return GetEnvironment() == "production"
}

func GetAPIHost() string {
	if os.Getenv("API_HOST") == "" {
		return "0.0.0.0"
	}

	return os.Getenv("API_HOST")
}

func GetAPIPort() string {
	if os.Getenv("API_PORT") == "" {
		return "5000"
	}

	return os.Getenv("API_PORT")
}

func GetDBName() string {
	if os.Getenv("DB_NAME") == "" {
		return "dse"
	}

	return os.Getenv("DB_NAME")
}

func GetDBUser() string {
	if os.Getenv("DB_USER") == "" {
		return "postgres"
	}

	return os.Getenv("DB_USER")
}

func GetDBPass() string {
	if os.Getenv("DB_PASS") == "" {
		return "postgres"
	}

	return os.Getenv("DB_PASS")
}

func GetDBHost() string {
	if os.Getenv("DB_HOST") == "" {
		return "localhost"
	}

	return os.Getenv("DB_HOST")
}

func GetDBPort() string {
	if os.Getenv("DB_PORT") == "" {
		return "5432"
	}

	return os.Getenv("DB_PORT")
}

func GetDBSSL() string {
	if os.Getenv("DB_SSL") == "" {
		return "disable"
	}

	if os.Getenv("DB_SSL") == "1" {
		return "enable"
	}

	if os.Getenv("DB_SSL") == "0" {
		return "disable"
	}

	return "disable"
}

// ------------------------------------------------------------
// : Functions
// ------------------------------------------------------------
func Load() {
	godotenv.Load(".env")

	env := Environment{
		ENVIRONMENT: os.Getenv("ENVIRONMENT"),

		API_HOST: os.Getenv("API_HOST"),
		API_PORT: os.Getenv("API_PORT"),

		DB_NAME: os.Getenv("DB_NAME"),
		DB_USER: os.Getenv("DB_USER"),
		DB_PASS: os.Getenv("DB_PASS"),
		DB_HOST: os.Getenv("DB_HOST"),
		DB_PORT: os.Getenv("DB_PORT"),
		DB_SSL:  os.Getenv("DB_SSL"),
	}

	rhost := regexp.MustCompile(`^(localhost|(([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,})|(\d{1,3}\.){3}\d{1,3})$`)
	rport := regexp.MustCompile(`^\d{1,5}$`)

	validation := v.New()
	validation.Is(v.String(env.ENVIRONMENT, "environment").EqualTo("local").Or().EqualTo("production"))
	
	validation.Is(v.String(env.API_HOST, "api_host").MatchingTo(rhost))
	validation.Is(v.String(env.API_PORT, "api_port").MatchingTo(rport))

	validation.Is(v.String(env.DB_NAME, "db_name").Not().Empty())
	validation.Is(v.String(env.DB_USER, "db_user").Not().Empty())
	validation.Is(v.String(env.DB_PASS, "db_pass").Not().Empty())
	validation.Is(v.String(env.DB_HOST, "db_host").MatchingTo(rhost))
	validation.Is(v.String(env.DB_PORT, "db_port").MatchingTo(rport))
	validation.Is(v.String(env.DB_SSL, "db_ssl").EqualTo("0").Or().EqualTo("1"))

	if validation.Valid() == false {
		for _, err := range validation.Errors() {
			field   := err.Name()
			message := arrays.First[string](err.Messages())

			logger.Error().Str("field", field).Str("message", *message).Msg("Failed to validate environment variables")
		}

		logger.Fatal().Msg("Failed to validate environment variables")
	}
}
