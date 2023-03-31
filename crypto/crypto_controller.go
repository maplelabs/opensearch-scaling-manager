package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base32"
	mrand "math/rand"
	"os"
	"strings"
	"time"

	ansibleutils "github.com/maplelabs/opensearch-scaling-manager/ansible_scripts"
	"github.com/maplelabs/opensearch-scaling-manager/config"
	"github.com/maplelabs/opensearch-scaling-manager/logger"
	osutils "github.com/maplelabs/opensearch-scaling-manager/opensearchUtils"
	utils "github.com/maplelabs/opensearch-scaling-manager/utilities"
)

var log = new(logger.LOG)
var EncryptionSecret string
var seed = time.Now().Unix()
var SecretFilepath = ".secret.txt"

func init() {
	log.Init("logger")
	log.Info.Println("Crypto module initiated")

	mrand.Seed(seed)
	configStruct := parseConfig()

	osutils.InitializeOsClient(configStruct.ClusterDetails.OsCredentials.OsAdminUsername, configStruct.ClusterDetails.OsCredentials.OsAdminPassword)
}

func parseConfig() config.ConfigStruct {
	configStruct, err := config.GetConfig()
	if err != nil {
		log.Error.Println("Error validating config file", err)
		panic(err)
	}
	if _, err = os.Stat(SecretFilepath); err == nil {
		EncryptionSecret = GetEncryptionSecret()
		GetDecryptedOsCreds(&configStruct.ClusterDetails.OsCredentials)
		GetDecryptedCloudCreds(&configStruct.ClusterDetails.CloudCredentials)
	}
	return configStruct
}

// Initializing logger module
func Initialize() {
	configStruct := parseConfig()
	UpdateSecretAndEncryptCreds(true, configStruct)
}

// bytes is used when creating ciphers for the string
var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

// Generate a random string of length 16
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

func GenerateAndScrambleSecret() {
	EncryptionSecret = GeneratePassword()
	f, err := os.Create(SecretFilepath)
	if err != nil {
		log.Panic.Println("Error while creating secret file in master node: ", err)
		panic(err)
	}
	defer f.Close()
	scrambled_secret := Encode([]byte(getScrambledOrOriginalSecret(EncryptionSecret, true)))
	_, err = f.WriteString(scrambled_secret)
	if err != nil {
		log.Panic.Println("Error while writing secret in the master node : ", err)
		panic(err)
	}
}

func GetEncryptionSecret() string {
	data, err := os.ReadFile(SecretFilepath)
	if err != nil {
		log.Panic.Println("Error reading the secret file")
		panic(err)
	}
	decoded_data, decodeErr := Decode(string(data))
	if decodeErr != nil {
		log.Error.Println("Error decoding the data: ", decodeErr)
		panic(decodeErr)
	}
	return getScrambledOrOriginalSecret(string(decoded_data), false)
}

func GetEncryptedOsCred(osCred *config.OsCredentials) error {
	var err error

	osCred.OsAdminUsername, err = GetEncryptedData(osCred.OsAdminUsername)
	if err != nil {
		return err
	}

	osCred.OsAdminPassword, err = GetEncryptedData(osCred.OsAdminPassword)
	if err != nil {
		return err
	}

	return nil
}

func GetEncryptedCloudCred(cloudCred *config.CloudCredentials) error {
	var err error

	cloudCred.SecretKey, err = GetEncryptedData(cloudCred.SecretKey)
	if err != nil {
		return err
	}

	cloudCred.AccessKey, err = GetEncryptedData(cloudCred.AccessKey)
	if err != nil {
		return err
	}

	return nil
}

func GetDecryptedOsCreds(osCred *config.OsCredentials) {

	os_admin_username := GetDecryptedData(osCred.OsAdminUsername)
	if os_admin_username != "" {
		osCred.OsAdminUsername = os_admin_username
	}

	os_admin_password := GetDecryptedData(osCred.OsAdminPassword)
	if os_admin_password != "" {
		osCred.OsAdminPassword = os_admin_password
	}

}

func GetDecryptedCloudCreds(cloudCred *config.CloudCredentials) {

	secret_key := GetDecryptedData(cloudCred.SecretKey)
	if secret_key != "" {
		cloudCred.SecretKey = secret_key
	}

	access_key := GetDecryptedData(cloudCred.AccessKey)
	if access_key != "" {
		cloudCred.AccessKey = access_key
	}

}

func UpdateEncryptedCred(initialRun bool, config_struct config.ConfigStruct) error {
	copyCreds := config_struct.ClusterDetails.OsCredentials
	OsCredErr := GetEncryptedOsCred(&config_struct.ClusterDetails.OsCredentials)
	if OsCredErr != nil {
		log.Panic.Println("Error getting the encrypted config struct : ", OsCredErr)
		panic(OsCredErr)
	}

	CloudCredErr := GetEncryptedCloudCred(&config_struct.ClusterDetails.CloudCredentials)
	if CloudCredErr != nil {
		log.Panic.Println("Error getting the encrypted config struct : ", CloudCredErr)
		panic(CloudCredErr)
	}

	err := config.UpdateConfigFile(config_struct)
	if err != nil {
		log.Panic.Println("Error updating the encrypted config struct : ", err)
		panic(err)
	}

	// initialize new os client connection with the updated creds
	if !initialRun {
		osutils.InitializeOsClient(copyCreds.OsAdminUsername, copyCreds.OsAdminPassword)
	}
	return nil
}

func DecryptCredsAndInitializeOs(config_struct config.ConfigStruct) {
	GetDecryptedOsCreds(&config_struct.ClusterDetails.OsCredentials)
	osutils.InitializeOsClient(config_struct.ClusterDetails.OsCredentials.OsAdminUsername, config_struct.ClusterDetails.OsCredentials.OsAdminPassword)
}

