package env

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"

	root "github.com/danielmesquitta/flight-api"
	"github.com/danielmesquitta/flight-api/internal/domain/errs"
	"github.com/danielmesquitta/flight-api/internal/pkg/validator"
	"github.com/spf13/viper"
)

const defaultEnvFileName = ".env"

type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentProduction  Environment = "production"
	EnvironmentStaging     Environment = "staging"
	EnvironmentTest        Environment = "test"
)

type Env struct {
	v validator.Validator

	Environment             Environment `mapstructure:"ENVIRONMENT"                 validate:"required,oneof=development production staging test"`
	Port                    string      `mapstructure:"PORT"`
	RedisDatabaseURL        string      `mapstructure:"REDIS_DATABASE_URL"          validate:"required"`
	JWTAccessTokenSecretKey string      `mapstructure:"JWT_ACCESS_TOKEN_SECRET_KEY" validate:"required"`
	AmadeusAPIKey           string      `mapstructure:"AMADEUS_API_KEY"             validate:"required"`
	AmadeusAPISecret        string      `mapstructure:"AMADEUS_API_SECRET"          validate:"required"`
	SerpAPIKey              string      `mapstructure:"SERP_API_KEY"                validate:"required"`
	DuffelAPIKey            string      `mapstructure:"DUFFEL_API_KEY"              validate:"required"`
}

func NewEnv(v validator.Validator) *Env {
	e := &Env{
		v: v,
	}

	if err := e.loadEnv(); err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	return e
}

func (e *Env) loadEnv() error {
	envFile, err := e.getEnvFile()
	if err != nil {
		return errs.New(err)
	}

	viper.SetConfigType("env")

	if err := viper.ReadConfig(bytes.NewBuffer(envFile)); err != nil {
		return errs.New(err)
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&e); err != nil {
		return errs.New(err)
	}

	if err := e.validate(); err != nil {
		return errs.New(err)
	}

	return nil
}

func (e *Env) getEnvFile() (envFile []byte, err error) {
	environment := os.Getenv("ENVIRONMENT")

	if environment != "" {
		envFileName := fmt.Sprintf("%s.%s", defaultEnvFileName, environment)
		envFile, err = root.Env.ReadFile(envFileName)
		if err == nil {
			return envFile, nil
		}
	}

	envFile, err = root.Env.ReadFile(defaultEnvFileName)
	if err != nil {
		return nil, errs.New(err)
	}

	return envFile, nil
}

func (e *Env) validate() error {
	if err := e.v.Validate(e); err != nil {
		return err
	}
	if e.Environment == "" {
		e.Environment = EnvironmentDevelopment
	}
	if e.Port == "" {
		e.Port = "8080"
	}
	return nil
}
