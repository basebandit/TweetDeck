package main

import (
	"fmt"
	"log"
	"os"

	"ekraal.org/avatarlysis/cmd/admin/commands"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		if errors.Cause(err) != commands.ErrHelp {
			log.Printf("error: %s", err)
		}
		os.Exit(1)
	}
}

func run() error {
	//Fix this asap..We aonly need two args provided ./avartalysis-admin keygen ~/.avatarlysis
	if len(os.Args) == 2 {
		fmt.Println("provide file path to keygen")
		return errors.New("key generation destination directory missing")
	}

	switch os.Args[1] {
	case "keygen":

		if err := commands.KeyGen(os.Args[2]); err != nil {
			return errors.Wrap(err, "key generation")
		}

	default:
		fmt.Println("keygen: generate a set of private/public key files")
		return commands.ErrHelp
	}
	return nil
}
