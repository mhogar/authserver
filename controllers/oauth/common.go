package oauth

type passwordGrantBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Scope    string `json:"scope"`
}

type authorizationCodeGrantBody struct {
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
