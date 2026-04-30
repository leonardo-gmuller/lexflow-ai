package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Environment string

const (
	EnvTest       Environment = "test"
	EnvLocal      Environment = "local"
	EnvHomolog    Environment = "homolog"
	EnvSandbox    Environment = "sandbox"
	EnvProduction Environment = "production"
)

type Config struct {
	Environment Environment `required:"true" envconfig:"ENVIRONMENT"`
	Development bool        `required:"true" envconfig:"DEVELOPMENT"`

	App    App
	Server Server
	Rag    RAGConfig

	CircuitBreaker CircuitBreaker

	// Infra
	Postgres Postgres
	Redis    RedisConfig
	OpenAI   OpenAIConfig
	Worker   WorkerConfig
	Storage  StorageConfig
}

type App struct {
	Name                    string        `required:"true" envconfig:"APP_NAME"`
	ID                      string        `required:"true" envconfig:"APP_ID"`
	GracefulShutdownTimeout time.Duration `required:"true" envconfig:"APP_GRACEFUL_SHUTDOWN_TIMEOUT"`
	Version                 string        `required:"false"`
}

type Server struct {
	SwaggerHost  string        `required:"true" envconfig:"SERVER_SWAGGER_HOST"`
	Address      string        `required:"true" envconfig:"SERVER_ADDRESS"`
	ReadTimeout  time.Duration `required:"true" envconfig:"SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `required:"true" envconfig:"SERVER_WRITE_TIMEOUT"`
}

type OpenAIConfig struct {
	APIKey string `required:"true" envconfig:"OPENAI_API_KEY"`

	Model          string `envconfig:"OPENAI_MODEL"`
	EmbeddingModel string `envconfig:"OPENAI_EMBEDDING_MODEL"`

	Timeout     int     `envconfig:"OPENAI_TIMEOUT_SECONDS"`
	Temperature float32 `envconfig:"OPENAI_TEMPERATURE"` // 0 to 100
	TopP        float32 `envconfig:"OPENAI_TOP_P"`       // 0 to 100
}

type CircuitBreaker struct {
	Timeout                time.Duration `required:"true" envconfig:"CIRCUIT_BREAKER_TIMEOUT"`
	SleepWindow            time.Duration `required:"true" envconfig:"CIRCUIT_BREAKER_SLEEP_WINDOW"`
	MaxConcurrentRequests  int           `required:"true" envconfig:"CIRCUIT_BREAKER_MAX_CONCURRENT_REQUESTS"`
	RequestVolumeThreshold int           `required:"true" envconfig:"CIRCUIT_BREAKER_REQUEST_VOLUME_THRESHOLD"`
	ErrorPercentThreshold  int           `required:"true" envconfig:"CIRCUIT_BREAKER_ERROR_PERCENT_THRESHOLD"`
}

type Postgres struct {
	Host         string `envconfig:"DB_HOST"     default:"localhost"`
	User         string `envconfig:"DB_USER"     default:"postgres"`
	Password     string `envconfig:"DB_PASSWORD" default:"postgres"`
	DatabaseName string `envconfig:"DB_NAME"     default:"chatbot_go_ai"`
	Port         string `envconfig:"DB_PORT"     default:"5432"`
}

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST"     default:"localhost"`
	Port     string `envconfig:"REDIS_PORT"     default:"6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	DB       int    `envconfig:"REDIS_DB"       default:"0"`
}

type WorkerConfig struct {
	Concurrency  int           `envconfig:"WORKER_CONCURRENCY" default:"5"`
	PollInterval time.Duration `envconfig:"WORKER_POLL_INTERVAL" default:"5s"`
}

type StorageConfig struct {
	BasePath string `envconfig:"STORAGE_BASE_PATH" default:"./storage"`
}

type RAGConfig struct {
	ChunkSize    int `envconfig:"RAG_CHUNK_SIZE" default:"500"`
	ChunkOverlap int `envconfig:"RAG_CHUNK_OVERLAP" default:"50"`

	TopK int `envconfig:"RAG_TOP_K" default:"5"`
}

func New(version string) (Config, error) {
	const operation = "Config.New"

	var cfg Config

	if err := godotenv.Load(); err != nil {
		fmt.Println(".env file not found, using system environment variables")
	}

	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("%s -> %w", operation, err)
	}

	cfg.App.Version = version

	return cfg, nil
}
