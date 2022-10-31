package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type envCfg struct {
	ClientID      string `env:"ClientID"`
	ClientSecret  string `env:"ClientSecret"`
	AuthServerURL string `env:"AuthServerURL"`
	PortClient    string `env:"PortClient"`
	AuthClientURL string `env:"AuthClientURL"`
	RedirectURL   string `env:"RedirectURL"`
	AuthState     string `env:"AuthState"`
}

var cache = map[string]string{
	"state":        "",
	"redirect_uri": "",
}

func setupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Requested-With, Authorization")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := envCfg{}
	if err := env.Parse(&config); err != nil {
		fmt.Printf("%+v\n", err)
	}

	var (
		oauthConfig = oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Scopes:       []string{"all"},
			RedirectURL:  config.RedirectURL,
			Endpoint: oauth2.Endpoint{
				AuthURL:  config.AuthServerURL + "/oauth/authorize",
				TokenURL: config.AuthServerURL + "/oauth/token",
			},
		}
		globalToken *oauth2.Token
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}

		var redirect_uri = r.URL.Query().Get("redirect_uri")
		redirect_uri, _ = url.QueryUnescape(redirect_uri)
		if redirect_uri != "" {
			cache["redirect_uri"] = redirect_uri
			cache["state"] = r.URL.Query().Get("state")
		}

		u := oauthConfig.AuthCodeURL(config.AuthState,
			oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256("s256example")),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"))
		http.Redirect(w, r, u, http.StatusFound)
	})

	http.HandleFunc("/oauth2", func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}

		r.ParseForm()
		state := r.Form.Get("state")
		if state != config.AuthState {
			http.Error(w, "State invalid", http.StatusBadRequest)
			return
		}

		code := r.Form.Get("code")
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": "test",
			"password": "test",
		})

		tokenString, error := token.SignedString([]byte(config.ClientSecret))
		if error != nil {
			fmt.Println(error)
		}

		url := cache["redirect_uri"] +
			"#state=" + cache["state"] +
			"&access_token=" + tokenString +
			"&token_type=Bearer"
		http.Redirect(w, r, url, http.StatusFound)
	})

	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}

		if globalToken == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		globalToken.Expiry = time.Now()
		token, err := oauthConfig.TokenSource(context.Background(), globalToken).Token()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		globalToken = token
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	http.HandleFunc("/listar-itens", func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}

		tokenStr := r.Header.Get("Authorization")
		token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return []byte(config.ClientSecret), nil
		})

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			w.Write([]byte("Vitrola, Fita cassete, Zip drive"))
		} else {
			w.Write([]byte("Acesso não autorizado"))
		}
	})

	log.Println("Client is running at " + config.PortClient + " port. Please open " + config.AuthClientURL)
	log.Fatal(http.ListenAndServeTLS(":"+config.PortClient, "client.crt", "client.key", nil))
}

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}
