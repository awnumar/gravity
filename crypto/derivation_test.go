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

	values := []string{
		"hZlSQx9fRk9njp7F1NEs+uVowIU/7DcCXZiuW3kAG2g=",
		"PT3TYsQ23cJaxhi250QKdUMVUGGEoddspT9nAeSFoj0=",
		"jKNwKnosCmInggyaqFX/OWehVjtWIywjiyTvgCRZ+T8="}

	for i, v := range values {
		actualValue, _ := base64.StdEncoding.DecodeString(v)
		if !bytes.Equal(DeriveIdentifierN(rootIdentifier, i), actualValue) {
			t.Errorf("When n=%d, derivedIdentifierN != actualValue", i)
		}
	}
}
