package auth

import (
	"net/http"
	"time"

	"prediction/internal/config"
)

type CookieManager struct {
	config config.AuthConfig
}

func NewCookieManager(config config.AuthConfig) *CookieManager {
	return &CookieManager{config: config}
}

func (m *CookieManager) SetSessionCookie(w http.ResponseWriter, rawToken string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     m.config.CookieName,
		Value:    rawToken,
		Path:     m.config.CookiePath,
		Domain:   m.config.CookieDomain,
		HttpOnly: true,
		Secure:   m.config.CookieSecure,
		SameSite: m.config.CookieSameSite,
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
	})
}

func (m *CookieManager) ReadSessionCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(m.config.CookieName)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
