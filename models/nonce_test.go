package models

import (
    "testing"
    "time"
)

func TestCreateNonce(t *testing.T) {
    n, err := CreateNonce()
    if err != nil {
        t.Fatal("TestCreateNonce error:", err)
    }

    result := &Nonce{}
    err = Db.SelectOne(result, "select * from nonces where nonce=$1", n.Nonce)
    if err != nil {
        t.Error("TestCreateNonce error:", err)
    }
    result.ExpiresAt = n.ExpiresAt
    if *n != *result {
        t.Errorf("TestCreateNonce: got (%+v), wanted (%+v)", *result, *n)
    }
}

func TestNonceValid(t *testing.T) {
    nonce1, err := generate()
    validNonce := &Nonce{
        Id: "TestNonceValid-1",
        ExpiresAt: time.Now().Add(time.Minute),
        Nonce: nonce1,
    }
    err = Db.Insert(validNonce)
    if err != nil {
        t.Fatal("TestNonceValid error:", err)
    }
    result := NonceValid(validNonce.Nonce)
    if result != true {
        t.Errorf("TestNonceValid: got (%t), want (%t)", result, true)
    }

    nonce2, err := generate()
    invalidNonce := &Nonce{
        Id: "TestNonceValid-2",
        ExpiresAt: time.Now().Add(-time.Minute),
        Nonce: nonce2,
    }
    err = Db.Insert(invalidNonce)
    if err != nil {
        t.Fatal("TestNonceValid error:", err)
    }
    result = NonceValid(invalidNonce.Nonce)
    if result != false {
        t.Errorf("TestNonceValid: got (%t), want (%t)", result, false)
    }
}
