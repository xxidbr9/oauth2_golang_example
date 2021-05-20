package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// request scope here
	fbOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
		ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
		Scopes:       []string{"public_profile", "email"},
		RedirectURL:  "http://localhost:8080/callback/facebook",
		Endpoint:     facebook.Endpoint,
	}
}

func handleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	Url, err := url.Parse(fbOauthConfig.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", fbOauthConfig.ClientID)
	parameters.Add("scope", strings.Join(fbOauthConfig.Scopes, " "))
	parameters.Add("redirect_uri", fbOauthConfig.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	Url.RawQuery = parameters.Encode()
	url := Url.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleFacebookCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	token, err := fbOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("fbOauthConfig.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// return in here
	resp, err := http.Get("https://graph.facebook.com/me?access_token=" +
		url.QueryEscape(token.AccessToken) + "&fields=id,name,email")
	if err != nil {
		fmt.Printf("Get: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp)
	fmt.Println(url.QueryEscape(token.AccessToken))
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ReadAll: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	log.Printf("parseResponseBody: %s\n", string(response))
	fmt.Fprintf(w, "Content: %s\n", string(response))

	// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
