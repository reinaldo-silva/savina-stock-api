package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	ServerPort string
	JwtSecret  string
}

type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
}

type S3Config struct {
	Region     string
	BucketName string
	AccessKey  string
	SecretKey  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	config := &Config{
		DBHost:     getEnv("DB_HOST", "postgres"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "secret"),
		DBName:     getEnv("DB_NAME", "stock_db"),
		DBPort:     getEnv("DB_PORT", "5432"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		JwtSecret:  getEnv("JWT_SECRET", "12345"),
	}

	return config
}

func LoadCloudinaryConfig() CloudinaryConfig {
	return CloudinaryConfig{
		CloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		APIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		APISecret: os.Getenv("CLOUDINARY_API_SECRET"),
	}
}

func LoadS3Config() S3Config {
	return S3Config{
		Region:     os.Getenv("AWS_REGION"),
		BucketName: os.Getenv("AWS_BUCKET_NAME"),
		AccessKey:  os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey:  os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)

	if exists {
		return value
	}
	return fallback
}
