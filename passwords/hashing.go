package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type Params struct {
	Time    uint32 // iterations
	Memory  uint32 // KiB
	Threads uint8
	KeyLen  uint32 // bytes
	SaltLen uint32 // bytes
}

// DefaultParams https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
var DefaultParams = Params{
	Time:    2,
	Memory:  19 * 1024, // 19 MiB
	Threads: 1,
	KeyLen:  32,
	SaltLen: 16,
}

// HashPassword generates a random salt and returns (salt, hash) as base64 strings.
func HashPassword(password string, p Params) (saltB64 string, hashB64 string, err error) {
	if password == "" {
		return "", "", errors.New("password must not be empty")
	}
	if p.Time == 0 || p.Memory == 0 || p.Threads == 0 || p.KeyLen == 0 || p.SaltLen == 0 {
		return "", "", errors.New("invalid argon2 parameters")
	}

	salt := make([]byte, p.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", "", fmt.Errorf("read salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, p.Time, p.Memory, p.Threads, p.KeyLen)

	return base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
		nil
}

// VerifyPassword recomputes the hash using stored salt + params and compares in constant time.
func VerifyPassword(password, saltB64, hashB64 string, p Params) (bool, error) {
	salt, err := base64.RawStdEncoding.DecodeString(saltB64)
	if err != nil {
		return false, fmt.Errorf("decode salt: %w", err)
	}
	wantHash, err := base64.RawStdEncoding.DecodeString(hashB64)
	if err != nil {
		return false, fmt.Errorf("decode hash: %w", err)
	}

	gotHash := argon2.IDKey([]byte(password), salt, p.Time, p.Memory, p.Threads, uint32(len(wantHash)))

	ok := subtle.ConstantTimeCompare(gotHash, wantHash) == 1
	return ok, nil
}
