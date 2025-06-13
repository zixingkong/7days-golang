package main

// $ curl http://localhost:9999/
// URL.Path = "/"
// $ curl http://localhost:9999/hello
// Header["Accept"] = ["*/*"]
// Header["User-Agent"] = ["curl/7.54.0"]

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)      // 处理根路径请求
	http.HandleFunc("/hello", helloHandler) // 处理/hello路径请求
	http.HandleFunc("/demo", demoHandler)   // 新增：处理/demo路径请求
	log.Fatal(http.ListenAndServe(":9999", nil))
}

// handler echoes r.URL.Path
func indexHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

// handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, req *http.Request) {

	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}

// 新增：处理/demo路径请求
func demoHandler(w http.ResponseWriter, req *http.Request) {
	
	fmt.Fprintf(w, "This is a demo page\n")
}
