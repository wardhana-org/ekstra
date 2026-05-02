package auth

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argon2Version = argon2.Version

	maxPasswordBytes = 1024

	minMemory      uint32 = 19 * 1024
	maxMemory      uint32 = 128 * 1024
	maxIterations  uint32 = 10
	maxParallelism uint8  = 4
)

type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var DefaultArgon2Params = Argon2Params{
	Memory:      19 * 1024,
	Iterations:  2,
	Parallelism: 1,
	SaltLength:  16,
	KeyLength:   32,
}

func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("password is required")
	}

	if len(password) > maxPasswordBytes {
		return "", errors.New("password is too long")
	}

	salt, err := randomBytes(DefaultArgon2Params.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		DefaultArgon2Params.Iterations,
		DefaultArgon2Params.Memory,
		DefaultArgon2Params.Parallelism,
		DefaultArgon2Params.KeyLength,
	)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2Version,
		DefaultArgon2Params.Memory,
		DefaultArgon2Params.Iterations,
		DefaultArgon2Params.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func VerifyPassword(encodedHash string, password string) error {
	if len(password) == 0 {
		return errors.New("password is required")
	}

	if len(password) > maxPasswordBytes {
		return errors.New("password is too long")
	}

	params, salt, expectedHash, err := decodeHash(encodedHash)
	if err != nil {
		return err
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	if subtle.ConstantTimeCompare(actualHash, expectedHash) != 1 {
		return errors.New("password does not match")
	}

	return nil
}

func PasswordNeedsRehash(encodedHash string) (bool, error) {
	params, _, _, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	return params.Memory != DefaultArgon2Params.Memory ||
		params.Iterations != DefaultArgon2Params.Iterations ||
		params.Parallelism != DefaultArgon2Params.Parallelism ||
		params.SaltLength != DefaultArgon2Params.SaltLength ||
		params.KeyLength != DefaultArgon2Params.KeyLength, nil
}

func decodeHash(encodedHash string) (*Argon2Params, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[0] != "" {
		return nil, nil, nil, errors.New("invalid hash format")
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, errors.New("unsupported password hash algorithm")
	}

	versionPart := strings.TrimPrefix(parts[2], "v=")
	version, err := strconv.Atoi(versionPart)
	if err != nil {
		return nil, nil, nil, errors.New("invalid argon2 version")
	}

	if version != argon2Version {
		return nil, nil, nil, errors.New("unsupported argon2 version")
	}

	params, err := parseParams(parts[3])
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, errors.New("invalid salt encoding")
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, errors.New("invalid hash encoding")
	}

	params.SaltLength = uint32(len(salt))
	params.KeyLength = uint32(len(hash))

	if params.SaltLength < DefaultArgon2Params.SaltLength {
		return nil, nil, nil, errors.New("salt is too short")
	}

	if params.KeyLength < DefaultArgon2Params.KeyLength {
		return nil, nil, nil, errors.New("hash is too short")
	}

	return params, salt, hash, nil
}

func parseParams(paramString string) (*Argon2Params, error) {
	params := &Argon2Params{}

	for _, item := range strings.Split(paramString, ",") {
		keyValue := strings.SplitN(item, "=", 2)
		if len(keyValue) != 2 {
			return nil, errors.New("invalid argon2 parameter format")
		}

		switch keyValue[0] {
		case "m":
			value, err := strconv.ParseUint(keyValue[1], 10, 32)
			if err != nil {
				return nil, errors.New("invalid argon2 memory parameter")
			}
			params.Memory = uint32(value)

		case "t":
			value, err := strconv.ParseUint(keyValue[1], 10, 32)
			if err != nil {
				return nil, errors.New("invalid argon2 iteration parameter")
			}
			params.Iterations = uint32(value)

		case "p":
			value, err := strconv.ParseUint(keyValue[1], 10, 8)
			if err != nil {
				return nil, errors.New("invalid argon2 parallelism parameter")
			}
			params.Parallelism = uint8(value)

		default:
			return nil, errors.New("unknown argon2 parameter")
		}
	}

	if params.Memory < minMemory || params.Memory > maxMemory {
		return nil, errors.New("argon2 memory parameter outside allowed range")
	}

	if params.Iterations == 0 || params.Iterations > maxIterations {
		return nil, errors.New("argon2 iteration parameter outside allowed range")
	}

	if params.Parallelism == 0 || params.Parallelism > maxParallelism {
		return nil, errors.New("argon2 parallelism parameter outside allowed range")
	}

	return params, nil
}
