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

	pairs := strings.Split(initData, "&")
	data := make(map[string]string)
	var checkStrings []string

	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := kv[0]
		value := kv[1]

		if key == "hash" || key == "signature" {
			continue
		}

		data[key] = value

		decodedValue, err := url.QueryUnescape(value)
		if err != nil {
			fmt.Println("[auth] ⚠️ Failed to decode value for key:", key, err)
			continue
		}

		checkStrings = append(checkStrings, key+"="+decodedValue)
	}

	sort.Strings(checkStrings)
	dataCheckString := strings.Join(checkStrings, "\n")

	secret := sha256.Sum256([]byte(botToken))
	mac := hmac.New(sha256.New, secret[:])
	mac.Write([]byte(dataCheckString))
	expected := hex.EncodeToString(mac.Sum(nil))

	parsedHash := ""
	for _, pair := range pairs {
		if strings.HasPrefix(pair, "hash=") {
			parsedHash = strings.TrimPrefix(pair, "hash=")
		}
		if strings.HasPrefix(pair, "signature=") && parsedHash == "" {
			parsedHash = strings.TrimPrefix(pair, "signature=")
		}
	}

	fmt.Println("[auth] 📋 Sorted data strings:")
	for _, s := range checkStrings {
		fmt.Println("  ", s)
	}
	fmt.Println("[auth] 🔐 dataCheckString:\n" + dataCheckString)
	fmt.Println("[auth] 🧾 Provided hash from initData:", parsedHash)
	fmt.Println("[auth] 🔑 Expected hash:", expected)

	return data, hmac.Equal([]byte(expected), []byte(parsedHash))
}
