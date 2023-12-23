package cryptography

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func (c *Crypto) generateRsaKeyPair(size int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, size)

	return privkey, &privkey.PublicKey, err
}

func (c *Crypto) encryptPKCS1RsaPrivateKey(privKey *rsa.PrivateKey) (string, error) {
	privBlock := pem.Block{}
	privBlock.Type = "RSA PRIVATE KEY"
	privBlock.Bytes = x509.MarshalPKCS1PrivateKey(privKey)

	var privateKeyBuffer bytes.Buffer
	privateKeyBufferWriter := bufio.NewWriter(&privateKeyBuffer)
	err := pem.Encode(privateKeyBufferWriter, &privBlock)
	if err != nil {
		return "", err
	}

	privateKey := make([]byte, privateKeyBufferWriter.Buffered())
	privateKeyBufferWriter.Flush()
	privateKeyBuffer.Read(privateKey)

	privateKeyString := base64.StdEncoding.EncodeToString(privateKey)

	return privateKeyString, nil
}

func (c *Crypto) encryptPKCS1RsaPublicKey(pubKey *rsa.PublicKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", err
	}

	pubBlock := pem.Block{}
	pubBlock.Type = "PUBLIC KEY"
	pubBlock.Bytes = publicKeyBytes

	var publicKeyBuffer bytes.Buffer
	publicKeyBufferWriter := bufio.NewWriter(&publicKeyBuffer)
	err = pem.Encode(publicKeyBufferWriter, &pubBlock)
	if err != nil {
		return "", err
	}

	publicKey := make([]byte, publicKeyBufferWriter.Buffered())
	publicKeyBufferWriter.Flush()
	publicKeyBuffer.Read(publicKey)

	publicKeyString := base64.StdEncoding.EncodeToString(publicKey)

	return publicKeyString, nil
}

func (c *Crypto) decryptPKCS1RsaPrivateKey(privateKeyStr string) (privKey *rsa.PrivateKey, err error) {
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKey)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("not a private key")
	}

	privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, nil

}

func (c *Crypto) GenerateRSAKey(bitSize int) (string, string, error) {
	priv, pub, err := c.generateRsaKeyPair(bitSize)
	if err != nil {
		return "", "", err
	}

	// private key
	privStr, err := c.encryptPKCS1RsaPrivateKey(priv)
	if err != nil {
		return "", "", err
	}

	// public key
	pubStr, err := c.encryptPKCS1RsaPublicKey(pub)
	if err != nil {
		return "", "", err
	}

	return privStr, pubStr, nil

}
