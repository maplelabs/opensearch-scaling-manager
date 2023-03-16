package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base32"
	ansibleutils "github.com/maplelabs/opensearch-scaling-manager/ansible_scripts"
	"github.com/maplelabs/opensearch-scaling-manager/config"
	"github.com/maplelabs/opensearch-scaling-manager/logger"
	osutils "github.com/maplelabs/opensearch-scaling-manager/opensearchUtils"
	utils "github.com/maplelabs/opensearch-scaling-manager/utilities"
	mrand "math/rand"
	"os"
	"strings"
	"time"
)

var log = new(logger.LOG)
var EncryptionSecret string
var seed = time.Now().Unix()
var SecretFilepath = ".secret.txt"

// Input:
//
// Description:
//
//	     Initializes Crypto module
//		Reads the config.yaml file, if it is a fresh install of scaling manager, the SecretFilepath won't be present.
//		If present, it uses the file to decrypt the credentials present in config.yaml
//		It calls DecryptCredsAndInitializeOs to Initialize the Opensearch client with one try
//		Then Updates the Secret file, and ecrypted credentials into config.yaml
//
// Return:
func init() {
	log.Init("logger")
	log.Info.Println("Crypto module initiated")
	mrand.Seed(seed)
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

	DecryptCredsAndInitializeOs(1)

	UpdateSecretAndEncryptCreds(true, configStruct)
}

// bytes is used when creating ciphers for the string
var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

// Input:
//
// Description:
// 	Generate a random string of length 16
//
// Return:
//	(string): Returns the random string generated as password

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

// Input:
//
// Description:
//
//	Scrambles the EncryptionSecret to store in SecretFilepath
//
// Return:
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

// Input:
//
// Description:
//
//	Read the EncryptionSecret from SecretFilepath, decode the scrmabled secret to original secret and return it
//
// Return:
//
//	(string): Returns the EncryptionSecret
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

// Input:
//
//	osCred (*config.OsCredentials): Pointer to the Opensearch credentials from the config struct
//
// Description:
//
//	Encrypt the Opensearch Credentials passed as a pointer variable
//
// Return:
//
//	(error): Error if any while Encrypting
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

// Input:
//
//	cloudCred (*config.CloudCredentials): Pointer to the Cloud credentials from the config struct
//
// Description:
//
//	Encrypt the Cloud Credentials passed as a pointer variable
//
// Return:
//
//	(error): Error if any while Encrypting
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

	cloudCred.RoleArn, err = GetEncryptedData(cloudCred.RoleArn)
	if err != nil {
		return err
	}

	return nil
}

// Input:
//
//	osCred (*config.OsCredentials): Pointer to the Opensearch credentials from the config struct
//
// Description:
//
//	Decrypt the Opensearch Credentials passed as a pointer variable
//
// Return:
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

// Input:
//
//	cloudCred (*config.CloudCredentials): Pointer to the Cloud credentials from the config struct
//
// Description:
//
//	Decrypt the Cloud Credentials passed as a pointer variable
//
// Return:
func GetDecryptedCloudCreds(cloudCred *config.CloudCredentials) {

	secret_key := GetDecryptedData(cloudCred.SecretKey)
	if secret_key != "" {
		cloudCred.SecretKey = secret_key
	}

	access_key := GetDecryptedData(cloudCred.AccessKey)
	if access_key != "" {
		cloudCred.AccessKey = access_key
	}

	role_arn := GetDecryptedData(cloudCred.RoleArn)
	if role_arn != "" {
		cloudCred.RoleArn = role_arn
	}

}

// Input:
//
//		initialRun (bool): A bool value which says if this function is called from init() or in between while the application is running
//	     config_struct (config.ConfigStruct): Config structure from the config.yaml
//
// Description:
//
//	Encrypt the credentials and update the configuration file [config.yaml] and also Initialize the Opensearch client with new credentials
//
// Return:
//
//	(error): Error if any while Updating the details
func UpdateEncryptedCred(initialRun bool, config_struct config.ConfigStruct) error {
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
		DecryptCredsAndInitializeOs(60)
	}
	return nil
}

// Input:
//
//	try (int): An integer which says the number of times this function has to be retried on any failure
//
// Description:
//
//	     Read the configuration file, get the decrypted Opensearch credentials and Initialize the Opensearch Client
//		This also retries for the number of times specified in the interval of 10 seconds if there is a failure
//
// Return:
func DecryptCredsAndInitializeOs(try int) {
	configStruct, err := config.GetConfig()
	if err != nil {
		log.Error.Println("Error validating config file", err)
		panic(err)
	}
	GetDecryptedOsCreds(&configStruct.ClusterDetails.OsCredentials)

	osErr := osutils.InitializeOsClient(configStruct.ClusterDetails.OsCredentials.OsAdminUsername, configStruct.ClusterDetails.OsCredentials.OsAdminPassword)
	if osErr != nil && try > 0 {
		log.Error.Println("Retrying #", try, " on error: ", osErr)
		time.Sleep(time.Duration(10) * time.Second)
		DecryptCredsAndInitializeOs(try - 1)
	} else if osErr != nil {
		log.Error.Println("Unable to connect to Opensearch even after maximum retries!!")
		log.Error.Println("Please check the Opensearch service or recheck username/password and restart scaling manager...")
		panic(err)
	}

}

// Input:
//
//	initial_run (bool): A bool value which says if this function is called from init() or in between while the application is running
//	config_struct (config.ConfigStruct): Config structure from the config.yaml
//
// Description:
//
//	If it is called by init(), i.e., when initial_run = true, it checks if the current node is master, if yes,
//	It regenerates the secret file, updates the encrypted credentials
//	And populates the secret file and updated config file across all the other nodes in the clusteer using ansible script
//	If not called from init(), it will be called if it is a master node. In this case, it generated the new secret file, updates the encrypted credentials.
//	Populates this information to all the other nodes in the cluster.
//
// Return:
//
//	(error): Error if any
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

