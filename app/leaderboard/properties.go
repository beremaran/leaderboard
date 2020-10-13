package leaderboard

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

const DefaultLeaderboardPrefixKey = "USER_RANKING_"

type Properties struct {
	HttpPort              int
	MysqlConnectionString string
	RedisHost             string
	RedisPassword         string
	RedisDB               int
	RedisCluster          bool
	LeaderboardKeyPrefix  string
}

func LoadProperties() (*Properties, error) {
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	p := &Properties{
		HttpPort:              getInteger("HTTP_PORT", 1323),
		MysqlConnectionString: os.Getenv("MYSQL_CONNECTION_STRING"),
		RedisHost:             os.Getenv("REDIS_HOST"),
		RedisPassword:         os.Getenv("REDIS_PASSWORD"),
		RedisDB:               getInteger("REDIS_DB", 0),
		RedisCluster:          getBool("REDIS_CLUSTER", false),
		LeaderboardKeyPrefix:  getOrDefault("LEADERBOARD_KEY_PREFIX", DefaultLeaderboardPrefixKey),
	}

	return p, nil
}

func getOrDefault(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists || len(value) == 0 {
		return defaultValue
	}

	return value
}

func getInteger(key string, defaultValue int) int {
	integerValue, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return defaultValue
	}

	return integerValue
}

func getBool(key string, defaultValue bool) bool {
	strValue := os.Getenv(key)
	if len(strValue) == 0 {
		return defaultValue
	}

	return strings.ToLower(strValue) == "true"
}
