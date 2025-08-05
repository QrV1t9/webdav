package users

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strconv"
	"strings"
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func CheckPassword(userPass, pass string) (bool, error) {
	if strings.HasPrefix(userPass, "$argon2id") {
		r, err := comparePasswordAndHash(pass, userPass)
		if err != nil {
			return false, err
		}
		return r, nil
	}
	return userPass == pass, nil
}

func comparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (p *params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")

	if len(vals) != 6 || vals[1] != "argon2id" {
		return nil, nil, nil, ErrInvalidHash
	}

	versionStr := strings.TrimPrefix(vals[2], "v=")
	version, err := strconv.Atoi(versionStr)
	if version != argon2.Version {
		return nil, nil, nil, fmt.Errorf("%w: expected version %d, got %d", ErrIncompatibleVersion, argon2.Version, version)
	}

	p = &params{}
	paramStr := strings.Split(vals[3], ",")
	if len(paramStr) != 3 {
		return nil, nil, nil, fmt.Errorf("%w: expected 3 parameters (m,t,p), got %d: %v", ErrInvalidHash, len(paramStr), vals[3])
	}

	if !strings.HasPrefix(paramStr[0], "m=") {
		return nil, nil, nil, fmt.Errorf("%w: invalid memory parameter format: %s", ErrInvalidHash, paramStr[0])
	}
	memoryStr := strings.TrimPrefix(paramStr[0], "m=")
	memory, err := strconv.ParseUint(memoryStr, 10, 32)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: invalid memory value: %v", ErrInvalidHash, err)
	}
	p.memory = uint32(memory)

	if !strings.HasPrefix(paramStr[1], "t=") {
		return nil, nil, nil, fmt.Errorf("%w: invalid iterations parameter format: %s", ErrInvalidHash, paramStr[1])
	}
	iterationsStr := strings.TrimPrefix(paramStr[1], "t=")
	iterations, err := strconv.ParseUint(iterationsStr, 10, 32)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: invalid iterations value: %v", ErrInvalidHash, err)
	}
	p.iterations = uint32(iterations)

	if !strings.HasPrefix(paramStr[2], "p=") {
		return nil, nil, nil, fmt.Errorf("%w: invalid parallelism parameter format: %s", ErrInvalidHash, paramStr[2])
	}
	parallelismStr := strings.TrimPrefix(paramStr[2], "p=")
	parallelism, err := strconv.ParseUint(parallelismStr, 10, 8)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: invalid parallelism value: %v", ErrInvalidHash, err)
	}
	p.parallelism = uint8(parallelism)

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: invalid salt format: %v", ErrInvalidHash, err)
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: invalid hash format: %v", ErrInvalidHash, err)
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
