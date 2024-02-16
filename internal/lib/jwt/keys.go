package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
)

func GenerateKeys() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Ошибка при генерации приватного ключа:", err)
		return err
	}

	publicKey := privateKey.PublicKey

	privateKeyPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	publicKeyPem, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		fmt.Println("Ошибка маршалинга публичного ключа:", err)
		return nil
	}

	privateKeyFile, err := os.Create("internal/lib/jwt/private_key.pem")
	if err != nil {
		return err
	}
	defer func(privateKeyFile *os.File) {
		err := privateKeyFile.Close()
		if err != nil {

		}
	}(privateKeyFile)

	publicKeyFile, err := os.Create("internal/lib/jwt/public_key.pem")
	if err != nil {
		return err
	}
	defer func(publicKeyFile *os.File) {
		err := publicKeyFile.Close()
		if err != nil {

		}
	}(publicKeyFile)

	err = pem.Encode(privateKeyFile, privateKeyPem)
	if err != nil {
		return err
	}

	err = pem.Encode(publicKeyFile, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyPem,
	})
	if err != nil {
		return err
	}

	return nil
}

func PublicKeyToRsa() (*rsa.PublicKey, error) {
	pemData, err := os.ReadFile("./public_key.pem")
	if err != nil {
		log.Fatal("Ошибка чтения файла PEM:", err)
	}

	block, _ := pem.Decode(pemData)
	if block == nil {
		log.Fatal("Не удалось извлечь блок PEM")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal("Ошибка парсинга открытого ключа:", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		log.Fatal("Не удалось преобразовать в тип *rsa.PublicKey")
	}

	return rsaPublicKey, err
}

func PrivateKeyToRsa() (*rsa.PrivateKey, error) {
	privateKeyFile := "./private_key.pem"

	keyBytes, err := os.ReadFile(privateKeyFile)
	if err != nil {
		log.Fatalf("Failed to read file with private key %v", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		log.Fatalf("Failed to parse key: %v", err)
	}

	return privateKey, nil
}
