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
		value := v[0]
		data[k] = value
		checkStrings = append(checkStrings, k+"="+value)
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

	// Логируем всё
	fmt.Println("[auth] 📋 Sorted data strings:")
	for _, s := range checkStrings {
		fmt.Println("  ", s)
	}
	fmt.Println("[auth] 🔐 dataCheckString:\n" + dataCheckString)
	fmt.Println("[auth] 🔑 Expected hash:", expected)
	fmt.Println("[auth] 🆚 Provided hash:", actual)

	return data, hmac.Equal([]byte(expected), []byte(actual))
}
