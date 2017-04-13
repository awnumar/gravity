package crypto

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestDeriveSecureValues(t *testing.T) {
	masterPassword := []byte("yellow submarine")
	identifier := []byte("yellow submarine")

	masterKey, rootIdentifier := DeriveSecureValues(masterPassword, identifier, map[string]int{"N": 18, "r": 16, "p": 1})

	actualMasterKey, _ := base64.StdEncoding.DecodeString("IQ0m0/Z7Oy/rvm67Pi0nj2Zk8N0u0Ba+t/uyhPVxTF8=")
	actualRootIdentifier, _ := base64.StdEncoding.DecodeString("FIRp7dJQ2RvA7jsQX1DFWxxit6t9ERMyCSloA8iRmU4=")

	if !bytes.Equal(masterKey[:], actualMasterKey) {
		t.Error("Derived master key != actual value")
	}

	if !bytes.Equal(rootIdentifier, actualRootIdentifier) {
		t.Error("Derived root identifier != actual value")
	}
}

func TestDeriveIdentifierN(t *testing.T) {
	rootIdentifier, _ := base64.StdEncoding.DecodeString("FIRp7dJQ2RvA7jsQX1DFWxxit6t9ERMyCSloA8iRmU4=")

	// Test when n = 0
	actualValue, _ := base64.StdEncoding.DecodeString("hZlSQx9fRk9njp7F1NEs+uVowIU/7DcCXZiuW3kAG2g=")
	if !bytes.Equal(DeriveIdentifierN(rootIdentifier, 0), actualValue) {
		t.Error("When n=0, derivedIdentifierN != actualValue")
	}

	// Test when n = 1
	actualValue, _ = base64.StdEncoding.DecodeString("PT3TYsQ23cJaxhi250QKdUMVUGGEoddspT9nAeSFoj0=")
	if !bytes.Equal(DeriveIdentifierN(rootIdentifier, 1), actualValue) {
		t.Error("When n=1, derivedIdentifierN != actualValue")
	}

	// Test when n = 2
	actualValue, _ = base64.StdEncoding.DecodeString("jKNwKnosCmInggyaqFX/OWehVjtWIywjiyTvgCRZ+T8=")
	if !bytes.Equal(DeriveIdentifierN(rootIdentifier, 2), actualValue) {
		t.Error("When n=2, derivedIdentifierN != actualValue")
	}
}
