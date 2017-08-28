package main

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/alexedwards/scs/engine/cookiestore"
	"github.com/alexedwards/scs/session"
	"github.com/mkozhukh/dsfs/config"
	log "github.com/sirupsen/logrus"
)

func sessionMiddleware() func(h http.Handler) http.Handler {
	// HMAC authentication key (hexadecimal representation of 32 random bytes)
	var hmacKey = []byte(fmt.Sprintf("%x", md5.Sum([]byte(config.Config.Server.Secret))))
	// AES encryption key (hexadecimal representation of 16 random bytes)
	var blockKey = hmacKey

	// Create a new keyset using your authentication and encryption secret keys.
	keyset, err := cookiestore.NewKeyset(hmacKey, blockKey)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new CookieStore instance using the keyset.
	engine := cookiestore.New(keyset)

	return session.Manage(engine)
}
