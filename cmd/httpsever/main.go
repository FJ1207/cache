package main

import (
//	"fmt"
	"log"
	"net/http"
)

type server int
func (hs *server) ServeHTTP(w http.ResponseWriter,r *http.Request)  {
	log.Println(r.URL.Path)
	w.Write([]byte("Hello World!"))
}

func main() {
	var s server
	http.ListenAndServe("localhost:9999", &s)
}



//第二种
// func Hellohandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hello!")
// }
// func printhandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "王老师，我肯定能学好的！别生气奥")
// }

// func main() {
// 	http.Handle("/", http.HandlerFunc(Hellohandler))
// 	// http.HandleFunc("/", Hellohandler)
// 	http.HandleFunc("/love", LoveWDYhandler)
// 	http.HandleFunc("/love/say", printhandler)
// 	fmt.Println("server is running...")
// 	if err := http.ListenAndServe("192.168.52.1:80", nil); err != nil {
// 		fmt.Println(err)
// 	}
// }

// func LoveWDYhandler(w http.ResponseWriter, r *http.Request) {
// 	s := "fj♡wdy"
// 	fmt.Fprintln(w, s)
// }
