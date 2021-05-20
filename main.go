package main

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

var (
	googleOauthConfig *oauth2.Config
	fbOauthConfig     *oauth2.Config
	ghOauthConfig     *oauth2.Config
	oauthStateString  = "pseudo-random"
)

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login/google", handleGoogleLogin)
	http.HandleFunc("/callback/google", handleGoogleCallback)
	http.HandleFunc("/login/facebook", handleFacebookLogin)
  http.HandleFunc("/callback/facebook", handleFacebookCallback)
	http.HandleFunc("/login/github", handleGithubLogin)
  http.HandleFunc("/callback/github", handleGithubCallback)
	fmt.Println(http.ListenAndServe(":8080", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html>
	<script>
  window.fbAsyncInit = function() {
    FB.init({
      appId      : '474256587200310',
      xfbml      : true,
      version    : 'v10.0'
    });
    FB.AppEvents.logPageView();
  };

  (function(d, s, id){
     var js, fjs = d.getElementsByTagName(s)[0];
     if (d.getElementById(id)) {return;}
     js = d.createElement(s); js.id = id;
     js.src = "https://connect.facebook.net/en_US/sdk.js";
     fjs.parentNode.insertBefore(js, fjs);
   }(document, 'script', 'facebook-jssdk'));
</script>
<body>
	<a href="/login/google">Google Log In</a><br>
	<a href="/login/facebook">Facebook Log In</a><br>
	<a href="/login/github">Github Log In</a><br>
</body>
</html>`

	fmt.Fprintf(w, htmlIndex)
}

