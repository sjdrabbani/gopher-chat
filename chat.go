package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	pusher "github.com/pusher/pusher-http-go"
)

var client = pusher.Client{
	AppId:   os.Getenv("PUSHER_APP_ID"),
	Key:     os.Getenv("PUSHER_APP_KEY"),
	Secret:  os.Getenv("PUSHER_APP_SECRET"),
	Cluster: os.Getenv("PUSHER_APP_CLUSTER"),
	Secure:  true,
}

type user struct {
	Name  string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}

func registerNewUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var newUser user

	err = json.Unmarshal(body, &newUser)
	if err != nil {
		panic(err)
	}

	client.Trigger("update", "new-user", newUser)
	json.NewEncoder(w).Encode(newUser)
}

func pusherAuth(w http.ResponseWriter, r *http.Request) {
	params, _ := ioutil.ReadAll(r.Body)
	response, err := client.AuthenticatePrivateChannel(params)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, string(response))
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/new/user", registerNewUser)
	http.HandleFunc("/pusher/auth", pusherAuth)

	log.Fatal(http.ListenAndServe(":8090", nil))
}
