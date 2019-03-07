package settings

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"github.com/google/logger"
	"io/ioutil"
	"os"
	"sync"
)

var jwtTokenOnce = sync.Once{}
var secretKey *rsa.PrivateKey

func generateKey() {
	logger.Infoln("Generating a new secret key")

	var err error
	secretKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logger.Fatal("Failed to generate secret key:", err)
	}
}

func GetSecretKey() *rsa.PrivateKey {
	jwtTokenOnce.Do(func() {
		if StoreKey {
			if keyBytes, err := ioutil.ReadFile(SecretPath); err == nil {
				logger.Info("Reading secret key from secret.rsa")
				secretKey, err = x509.ParsePKCS1PrivateKey(keyBytes)
				if err != nil {
					logger.Fatal("Failed to read secret key:", err)
				}
			} else {
				if os.IsNotExist(err) {
					logger.Infoln("Storing secret key in", SecretPath)
					generateKey()
					err = ioutil.WriteFile(
						SecretPath, x509.MarshalPKCS1PrivateKey(secretKey), 0644)
					if err != nil {
						logger.Errorln("Failed to save secret key in", SecretPath, err)
					}

				} else {
					logger.Errorf("Couldn't open %s: %s\n", SecretPath, err)
				}
			}
		} else {
			generateKey()
		}
	})
	return secretKey
}
