package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"io/ioutil"
	"net/http"
)

// var accessToken string

var GithubConfig = &oauth2.Config{
	ClientID:     "Ov23li7uTl01RGfo9LJT",
	ClientSecret: "78dc723c39b9c54f8b0808b8710640c24bfdc89c",
	RedirectURL:  "http://localhost:8000/callback",
	Endpoint:     github.Endpoint,
	Scopes:       []string{"repo"},
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/auth", generateAuthURL).Methods("GET")
	router.HandleFunc("/callback", handleCallback).Methods("GET")

	router.HandleFunc("/repos", getRepos).Methods("GET")

	http.ListenAndServe(":8000", router)
}

func generateAuthURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HomeHandler")
	url := GithubConfig.AuthCodeURL("random_state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, _ := GithubConfig.Exchange(r.Context(), code)

	fmt.Println("token", token.AccessToken)
	response := fmt.Sprintf(`
        <html>
            <body>
                <h1>GitHub OAuth</h1>
                <p><strong>Access Token:</strong><code> %s</code></p>
            </body>
        </html>`, token.AccessToken)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func getRepos(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("Authorization")
	url := "https://api.github.com/repos/zohaib-shamshad/uniqode-storefront"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
