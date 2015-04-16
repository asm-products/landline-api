package models

import "testing"

func TestSSORequestIsValid(t *testing.T) {
	secret := "secret"
	message := "a secret message"
	signature := Sign([]byte(secret), []byte(message))

	validRequest := SSORequest{
		Payload:   message,
		Signature: signature,
	}
	if !validRequest.IsValid(secret) {
		t.Errorf("TestSSORequestIsValid: should be valid - %+v", validRequest)
	}

	invalidRequest := SSORequest{
		Payload:   message,
		Signature: "someothersignature",
	}
	if invalidRequest.IsValid(secret) {
		t.Errorf("TestSSORequestIsValid: should be invalid - %+v", invalidRequest)
	}
}
