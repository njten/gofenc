package crypto

import (
	"crypto/rand"
	"errors"
	"os"

	"golang.org/x/crypto/argon2"
)

type Argon2Params struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	Salt    []byte
}

func DefaultArgon2Params(salt []byte) Argon2Params {
	return Argon2Params{
		Time:    3,
		Memory:  65536,
		Threads: 4,
		Salt:    salt,
	}
}

func GenerateRandom(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, errors.New("failed to generate random bytes")
	}
	return b, err

}

func DeriveKey(secret string, isKeyfile bool, params Argon2Params) ([]byte, error) {
	var secretBytes []byte

	if isKeyfile {
		data, err := os.ReadFile(secret)
		if err != nil {
			return nil, errors.New("failed to read keyfile" + err.Error())
		}
		secretBytes = data
	} else {
		secretBytes = []byte(secret)
	}

	if len(secretBytes) == 0 {
		return nil, errors.New("secret or keyfile must not be empty")
	}

	key := argon2.IDKey(
		secretBytes,
		params.Salt,
		params.Time,
		params.Memory,
		params.Threads,
		32,
	)

	return key, nil
}
