package main

import (
	"net/http"
	"net/url"
)

func main() {
	/*http.HandleFunc("/push", pushHandler)*/
	http.HandleFunc("/test", testHandler)
	http.ListenAndServe(":50080", nil)
}

/*func pushHandler(w http.ResponseWriter, r *http.Request) {

}*/

func testHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		myUrl, _ := url.Parse(r.URL.String())

		param, _ := url.ParseQuery(myUrl.RawQuery)

		name := param.Get("name")
		hard := param.Get("hard")
		app := param.Get("app")

		w.Write([]byte("name: " + name + ", hard: " + hard + ", app: " + app))
	}

}
