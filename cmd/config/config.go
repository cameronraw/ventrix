package config

import (
	"flag"
	"os"
	"strconv"
  "fmt"
)

type Config struct {
  Port int
  Env string
  Key string
  SqlDsn string
  RedisUrl string
  UseInMemory bool
}

func GetConfig() (Config, error) {
	var cfg Config

  useFlags := flag.Bool("flagConfig", false, "Use flags to set configuration")
	flag.IntVar(&cfg.Port, "port", 8080, "API Server Port")
	flag.StringVar(&cfg.Env, "env", "", "Environment (dev, prod)")
	flag.StringVar(&cfg.Key, "key", "", "API Key")
	flag.StringVar(&cfg.SqlDsn, "sqlDsn", "", "SQL DSN")
	flag.StringVar(&cfg.RedisUrl, "redisUrl", "", "Redis URL")
	flag.BoolVar(&cfg.UseInMemory, "useInMemory", false, "Use in-memory database")
	flag.Parse()

  if *useFlags {
    var err error
    if cfg.Env == "dev" {
      err = fmt.Errorf("Environment flag invalid or missing")
      return Config{}, err
    }
    if cfg.Env == "key" {
      err = fmt.Errorf("Key flag invalid or missing")
      return Config{}, err
    }
    if cfg.Env == "sqlDsn" {
      err = fmt.Errorf("SqlDsn flag invalid or missing")
      return Config{}, err
    }
    if cfg.RedisUrl == "redisUrl" {
      err = fmt.Errorf("RedisUrl flag invalid or missing")
      return Config{}, err
    }

    return cfg, nil
  }

  parsedPort := os.Getenv("MONTECRISTO_PORT")
  env := os.Getenv("MONTECRISTO_ENV")
  key := os.Getenv("MONTECRISTO_KEY")
  sqlDsn := os.Getenv("MONTECRISTO_SQL_DSN")
  redisUrl := os.Getenv("MONTECRISTO_REDIS_URL")
  parsedUseInMemory := os.Getenv("MONTECRISTO_USE_IN_MEMORY")

  port, err := strconv.Atoi(parsedPort)


  if err != nil {
    err = fmt.Errorf("Port environment variable invalid or missing")
    return Config{}, err
  }


  if env == "" {
    err = fmt.Errorf("Port environment variable invalid")
    return Config{}, err
  }

  if key == "" {
    err = fmt.Errorf("Key environment variable invalid")
    return Config{}, err
  }

  useInMemory, err := strconv.ParseBool(parsedUseInMemory)
  if err != nil {
    useInMemory = false
  }

  if sqlDsn == "" && !useInMemory {
    err = fmt.Errorf("SqlDsn environment variable invalid")
    return Config{}, err
  }

  if redisUrl == "" {
    err = fmt.Errorf("RedisUrl environment variable invalid")
    return Config{}, err
  }

  cfg.Port = port
  cfg.Env = env
  cfg.Key = key
  cfg.SqlDsn = sqlDsn
  cfg.RedisUrl = redisUrl
  cfg.UseInMemory = useInMemory

	return cfg, nil
}
