package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func ValidateTelegramInitData(initData, botToken string) (map[string]string, bool) {
	fmt.Println("[auth] 📥 Raw initData:", initData)
	fmt.Println("[auth] 🤖 Using botToken:", botToken)

	parsed, err := url.ParseQuery(initData)
	if err != nil {
		fmt.Println("[auth] ❌ Failed to parse initData:", err)
		return nil, false
	}

	data := make(map[string]string)
	var checkStrings []string

	for k, v := range parsed {
		if k == "hash" || k == "signature" {
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

	actual := parsed.Get("hash")
	if actual == "" {
		actual = parsed.Get("signature")
	}

	fmt.Println("[auth] 🧾 Provided hash from initData:", actual)
	fmt.Println("[auth] 🔑 Expected hash:", expected)

	return data, hmac.Equal([]byte(expected), []byte(actual))
}
