package main

import (
	"flag"
	"fmt"

	"github.com/lucifinil-long/stores/utils"
)

var (
	flDecode = flag.Bool("d", false, "decrypt ?")
	flText   = flag.String("s", "", "text string")
	flKey    = flag.String("k", utils.DefaultEncryptKey, "key string")
)

func main() {
	flag.Parse()
	if *flDecode {
		fmt.Println("orginal encrypted text:", *flText, "key:", *flKey)
		fmt.Println("Descrpted result:", utils.RC4Base64Descrypt(*flText, *flKey))
	} else {
		fmt.Println("orginal raw text:", *flText, "key:", *flKey)
		fmt.Println("Encrpted result:", utils.RC4Base64Encrypt(*flText, *flKey))
	}
}
