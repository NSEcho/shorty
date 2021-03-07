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

func (cfg *Config) SaveLink(url string) (string, error) {
	var shorted string
	err := cfg.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(cfg.Bucket)
		shorted = cfg.hashFn(url)
		return b.Put([]byte(shorted), []byte(url))
	})

	return shorted, err
}

func (cfg *Config) GetShorted(shortStr string) string {
	var url string
	_ = cfg.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(cfg.Bucket)
		url = string(b.Get([]byte(shortStr)))
		return nil
	})

	return url
}

type ConfigOption func(*Config)

func NewConfig(opts ...ConfigOption) *Config {
	cfg := Config{
		Bucket:  []byte("links.db"),
		Timeout: 1,
		Db:      nil,
		hashFn:  getHashed,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &cfg
}

/*
InitDatabase will initialize database with functional parameters passed to the function.
Functional parameters need to be in format like below

func example(bucketName string) ConfigOption {
	return func(cfg *Config) {
		cfg.Bucket = []byte(bucketName)
	}
}

*/
func InitDatabase(opts ...ConfigOption) (*Config, error) {
	cfg := NewConfig(opts...)
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

func getHashed(url string) string {
	byteURL := []byte(url)
	hash := fmt.Sprintf("%x", sha1.Sum(byteURL))

	return hash[0:10]
}
