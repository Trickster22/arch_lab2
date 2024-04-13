package cryp

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
)

type Data struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
	PublicKey string `json:"publickey"`
}

func GenKeys(name string) (*rsa.PrivateKey, error) {
	fmt.Printf("%s: генерируем пару ключей\n", name)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func SignMessage(message string, privateKey *rsa.PrivateKey) (string, error) {
	hash := sha256.Sum256([]byte(message))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}

	signatureHex := hex.EncodeToString(signature)
	return signatureHex, nil
}

func GetPubKeyPem(privateKey *rsa.PrivateKey) string {
	pubKeyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	pubKeyPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubKeyBytes})
	return string(pubKeyPem)
}

func VerifySignature(data Data) error {
	block, _ := pem.Decode([]byte(data.PublicKey))
	if block == nil {
		return fmt.Errorf("не получилось декодировать публичный ключ")
	}
	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("не получилось распарсить публичный ключ")
	}

	hash := sha256.Sum256([]byte(data.Message))
	singnature, err := hex.DecodeString(data.Signature)
	if err != nil {
		return fmt.Errorf("не получилось декодировать подпись")
	}

	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], singnature)
	if err != nil {
		return err
	}
	return nil
}
