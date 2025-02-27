package oidc_config

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"gopkg.in/square/go-jose.v2"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/openshift/rosa/pkg/helper"
)

const (
	defaultLengthRandomLabel = 4

	prefixForPrivateKeySecret     = "rosa-private-key"
	defaultPrefixForConfiguration = "oidc"
)

type OidcConfigInput struct {
	BucketName           string
	IssuerUrl            string
	PrivateKey           []byte
	PrivateKeyFilename   string
	DiscoveryDocument    string
	Jwks                 []byte
	PrivateKeySecretName string
}

func BuildOidcConfigInput(userPrefix, region string) (OidcConfigInput, error) {
	bucketName, err := GenerateBucketName(userPrefix)
	if err != nil {
		return OidcConfigInput{}, fmt.Errorf("There was a problem generating bucket name: %s", err)
	}

	privateKeySecretName := fmt.Sprintf("%s-%s", prefixForPrivateKeySecret, bucketName)
	bucketUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com", bucketName, region)
	privateKey, publicKey, err := CreateKeyPair()
	if err != nil {
		return OidcConfigInput{}, fmt.Errorf("There was a problem generating key pair: %s", err)
	}
	privateKeyFilename := fmt.Sprintf("%s.key", privateKeySecretName)
	discoveryDocument := GenerateDiscoveryDocument(bucketUrl)
	jwks, err := BuildJSONWebKeySet(publicKey)
	if err != nil {
		return OidcConfigInput{}, fmt.Errorf("There was a problem generating JSON Web Key Set: %s", err)
	}
	return OidcConfigInput{
		BucketName:           bucketName,
		IssuerUrl:            bucketUrl,
		PrivateKey:           privateKey,
		PrivateKeyFilename:   privateKeyFilename,
		DiscoveryDocument:    discoveryDocument,
		Jwks:                 jwks,
		PrivateKeySecretName: privateKeySecretName,
	}, nil
}

func GenerateBucketName(userPrefix string) (string, error) {
	randomLabel := helper.RandomLabel(defaultLengthRandomLabel)
	bucketName := fmt.Sprintf("%s-%s", defaultPrefixForConfiguration, randomLabel)
	if userPrefix != "" {
		bucketName = fmt.Sprintf("%s-%s", userPrefix, bucketName)
	}
	if !IsValidBucketName(bucketName) {
		return "", fmt.Errorf("The bucket name '%s' is not valid", bucketName)
	}

	return bucketName, nil
}

const (
	bucketNameRegex = "^[a-z][a-z0-9\\-]+[a-z0-9]$"
)

func IsValidBucketName(bucketName string) bool {
	if bucketName[0] == '.' || bucketName[len(bucketName)-1] == '.' {
		return false
	}
	if strings.HasPrefix(bucketName, "xn--") {
		return false
	}
	if strings.HasSuffix(bucketName, "-s3alias") {
		return false
	}
	if match, _ := regexp.MatchString("\\.\\.", bucketName); match {
		return false
	}
	// We don't support buckets with '.' in them
	match, _ := regexp.MatchString(bucketNameRegex, bucketName)
	return match
}

func CreateKeyPair() ([]byte, []byte, error) {
	bitSize := 4096

	// Generate RSA keypair
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to generate private key")
	}
	encodedPrivateKey := pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Generate public key from private keypair
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to generate public key from private")
	}
	encodedPublicKey := pem.EncodeToMemory(&pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubKeyBytes,
	})

	return encodedPrivateKey, encodedPublicKey, nil
}

const (
	discoveryDocumentTemplate = `{
	"issuer": "%s",
	"jwks_uri": "%s/keys.json",
	"response_types_supported": [
		"id_token"
	],
	"subject_types_supported": [
		"public"
	],
	"id_token_signing_alg_values_supported": [
		"RS256"
	],
	"claims_supported": [
		"aud",
		"exp",
		"sub",
		"iat",
		"iss",
		"sub"
	]
}`
)

func GenerateDiscoveryDocument(bucketURL string) string {
	return fmt.Sprintf(discoveryDocumentTemplate, bucketURL, bucketURL)
}

type JSONWebKeySet struct {
	Keys []jose.JSONWebKey `json:"keys"`
}

// BuildJSONWebKeySet builds JSON web key set from the public key
func BuildJSONWebKeySet(publicKeyContent []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKeyContent)
	if block == nil {
		return nil, errors.Errorf("Failed to decode PEM file")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse key content")
	}

	var alg jose.SignatureAlgorithm
	switch publicKey.(type) {
	case *rsa.PublicKey:
		alg = jose.RS256
	default:
		return nil, errors.Errorf("Public key is not of type RSA")
	}

	kid, err := keyIDFromPublicKey(publicKey)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to fetch key ID from public key")
	}

	var keys []jose.JSONWebKey
	keys = append(keys, jose.JSONWebKey{
		Key:       publicKey,
		KeyID:     kid,
		Algorithm: string(alg),
		Use:       "sig",
	})

	keySet, err := json.MarshalIndent(JSONWebKeySet{Keys: keys}, "", "    ")
	if err != nil {
		return nil, errors.Wrapf(err, "JSON encoding of web key set failed")
	}

	return keySet, nil
}

// keyIDFromPublicKey derives a key ID non-reversibly from a public key
// reference: https://github.com/kubernetes/kubernetes/blob/v1.21.0/pkg/serviceaccount/jwt.go#L89-L111
func keyIDFromPublicKey(publicKey interface{}) (string, error) {
	publicKeyDERBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to serialize public key to DER format")
	}

	hasher := crypto.SHA256.New()
	hasher.Write(publicKeyDERBytes)
	publicKeyDERHash := hasher.Sum(nil)

	keyID := base64.RawURLEncoding.EncodeToString(publicKeyDERHash)

	return keyID, nil
}
