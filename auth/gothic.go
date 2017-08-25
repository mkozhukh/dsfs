package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/alexedwards/scs/session"
	"github.com/go-chi/chi"
	"github.com/markbates/goth"
)

// SessionName is the key used to access the session store.
const sessionName = "_gothic_session"

/*
BeginAuthHandler will redirect the user to the appropriate authentication end-point
for the requested provider.
*/
func beginAuthHandler(res http.ResponseWriter, req *http.Request) {
	url, err := getAuthURL(res, req)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(res, err)
		return
	}

	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

// SetState sets the state string associated with the given request.
// If no state string is associated with the request, one will be generated.
// This state is sent to the provider and can be retrieved during the
// callback.
var setState = func(req *http.Request) string {
	state := req.URL.Query().Get("state")
	if len(state) > 0 {
		return state
	}

	return "state"

}

// GetState gets the state returned by the provider during the callback.
// This is used to prevent CSRF attacks, see
// http://tools.ietf.org/html/rfc6749#section-10.12
var getState = func(req *http.Request) string {
	return req.URL.Query().Get("state")
}

/*
GetAuthURL starts the authentication process with the requested provided.
It will return a URL that should be used to send users to.

I would recommend using the BeginAuthHandler instead of doing all of these steps
yourself, but that's entirely up to you.
*/
func getAuthURL(res http.ResponseWriter, req *http.Request) (string, error) {
	providerName, err := getProviderName(req)
	if err != nil {
		return "", err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}
	sess, err := provider.BeginAuth(setState(req))
	if err != nil {
		return "", err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}

	err = session.PutString(req, providerName+sessionName, sess.Marshal())

	if err != nil {
		return "", err
	}

	return url, err
}

/*
CompleteUserAuth does what it says on the tin. It completes the authentication
process and fetches all of the basic information about the user from the provider.

See https://github.com/markbates/goth/examples/main.go to see this in action.
*/
var completeUserAuth = func(res http.ResponseWriter, req *http.Request) (goth.User, error) {
	providerName, err := getProviderName(req)
	if err != nil {
		return goth.User{}, err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return goth.User{}, err
	}

	value, err := session.GetString(req, providerName+sessionName)
	if err != nil {
		return goth.User{}, err
	}

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		return goth.User{}, err
	}

	err = validateState(req, sess)
	if err != nil {
		return goth.User{}, err
	}

	user, err := provider.FetchUser(sess)
	if err == nil {
		// user can be found with existing session data
		return user, err
	}

	// get new token and retry fetch
	_, err = sess.Authorize(provider, req.URL.Query())
	if err != nil {
		return goth.User{}, err
	}

	err = session.PutString(req, providerName+sessionName, sess.Marshal())

	if err != nil {
		return goth.User{}, err
	}

	return provider.FetchUser(sess)
}

// validateState ensures that the state token param from the original
// AuthURL matches the one included in the current (callback) request.
func validateState(req *http.Request, sess goth.Session) error {
	rawAuthURL, err := sess.GetAuthURL()
	if err != nil {
		return err
	}

	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		return err
	}

	originalState := authURL.Query().Get("state")
	if originalState != "" && (originalState != req.URL.Query().Get("state")) {
		return errors.New("state token mismatch")
	}
	return nil
}

// Logout invalidates a user session.
func logoutHandler(res http.ResponseWriter, req *http.Request) error {
	err := session.Destroy(res, req)
	if err != nil {
		return errors.New("Could not delete user session ")
	}
	return nil
}

func getProviderName(req *http.Request) (string, error) {
	provider := chi.URLParam(req, "provider")
	if provider != "" {
		return provider, nil
	}

	// if not found then return an empty string with the corresponding error
	return "", errors.New("you must select a provider")
}