func UpdateSecretAndEncryptCreds(initial_run bool, config_struct config.ConfigStruct) error {
	if initial_run {
		if utils.CheckIfMaster(context.Background(), "") {
			GenerateAndScrambleSecret()
			UpdateEncryptedCred(initial_run, config_struct)
			//ansible logic to copy the secret and config
			hostFileName := "broadcast_hosts"
			utils.HostsWithCurrentNodes(hostFileName, config_struct.ClusterDetails)
			err := ansibleutils.UpdateWithTags(config_struct.ClusterDetails.SshUser, hostFileName, []string{"update_secret", "update_config"})
			if err != nil {
				log.Error.Println(err)
				log.Error.Println("Unable to update config.yaml and .secret.txt on the other node")
				panic(err)
			}
		}
	} else {
		GetDecryptedOsCreds(&config_struct.ClusterDetails.OsCredentials)
		GetDecryptedCloudCreds(&config_struct.ClusterDetails.CloudCredentials)
		GenerateAndScrambleSecret()
		UpdateEncryptedCred(initial_run, config_struct)
		//ansible logic to copy the secret and config
		hostFileName := "broadcast_hosts"
		utils.HostsWithCurrentNodes(hostFileName, config_struct.ClusterDetails)
		err := ansibleutils.UpdateWithTags(config_struct.ClusterDetails.SshUser, hostFileName, []string{"update_secret", "update_config"})
		if err != nil {
			log.Error.Println(err)
			log.Error.Println("Unable to update config.yaml and .secret.txt on the other node")
			panic(err)
		}

	}

	return nil
}

func OsCredsMismatch(currOsCred config.OsCredentials, prevOsCred config.OsCredentials) bool {
	if (currOsCred.OsAdminUsername != prevOsCred.OsAdminUsername) || (currOsCred.OsAdminPassword != currOsCred.OsAdminPassword) {
		return true
	}
	return false
}

func CloudCredsMismatch(currCloudCred config.CloudCredentials, prevCloudCred config.CloudCredentials) bool {
	if (currCloudCred.SecretKey != prevCloudCred.SecretKey) || (currCloudCred.AccessKey != prevCloudCred.AccessKey) {
		return true
	}
	return false
}

// Encode the given byte value
func Encode(b []byte) string {
	return base32.StdEncoding.EncodeToString(b)
}

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text, EncryptionSecret string) (string, error) {
	block, err := aes.NewCipher([]byte(EncryptionSecret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

// Decode the given string
func Decode(s string) ([]byte, error) {
	data, err := base32.StdEncoding.DecodeString(s)
	if err != nil {
		if !strings.Contains(err.Error(), "illegal base32 data at input") {
			log.Panic.Println("Error while decoding : ", err)
			panic(err)
		} else {
			return data, err
		}
	}
	return data, nil
}

// Decrypt method is to extract back the encrypted text
func Decrypt(text, EncryptionSecret string) (string, error) {
	block, err := aes.NewCipher([]byte(EncryptionSecret))
	if err != nil {
		log.Error.Println("Error while creating cipher during decryption : ", err)
		return "", err
	}
	cipherText, err := Decode(text)
	if err != nil {
		return "", nil
	}
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

// Creates an encrypted string : performs AES encryption using the defined secret
// and return base32 encoded string. Also checks if the encrypted string is able
// to be decrypted used the same secret.
func GetEncryptedData(toBeEncrypted string) (string, error) {
	encText, err := Encrypt(toBeEncrypted, EncryptionSecret)
	if err != nil {
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
func GetDecryptedData(encryptedString string) string {
	decrypted_txt, err := Decrypt(encryptedString, EncryptionSecret)
	if err != nil {
		log.Panic.Println("Error decrypting your encrypted text: ", err)
		panic(err)
	}
	return decrypted_txt
}

// Converts a 16 len string to 4*4 matrix
func stringToMatrix(str string) [4][4]string {
	var matrix [4][4]string
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			matrix[i][j] = string(str[i*4+j])
		}
	}
	return matrix
}

// Returns the transpose of the given matrix
func transpose(matrix [4][4]string) [4][4]string {
	var transposedMatrix [4][4]string
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			transposedMatrix[j][i] = matrix[i][j]
		}
	}
	return transposedMatrix
}

// Returns the matrix with interchanged rows
func reverse(matrix [4][4]string) [4][4]string {
	for i, j := 0, len(matrix)-1; i < j; i, j = i+1, j-1 {
		matrix[i], matrix[j] = matrix[j], matrix[i]
	}
	return matrix
}

// Returns the matrix with intergchanged diagonal values
func reverse_diag(matrix [4][4]string) [4][4]string {
	for i := 0; i < 4; i++ {
		temp := matrix[i][i]
		matrix[i][i] = matrix[i][4-i-1]
		matrix[i][4-i-1] = temp
	}
	return matrix
}

// Input :
// secret (string) : The string which needs to be scrambled or unscrambled
// scrambled (boolean) : True for scramble, false for unscramble
//
// Description :
// This function scrambles and unscrambles the given string by converting it
// into matrix and interchanging the values in it.
//
// Output :
// string : scrambled or unscrambled string as per the requirement
func getScrambledOrOriginalSecret(secret string, scrambled bool) string {
	var requiredArr []string
	matrix := stringToMatrix(secret)
	if scrambled {
		matrix = reverse_diag(reverse(transpose(matrix)))
	} else {
		matrix = transpose(reverse(reverse_diag(matrix)))
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			requiredArr = append(requiredArr, matrix[i][j])
		}
	}
	return strings.Join(requiredArr, "")
}
