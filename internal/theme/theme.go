package theme

import "net/http"

var themeCookieName = "brattonross_theme"

func GetTheme(r *http.Request) string {
	cookie, _ := r.Cookie(themeCookieName)
	if cookie != nil {
		return cookie.Value
	}
	return ""
}

func SetTheme(w http.ResponseWriter, theme string) {
	http.SetCookie(w, &http.Cookie{
		Name:     themeCookieName,
		Value:    theme,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		MaxAge:   60 * 60 * 24 * 365 * 100, // 100 years
		Path:     "/",
	})
}
