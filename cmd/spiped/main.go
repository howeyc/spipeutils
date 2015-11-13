package main

import (
        "flag"
        "io"
        "io/ioutil"
        "log"
        "net"
        
        "github.com/dchest/spipe"
)
var (
        fAddr    = flag.String("t", "", "target socket")
        sAddr    = flag.String("s", "", "source socket")
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
        
        s, _ := net.Listen("tcp", *sAddr)
        for {
                c, _ := s.Accept()
                
                // Dial.
                conn, err := spipe.Dial(key, "tcp", *fAddr)
                if err != nil {
                        log.Fatalf("Dial: %s", err)
                }
                defer conn.Close()
                errc := make(chan error, 1)
                go func() {
                        _, err := io.Copy(conn, c)
                        errc <- err
                }()
                go func() {
                        _, err := io.Copy(c, conn)
                        errc <- err
                }()
                <-errc
                conn.Close()
        }
}
