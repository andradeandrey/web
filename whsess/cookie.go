package whsess

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"net/http"

	"golang.org/x/crypto/nacl/secretbox"
)

const (
	nonceLength = 24
	KeyLength   = 32
)

type CookieOptions struct {
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

type CookieStore struct {
	Options CookieOptions
	Secret  [KeyLength]byte
}

var _ Store = (*CookieStore)(nil)

// NewCookieStore creates a secure cookie store with default settings.
// Configure the Options field further if additional settings are required.
func NewCookieStore(secretKey []byte) *CookieStore {
	rv := &CookieStore{
		Options: CookieOptions{
			Path:   "/",
			MaxAge: 86400 * 30}}
	copy(rv.Secret[:], secretKey)
	return rv
}

// Load implements the Store interface. Not expected to be used directly.
func (cs *CookieStore) Load(r *http.Request, namespace string) (rv SessionData,
	err error) {
	empty := SessionData{New: true, Values: map[interface{}]interface{}{}}
	c, err := r.Cookie(namespace)
	if err != nil || c.Value == "" {
		return empty, nil
	}
	data, err := base64.URLEncoding.DecodeString(c.Value)
	if err != nil {
		return empty, nil
	}
	var nonce [nonceLength]byte
	copy(nonce[:], data[:nonceLength])
	decrypted, ok := secretbox.Open(nil, data[nonceLength:], &nonce,
		&cs.Secret)
	if !ok {
		return empty, nil
	}
	err = gob.NewDecoder(bytes.NewReader(decrypted)).Decode(&rv.Values)
	if err != nil {
		return empty, nil
	}
	return rv, nil
}

// Save implements the Store interface. Not expected to be used directly.
func (cs *CookieStore) Save(w http.ResponseWriter, namespace string,
	s SessionData) error {
	var out bytes.Buffer
	err := gob.NewEncoder(&out).Encode(&s.Values)
	if err != nil {
		return err
	}
	var nonce [nonceLength]byte
	_, err = rand.Read(nonce[:])
	if err != nil {
		return err
	}
	value := base64.URLEncoding.EncodeToString(
		secretbox.Seal(nonce[:], out.Bytes(), &nonce, &cs.Secret))

	return setCookie(w, &http.Cookie{
		Name:     namespace,
		Value:    value,
		Path:     cs.Options.Path,
		Domain:   cs.Options.Domain,
		MaxAge:   cs.Options.MaxAge,
		Secure:   cs.Options.Secure,
		HttpOnly: cs.Options.HttpOnly})
}

func setCookie(w http.ResponseWriter, cookie *http.Cookie) error {
	v := cookie.String()
	if v == "" {
		return SessionError.New("invalid cookie %#v", cookie.Name)
	}
	w.Header().Add("Set-Cookie", v)
	return nil
}

// Clear implements the Store interface. Not expected to be used directly.
func (cs *CookieStore) Clear(w http.ResponseWriter, namespace string) error {
	return setCookie(w, &http.Cookie{
		Name:     namespace,
		Value:    "",
		Path:     cs.Options.Path,
		Domain:   cs.Options.Domain,
		MaxAge:   -1,
		Secure:   cs.Options.Secure,
		HttpOnly: cs.Options.HttpOnly})
}