// Input:
//
//	currOsCred (config.OsCredentials): Current Opensearch credentials in the config.yaml file
//	prevOsCred (config.OsCredentials): Previous Opensearch credentials
//
// Description:
//
//	Checks if there is a mismatch in the previous and current credentials
//
// Return:
//
//	(bool): Bool value which says if there was a mismatch or not
func OsCredsMismatch(currOsCred config.OsCredentials, prevOsCred config.OsCredentials) bool {
	if (currOsCred.OsAdminUsername != prevOsCred.OsAdminUsername) || (currOsCred.OsAdminPassword != prevOsCred.OsAdminPassword) {
		return true
	}
	return false
}

// Input:
//
//	currCloudCred (config.CloudCredentials): Current Cloud credentials in the config.yaml file
//	prevCloudCred (config.CloudCredentials): Previous Cloud credentials
//
// Description:
//
//	Checks if there is a mismatch in the previous and current credentials
//
// Return:
//
//	(bool): Bool value which says if there was a mismatch or not
func CloudCredsMismatch(currCloudCred config.CloudCredentials, prevCloudCred config.CloudCredentials) bool {
	if (currCloudCred.SecretKey != prevCloudCred.SecretKey) || (currCloudCred.AccessKey != prevCloudCred.AccessKey) || (currCloudCred.RoleArn != prevCloudCred.RoleArn) {
		return true
	}
	return false
}

// Input:
//
//	b ([]byte): Byte value to be encoded to string
//
// Description:
//
//	Encode the given byte value into string
//
// Return:
//
//	(string): Encoded value of byte as string
func Encode(b []byte) string {
	return base32.StdEncoding.EncodeToString(b)
}

// Input:
//
//	     text (string): String to be encrypted
//		EncryptionSecret (string): EncryptionSecret to be used to encrypt the text
//
// Description:
//
//	To encrypt or hide any classified text using the EncryptionSecret
//
// Return:
//
//	     (string): Encrypted value of the text
//		(error): Error if any
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

// Input:
//
//	s (string): String to be decoded
//
// Description:
//
//	Decode the given string
//
// Return:
//
//	([]byte): Decoded value as byte
//	(error): Error if any
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

// Input:
//
//	     text (string): String to be decrypted
//		EncryptionSecret (string): EncryptionSecret to be used for decryption
//
// Description:
//
//	Decrypt method is to extract back the encrypted text
//
// Return:
//
//	([]byte): Decrypted text
//	(error): Error if any
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

// Input :
//
//	toBeEncrypted (string): string to be encrypted
//
// Description :
//
//	Creates an encrypted string : performs AES encryption using the defined secret
//	and return base32 encoded string. Also checks if the encrypted string is able
//	to be decrypted used the same secret.
//
// Output :
//
//	(string, error) : Encrypted string and error if any
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

// Input :
//
//	encryptedString (string): string to be decrypted
//
// Description :
//
//	Return the decrypted string of the given encrypted string
//
// Output :
//
//	(string) : Decrypted string
func GetDecryptedData(encryptedString string) string {
	decrypted_txt, err := Decrypt(encryptedString, EncryptionSecret)
	if err != nil {
		log.Panic.Println("Error decrypting your encrypted text: ", err)
		panic(err)
	}
	return decrypted_txt
}

// Input :
//
//	str (string): string to be converted to matrix
//
// Description :
//
//	Converts a 16 len string to 4*4 matrix
//
// Output :
//
//	([4][4]string) : A matrix derived from the given string
func stringToMatrix(str string) [4][4]string {
	var matrix [4][4]string
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			matrix[i][j] = string(str[i*4+j])
		}
	}
	return matrix
}

// Input :
//
//	matrix ([4][4]string): Matrix to reverse
//
// Description :
//
//	Returns the transpose of the given matrix
//
// Output :
//
//	([4][4]string) : Updated matrix
func transpose(matrix [4][4]string) [4][4]string {
	var transposedMatrix [4][4]string
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			transposedMatrix[j][i] = matrix[i][j]
		}
	}
	return transposedMatrix
}

// Input :
//
//	matrix ([4][4]string): Matrix to reverse
//
// Description :
//
//	Returns the matrix with interchanged rows
//
// Output :
//
//	([4][4]string) : Updated matrix
func reverse(matrix [4][4]string) [4][4]string {
	for i, j := 0, len(matrix)-1; i < j; i, j = i+1, j-1 {
		matrix[i], matrix[j] = matrix[j], matrix[i]
	}
	return matrix
}

// Input :
//
//	matrix ([4][4]string): Matrix to reverse
//
// Description :
//
//	Returns the matrix with intergchanged diagonal values
//
// Output :
//
//	([4][4]string) : Reversed matrix
func reverse_diag(matrix [4][4]string) [4][4]string {
	for i := 0; i < 4; i++ {
		temp := matrix[i][i]
		matrix[i][i] = matrix[i][4-i-1]
		matrix[i][4-i-1] = temp
	}
	return matrix
}

// Input :
//
//	secret (string) : The string which needs to be scrambled or unscrambled
//	scrambled (boolean) : True for scramble, false for unscramble
//
// Description :
//
//	This function scrambles and unscrambles the given string by converting it
//	into matrix and interchanging the values in it.
//
// Output :
//
//	string : scrambled or unscrambled string as per the requirement
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
