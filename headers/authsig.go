package headers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// VerifyAuthSignature checks the HMAC-SHA256 signature for X-Auth headers.
//
// secret: the shared secret key (raw bytes)
// tsStr: the X-Auth-Ts header (string, UNIX epoch)
// claimsB64: the X-Auth-Claims header (base64-encoded JSON)
// sig: the X-Auth-Sig header (base64-encoded HMAC)
// maxSkew: allowed clock skew (e.g., 5 * time.Minute)
func VerifyAuthSignature(secret []byte, tsStr, claimsB64, sig string, maxSkew time.Duration) error {
	if tsStr == "" || claimsB64 == "" || sig == "" {
		return errors.New("missing required auth headers")
	}

	// Parse timestamp
	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid X-Auth-Ts: %w", err)
	}
	now := time.Now().Unix()
	if ts > now+int64(maxSkew.Seconds()) || ts < now-int64(maxSkew.Seconds()) {
		return errors.New("X-Auth-Ts out of allowed range")
	}

	// Compute expected signature
	msg := tsStr + "." + claimsB64
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(msg))
	expectedSig := h.Sum(nil)
	expectedSigB64 := base64.RawStdEncoding.EncodeToString(expectedSig)

	// Compare signatures (constant time)
	if !hmac.Equal([]byte(expectedSigB64), []byte(sig)) {
		return errors.New("invalid X-Auth-Sig")
	}
	return nil
}

// ComputeAuthSignature computes the X-Auth-Sig value for given inputs.
func ComputeAuthSignature(secret []byte, tsStr, claimsB64 string) string {
	msg := tsStr + "." + claimsB64
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(msg))
	return base64.RawStdEncoding.EncodeToString(h.Sum(nil))
}

