package main

import (
	"net/http"

	"github.com/alexedwards/scs/engine/cookiestore"
	"github.com/alexedwards/scs/session"
	log "github.com/sirupsen/logrus"
)

func sessionMiddleware() func(h http.Handler) http.Handler {
	// HMAC authentication key (hexadecimal representation of 32 random bytes)
	var hmacKey = []byte("f71dc7e58abab014ddad2652475056f185164d262869c8931b239de52711ba87")

	// AES encryption key (hexadecimal representation of 16 random bytes)
	var blockKey = []byte("911182cec2f206986c8c82440adb7d17")

	// Create a new keyset using your authentication and encryption secret keys.
	keyset, err := cookiestore.NewKeyset(hmacKey, blockKey)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new CookieStore instance using the keyset.
	engine := cookiestore.New(keyset)

	return session.Manage(engine)
}
