package main

//
//import (
//	"bufio"
//	"crypto/aes"
//	"crypto/cipher"
//	"crypto/ecdh"
//	"crypto/rand"
//	"crypto/sha256"
//	"errors"
//	"fmt"
//	"log"
//	"os"
//)
//
//type DHKeyPair struct {
//	PrivateKey *ecdh.PrivateKey
//	PublicKey  *ecdh.PublicKey
//}
//
//func GenerateDHKeyPair() (*DHKeyPair, error) {
//
//	curve := ecdh.P256()
//
//	privKey, err := curve.GenerateKey(rand.Reader)
//	if err != nil {
//		return nil, fmt.Errorf("failed to generate DH keys: %v", err)
//	}
//
//	return &DHKeyPair{
//		PrivateKey: privKey,
//		PublicKey:  privKey.PublicKey(),
//	}, nil
//}
//
//func ComputeSharedSecret(privKey *ecdh.PrivateKey, peerPubKey *ecdh.PublicKey) ([]byte, error) {
//	secret, err := privKey.ECDH(peerPubKey)
//	if err != nil {
//		return nil, fmt.Errorf("failed to compute shared secret: %v", err)
//	}
//	return secret, nil
//}
//
//func EncryptMessage(message []byte, sharedSecret []byte) ([]byte, []byte, error) {
//	key := sha256.Sum256(sharedSecret)
//
//	block, err := aes.NewCipher(key[:])
//	if err != nil {
//		return nil, nil, fmt.Errorf("failed to create cipher: %v", err)
//	}
//
//	gcm, err := cipher.NewGCM(block)
//	if err != nil {
//		return nil, nil, fmt.Errorf("failed to create GCM: %v", err)
//	}
//
//	nonce := make([]byte, gcm.NonceSize())
//	if _, err := rand.Read(nonce); err != nil {
//		return nil, nil, fmt.Errorf("failed to generate nonce: %v", err)
//	}
//
//	ciphertext := gcm.Seal(nil, nonce, message, nil)
//	return ciphertext, nonce, nil
//}
//
//func DecryptMessage(encryptedMessage []byte, nonce []byte, sharedSecret []byte) ([]byte, error) {
//	key := sha256.Sum256(sharedSecret)
//
//	block, err := aes.NewCipher(key[:])
//	if err != nil {
//		return nil, fmt.Errorf("failed to create cipher: %v", err)
//	}
//
//	gcm, err := cipher.NewGCM(block)
//	if err != nil {
//		return nil, fmt.Errorf("failed to create GCM: %v", err)
//	}
//
//	if len(nonce) != gcm.NonceSize() {
//		return nil, errors.New("incorrect nonce size")
//	}
//
//	plaintext, err := gcm.Open(nil, nonce, encryptedMessage, nil)
//	if err != nil {
//		return nil, fmt.Errorf("failed to decrypt message: %v", err)
//	}
//
//	return plaintext, nil
//}
//
//func main() {
//	aliceKeys, err := GenerateDHKeyPair()
//	if err != nil {
//		log.Fatal("Alice keygen failed:", err)
//	}
//
//	bobKeys, err := GenerateDHKeyPair()
//	if err != nil {
//		log.Fatal("Bob keygen failed:", err)
//	}
//
//	aliceSecret, err := ComputeSharedSecret(aliceKeys.PrivateKey, bobKeys.PublicKey)
//	if err != nil {
//		log.Fatal("Alice failed to compute secret:", err)
//	}
//
//	bobSecret, err := ComputeSharedSecret(bobKeys.PrivateKey, aliceKeys.PublicKey)
//	if err != nil {
//		log.Fatal("Bob failed to compute secret:", err)
//	}
//
//	fmt.Println("соо:")
//	scanner := bufio.NewScanner(os.Stdin)
//	scanner.Scan()
//	message := []byte(scanner.Text())
//	encrypted, nonce, err := EncryptMessage(message, aliceSecret)
//	if err != nil {
//		log.Fatal("Encryption failed:", err)
//	}
//
//	decrypted, err := DecryptMessage(encrypted, nonce, bobSecret)
//	if err != nil {
//		log.Fatal("Decryption failed:", err)
//	}
//
//	fmt.Printf("Original: %s\n", message)
//	fmt.Printf("Encrypted: %x\n", encrypted)
//	fmt.Printf("Decrypted: %s\n", decrypted)
//}
