package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"scaling_manager/logger"
)

var log = new(logger.LOG)

// Initializing logger module
func init() {
	log.Init("logger")
	log.Info.Println("FetchMetrics module initiated")
}

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

// This should be in an env file in production
const EncryptionSecret string = "encryptionsecret"

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text, EncryptionSecret string) (string, error) {
	block, err := aes.NewCipher([]byte(EncryptionSecret))
	if err != nil {
		log.Error.Println("Error while creating cipher during encryption : ", err)
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Error.Println("Error while decoding : ", err)
	}
	return data
}

// Decrypt method is to extract back the encrypted text
func Decrypt(text, EncryptionSecret string) (string, error) {
	block, err := aes.NewCipher([]byte(EncryptionSecret))
	if err != nil {
		log.Error.Println("Error while creating cipher during decryption : ", err)
		return "", err
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

// Creates an encrypted string : performs AES encryption using the defined secret
// and return base64 encoded string. Also checks if the encrypted string is able
// to be decrypted used the same secret.
func GetEncryptedData(toBeEncrypted string) (string, error) {
	encText, err := Encrypt(toBeEncrypted, EncryptionSecret)
	if err != nil {
		log.Error.Println("Error encrypting your classified text: ", err)
		return "", err
	} else {
		_, err := Decrypt(encText, EncryptionSecret)
		if err != nil {
			log.Error.Println("Error decrypting your encrypted text: ", err)
			return "", err
		}
	}
	return encText, nil
}

// Return the decrypted string of the given encrypted string
func GetDecryptedData(encryptedString string) (string, error) {
	decrypted_txt, err := Decrypt(encryptedString, EncryptionSecret)
	if err != nil {
		log.Error.Println("Error decrypting your encrypted text: ", err)
		return "", err
	}
	return decrypted_txt, nil
}
