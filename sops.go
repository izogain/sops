package sops

import (
	"fmt"
	"time"
)

type Metadata struct {
	LastModified              time.Time
	UnencryptedSuffix         string
	MessageAuthenticationCode string
	Version                   string
	KeySources                []KeySource
}

type KeySource struct {
	Name string
	Keys []*MasterKey
}

type MasterKey interface {
	Encrypt(dataKey string) error
	EncryptIfNeeded(dataKey string) error
	Decrypt() (string, error)
	NeedsRotation() bool
	ToString() string
}

type Store interface {
	Load(data, key string) error
	Dump(key string) (string, error)
}

func (m *Metadata) MasterKeyCount() int {
	count := 0
	for _, ks := range m.KeySources {
		count += len(ks.Keys)
	}
	return count
}

func (m *Metadata) RemoveMasterKeys(keys []MasterKey) {
	for _, ks := range m.KeySources {
		for i, k := range ks.Keys {
			for _, k2 := range keys {
				if (*k).ToString() == k2.ToString() {
					ks.Keys = append(ks.Keys[:i], ks.Keys[i+1:]...)
				}
			}
		}
	}
}

func (m *Metadata) UpdateMasterKeys(dataKey string) {
	for _, ks := range m.KeySources {
		for _, k := range ks.Keys {
			err := (*k).EncryptIfNeeded(dataKey)
			if err != nil {
				fmt.Println("[WARNING]: could not encrypt data key with master key ", (*k).ToString())
			}
		}
	}
}