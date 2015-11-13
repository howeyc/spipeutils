package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/dchest/spipe"
)

var (
	fAddr    = flag.String("t", "", "target socket")
	fKeyFile = flag.String("k", "", "key file name")
)

func main() {
	flag.Parse()
	if *fKeyFile == "" {
		flag.Usage()
		return
	}
	// Read key file.
	key, err := ioutil.ReadFile(*fKeyFile)
	if err != nil {
		log.Fatalf("key file: %s", err)
	}

	// Dial.
	conn, err := spipe.Dial(key, "tcp", *fAddr)
	if err != nil {
		log.Fatalf("Dial: %s", err)
	}
	defer conn.Close()

	errc := make(chan error, 1)

	go func() {
		_, err := io.Copy(conn, os.Stdin)
		errc <- err
	}()
	go func() {
		_, err := io.Copy(os.Stdout, conn)
		errc <- err
	}()
	<-errc
	conn.Close()
}

