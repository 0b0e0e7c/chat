package svc

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"chatting/component/auth"
	"chatting/service/userService/internal/config"
	"chatting/service/userService/internal/dao"
	"chatting/service/userService/internal/model"
)

type ServiceContext struct {
	Config  config.Config
	UserDao *dao.UserDao
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(
		&model.User{},
		&model.UserProfile{})
	if err != nil {
		logx.Errorf("failed to auto migrate: %v", err)
		return nil
	}

	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     c.CacheRedis.Host,
			Password: c.CacheRedis.Pass,
		},
	)

	err = parseCert(c.CertConfig)
	if err != nil {
		logx.Error(err)
		return nil
	}

	return &ServiceContext{
		Config:  c,
		UserDao: dao.NewUserDao(db, dao.WithCache(redisClient)),
	}
}

func parseCert(c config.KeyPair) error {
	logx.Info(c)

	// Check if cert file exists
	if _, err := os.Stat(c.CertFile); os.IsNotExist(err) {
		logx.Error("cert file not exists")
		return err
	}

	// Check if key file exists
	if _, err := os.Stat(c.KeyFile); os.IsNotExist(err) {
		logx.Error("key file not exists")
		return err
	}

	// Read and parse the public key
	pubKeyPEM, err := os.ReadFile(c.CertFile)
	if err != nil {
		logx.Error(err)
		return err
	}

	block, _ := pem.Decode(pubKeyPEM)
	if block == nil || block.Type != "PUBLIC KEY" {
		logx.Error("failed to decode PEM block containing public key")
		return err
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logx.Error("failed to parse public key: ", err)
		return err
	}

	var ok bool
	auth.PublicKey, ok = pubKey.(ed25519.PublicKey)
	if !ok {
		logx.Error("not an Ed25519 public key")
		return err
	}

	// Read and parse the private key
	privKeyPEM, err := os.ReadFile(c.KeyFile)
	if err != nil {
		logx.Error(err)
		return err
	}

	block, _ = pem.Decode(privKeyPEM)
	if block == nil || block.Type != "PRIVATE KEY" {
		logx.Error("failed to decode PEM block containing private key")
		return err
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		logx.Error("failed to parse private key: ", err)
		return err
	}

	auth.PrivateKey, ok = privKey.(ed25519.PrivateKey)
	if !ok {
		logx.Error("not an Ed25519 private key")
		return err
	}

	fmt.Println("Public and private keys loaded successfully")
	return nil
}
