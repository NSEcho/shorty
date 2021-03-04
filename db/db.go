package db

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

type Config struct {
	Bucket  []byte
	Timeout int
	Db      *bolt.DB
	hashFn  func(string) string
}

type ConfigOption func(*Config)

func NewConfig(opts ...func(*Config)) *Config {
	cfg := Config{}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &cfg
}

func InitDatabase() (*Config, error) {
	cfg := NewConfig(
		WithBucketName("links.db"),
		WithTimeout(2),
		WithHashedFunc(getHashed))
	var err error
	cfg.Db, err = bolt.Open(string(cfg.Bucket), 0600, &bolt.Options{Timeout: time.Duration(cfg.Timeout) * time.Second})
	if err != nil {
		return nil, err
	}

	return cfg, cfg.Db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(cfg.Bucket)
		return err
	})
}

func WithTimeout(timeout int) ConfigOption {
	return func(cfg *Config) {
		cfg.Timeout = timeout
	}
}

func WithBucketName(name string) ConfigOption {
	return func(cfg *Config) {
		cfg.Bucket = []byte(name)
	}
}

func WithHashedFunc(hashFn func(string) string) ConfigOption {
	return func(cfg *Config) {
		cfg.hashFn = hashFn
	}
}

func (cfg *Config) SaveLink(url string) (string, error) {
	var shorted string
	err := cfg.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(cfg.Bucket)
		shorted = cfg.hashFn(url)
		return b.Put([]byte(url), []byte(shorted))
	})

	return shorted, err
}

func (cfg *Config) GetShorted(shortStr string) (string, error) {
	var url string
	err := cfg.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(cfg.Bucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if string(v) == shortStr {
				url = string(k)
				return nil
			}
		}
		return fmt.Errorf("")
	})

	return url, err
}

func getHashed(url string) string {
	byteURL := []byte(url)
	hash := fmt.Sprintf("%x", sha1.Sum(byteURL))

	return hash[0:10]
}
