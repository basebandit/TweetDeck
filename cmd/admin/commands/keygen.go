package commands

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"path/filepath"

	"os"

	"github.com/pkg/errors"
)

//ErrHelp provides context that help was given
var ErrHelp = errors.New("provided help")

//KeyGen creates an X509 private/public key for auth tokens.
func KeyGen(dir string) error {

	//Generate a new private key.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return ErrHelp
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, os.FileMode(0700)); err != nil {
			return errors.Wrap(err, "creating dir")
		}
	}

	//Create a file for the private key information in PEM format.
	privateFile, err := os.Create(filepath.Join(dir, filepath.Base("private.pem")))
	if err != nil {
		return errors.Wrap(err, "creating private file")
	}
	defer privateFile.Close()

	//Construct a PEM block for the private key.
	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		return errors.Wrap(err, "encoding to private file")
	}
	//Write the private key to the private key file.
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return errors.Wrap(err, "marshaling public key")
	}
	//Create a file for the public key information in PEM form.
	publicFile, err := os.Create(filepath.Join(dir, filepath.Base("public.pem")))
	if err != nil {
		return errors.Wrap(err, "creating public file")
	}
	defer privateFile.Close()

	//Construct a PEM block for the public key.
	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	//Write the public key to the private key file.
	if err := pem.Encode(publicFile, &publicBlock); err != nil {
		return errors.Wrap(err, "encoding to public file")
	}

	fmt.Println("private and public key files generated")
	return nil
}
