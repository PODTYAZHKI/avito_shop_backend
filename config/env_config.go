package config

import (
	"fmt"
	"os"
	"regexp"
)

func getEnv(key, def string) string {
    val, ok := os.LookupEnv(key)
    if !ok {
        return def
    }
    return val
}

func DBConfig() string {


	postgresConn := getEnv("POSTGRES_CONN", "postgres://db_pg:db_postgres@db:5432/db_shop?sslmode=disable")
	postgresJDBC := getEnv("POSTGRES_JDBC_URL", "jdbc:postgresql://db:5432/db_shop")
	postgresUser := getEnv("POSTGRES_USERNAME", "db_pg")
	postgresPassword := getEnv("POSTGRES_PASSWORD", "db_postgres")
	postgresHost := getEnv("POSTGRES_HOST", "localhost")
	postgresPort := getEnv("POSTGRES_PORT", "5432")
	postgresDB := getEnv("POSTGRES_DATABASE", "db_shop")

	var config string

	if postgresConn != "" {
		config = postgresConn
		return config
	}
	if postgresJDBC != "" {
		configJDBC, err := jdbcToGormConfig(postgresJDBC)
		if err == nil {
			config = configJDBC
			return config
		}
	}
	config = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		postgresHost, postgresUser, postgresPassword, postgresDB, postgresPort)
	return config
}

func jdbcToGormConfig(jdbcURL string) (string, error) {
	re := regexp.MustCompile(`jdbc:postgresql://([^:]+):(\d+)/([^?]+)\?user=([^&]+)&password=([^&]+)`)
	matches := re.FindStringSubmatch(jdbcURL)
	if len(matches) != 6 {
		return "", fmt.Errorf("не удалось разобрать строку JDBC")
	}

	host := matches[1]
	port := matches[2]
	dbname := matches[3]
	user := matches[4]
	password := matches[5]

	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port), nil
}