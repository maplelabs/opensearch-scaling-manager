package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	mrand "math/rand"
	"os"
	"scaling_manager/logger"
	"time"
)

var log = new(logger.LOG)
var EncryptionSecret string
var seed = time.Now().Unix()

// Initializing logger module
func init() {
	log.Init("logger")
	log.Info.Println("Crypto module initiated")
	mrand.Seed(seed)
	err := CheckAndUpdateSecretFile("secret.txt")
	if err != nil {
		panic(err)
	}

}

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func GeneratePassword() string {
	mrand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "*@#$"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials
	length := 16
	buf := make([]byte, length)
	buf[0] = digits[mrand.Intn(len(digits))]
	buf[1] = specials[mrand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[mrand.Intn(len(all))]
	}
	mrand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	str := string(buf)
	return str
}

func CheckAndUpdateSecretFile(secret_filepath string) error {
	if _, err := os.Stat(secret_filepath); err == nil {
		data, err := os.ReadFile(secret_filepath)
		if err != nil {
			log.Panic.Println("Error reading the secret file")
			return err
		}
		EncryptionSecret = string(data)
	} else {
		EncryptionSecret = GeneratePassword()
		f, err := os.Create(secret_filepath)
		if err != nil {
			log.Panic.Println("Error while creating secret file : ", err)
			return err
		}
		defer f.Close()
		_, err = f.WriteString(EncryptionSecret)
		if err != nil {
			log.Panic.Println("Error while creating secret file : ", err)
			return err
		}
	}
	return nil
}

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
