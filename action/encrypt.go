package action

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/malnick/cryptorious/config"
)

type VaultSet struct {
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	SecureNote string `yaml:"secure_note"`
}

type Vault struct {
	Data map[string]*VaultSet `yaml:"data"`
	Path string
	Dir  string
}

func (vault *Vault) load() error {
	if _, err := os.Stat(vault.Path); err != nil {
		log.Warnf("%s not found, will create new Vault file.", vault.Path)
		return nil
	}
	yamlBytes, err := ioutil.ReadFile(vault.Path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlBytes, &vault.Data)
	if err != nil {
		return err
	}
	return nil
}

func (vault *Vault) writeToVault(vaultkey string) error {
	// Assumes .load() was called before executing.
	newYamlData, err := yaml.Marshal(&vault.Data)
	if err != nil {
		return err
	}
	if _, err := os.Stat(vault.Path); err != nil {
		log.Warnf("%s does not exist, writing new vault file.", vault.Path)
	}
	if err := ioutil.WriteFile(vault.Path, newYamlData, 0644); err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"vaultkey": vaultkey,
	}).Infof("Successfully wrote to %s", vault.Path)
	return nil
}

func Encrypt(vaultkey string, vs *VaultSet, c config.Config) error {
	keydata, err := ioutil.ReadFile(c.KeyPath)
	if err != nil {
		return err
	}
	log.Debug("Using key ", c.KeyPath)

	// Amend the Vault with the new data
	var vault = Vault{
		Data: make(map[string]*VaultSet),
		Path: c.VaultPath,
	}
	if err := vault.load(); err != nil {
		return err
	}

	if _, ok := vault.Data[vaultkey]; !ok {
		log.Warnf("Key not found, adding: %s", vaultkey)
		vault.Data[vaultkey] = vs
	} else {
		return errors.New("Key already found in vault, will not overwrite.")
	}

	if len(vs.Password) > 0 {
		if encoded, err := encryptValue(keydata, vs.Password); err == nil {
			vault.Data[vaultkey].Password = string(encoded)
		} else {
			return err
		}
	}

	if len(vs.SecureNote) > 0 {
		if encoded, err := encryptValue(keydata, vs.SecureNote); err == nil {
			vault.Data[vaultkey].SecureNote = string(encoded)
		} else {
			return err
		}
	}

	if len(vs.Username) > 0 {
		vault.Data[vaultkey].Username = vs.Username
	}

	if err := vault.writeToVault(vaultkey); err != nil {
		return err
	}

	return nil
}

func encryptValue(key []byte, plaintext string) ([]byte, error) {
	log.Debugf("Encrypting plaintext: %s", plaintext)
	var block cipher.Block

	if _, err := aes.NewCipher(key); err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// iv =  initialization vector
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte{}, errors.New("Errors encountered making initilization vector while encrypting plaintext.")
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	log.Debugf("Ciphertext %x", ciphertext)

	return ciphertext, nil
}

func createPublicKeyBlockCipher(pubData []byte) (interface{}, error) {
	// Create block cipher from RSA key
	block, _ := pem.Decode(pubData)
	// Ensure key is PEM encoded
	if block == nil {
		return nil, errors.New(fmt.Sprintf("Bad key data: %s, not PEM encoded", string(pubData)))
	}
	// Ensure this is actually a RSA pub key
	if got, want := block.Type, "RSA PUBLIC KEY"; got != want {
		return nil, errors.New(fmt.Sprintf("Unknown key type %q, want %q", got, want))
	}
	// Lastly, create the public key using the new block
	pubkey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubkey, nil
}
