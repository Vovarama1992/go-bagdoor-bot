package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
)

func ValidateTelegramInitData(initData, botToken string) (map[string]string, bool) {
	parsed, err := url.ParseQuery(initData)
	if err != nil {
		return nil, false
	}

	data := make(map[string]string)
	var checkStrings []string

	for k, v := range parsed {
		if k == "hash" {
			continue
		}
		data[k] = v[0]
		checkStrings = append(checkStrings, k+"="+v[0])
	}

	sort.Strings(checkStrings)
	dataCheckString := strings.Join(checkStrings, "\n")

	secret := sha256.Sum256([]byte(botToken))
	mac := hmac.New(sha256.New, secret[:])
	mac.Write([]byte(dataCheckString))
	expected := hex.EncodeToString(mac.Sum(nil))

	return data, expected == parsed.Get("hash")
}
