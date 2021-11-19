package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log"
)

// verify signature using RSA-PSS SHA256
func verify(pub *rsa.PublicKey, msg string, signature []byte) (err error) {
	h := sha256.New()
	h.Write([]byte(msg))
	d := h.Sum(nil)
	if rsa.VerifyPSS(pub, crypto.SHA256, d, signature, nil) == nil {
		fmt.Print("Skkkrrtttt")
	} else {
		fmt.Print("pahhhh :(")
	}
	return err
}

func main() {
	var sig = "d9zWtT05CIU704Coj4hy0SfcS3ekID5Ol+ezZEYf3OQEcftLV+vpQ/94qnpp1T43BLoIdgWecqDXVhKOL3vM8UXq9o2zAXaWijyXhesvh0k+LLn+T8Ln3vpBaI3esqlHDP3/GF966z1ImkfFja+dKUS9lz3uLa0XLG+HCN3ECpmwPP6DEt3Xu2Vy95KiIcpeT7L4kIrp/dtN9rvrpVF4jsfYYGmaO1GGozshtpYOlFtxQj3FaqWkfTdtXaMOm2ZuJkD0TC6sIW+pLUuvPDpGFrm8hRJgKMDws/aEmCNKi5QXKBs6KHKqmCnXk3jO1IR9bQwUT+cuJJ41N8Vzw2SyGA=="
	var pk = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA52xM4P4DTJwA9tqhkJV9bKRemII2K54OANBBvQA24/TIpxLCQhkliDqTcYP/NYgedMcQfJNNOTTuT5ulSWUpuY0gUPtm5/RzwmEau9Ja0bdlqpcKjy0mG5b+nxycV7v6fGvNQGu8f8Ln+Mc/XRMLqAzsfXFAUvmhHaJO9WLHbYHHnbjAsi7/K2EpwQbUOVjn8eEjCzst5MJSX/fLfXuAjJ207h33X38tSkvlv38+Sbd1QyJVkQxGt4NsAjJs5VIURfZQDInjit18+9M7Y4aZUYjEU3rFYB5fPX0PHiIv7DN9/C35l7V5Q8iYFa7mrq68hYe5hLDU3V6w3PoSg/L99wIDAQAB"
	var pkBytes, err = base64.StdEncoding.DecodeString(pk)
	var sigBytes, err1 = base64.StdEncoding.DecodeString(sig)
	var msg = "rt"
	publicKeyInterface, err := x509.ParsePKIXPublicKey(pkBytes)
	if err != nil || err1 != nil {
		log.Println("Could not parse DER encoded public key (encryption key)")
	}
	publicKey, isRSAPublicKey := publicKeyInterface.(*rsa.PublicKey)
	if !isRSAPublicKey {
		log.Println("Public key parsed is not an RSA public key")
	}
	verify(publicKey, msg, sigBytes)
}
