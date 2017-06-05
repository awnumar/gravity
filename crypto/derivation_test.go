package crypto

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/0xAwn/memguard"
)

func TestDeriveSecureValues(t *testing.T) {
	masterPassword, _ := memguard.NewFromBytes([]byte("yellow submarine"), false)
	identifier, _ := memguard.NewFromBytes([]byte("yellow submarine"), false)

	masterKey, rootIdentifier := DeriveSecureValues(masterPassword, identifier, map[string]int{"N": 18, "r": 16, "p": 1})

	actualMasterKey, _ := base64.StdEncoding.DecodeString("IQ0m0/Z7Oy/rvm67Pi0nj2Zk8N0u0Ba+t/uyhPVxTF8=")
	actualRootIdentifier, _ := base64.StdEncoding.DecodeString("FIRp7dJQ2RvA7jsQX1DFWxxit6t9ERMyCSloA8iRmU4=")

	if !bytes.Equal(masterKey.Buffer, actualMasterKey) {
		t.Error("Derived master key != actual value")
	}

	if !bytes.Equal(rootIdentifier.Buffer, actualRootIdentifier) {
		t.Error("Derived root identifier != actual value")
	}
}

func TestDeriveIdentifierN(t *testing.T) {
	rootIdentifierBytes, _ := base64.StdEncoding.DecodeString("FIRp7dJQ2RvA7jsQX1DFWxxit6t9ERMyCSloA8iRmU4=")
	rootIdentifier, _ := memguard.NewFromBytes(rootIdentifierBytes, false)

	values := []string{
		"pA095wqN05ms+VQVq+BjIowWQcL6NDw9DbcfMrzTYuk=",
		"iJ+nOpssBHjQYEooof4Ka6BtfXgsA3OZRkLUcNQ/u5Y=",
		"msqNW6pT9+EhpPuo76/tObIcFyqkj+w/0raBsja+Q6I="}

	for i, v := range values {
		actualValue, _ := base64.StdEncoding.DecodeString(v)
		if !bytes.Equal(DeriveIdentifierN(rootIdentifier, uint64(i)), actualValue) {
			t.Errorf("When n=%d, derivedIdentifierN != actualValue", i)
		}
	}
}

func TestDeriveMetaIdentifierN(t *testing.T) {
	rootIdentifierBytes, _ := base64.StdEncoding.DecodeString("FIRp7dJQ2RvA7jsQX1DFWxxit6t9ERMyCSloA8iRmU4=")
	rootIdentifier, _ := memguard.NewFromBytes(rootIdentifierBytes, false)

	values := []string{
		"/Om2e4K6GuC8HVsUcNoIAQtxbXRjZU6XVW6MRjrXVwU=",
		"TQkDMuXFyJfkR4dzRitLVS106s+/8GP9FHBtw6X0nHc=",
		"OKmgv/NCwMUm5TbrDNXV+PPGk6XEc1IhWzhSqEMawzQ="}

	for i, v := range values {
		actualValue, _ := base64.StdEncoding.DecodeString(v)
		if !bytes.Equal(DeriveMetaIdentifierN(rootIdentifier, -i-1), actualValue) {
			t.Errorf("When n=%d, derivedMetaIdentifierN != actualValue", i)
		}
	}
}
