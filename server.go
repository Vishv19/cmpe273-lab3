package main

import (
    "github.com/julienschmidt/httprouter"
    "net/http"
    "encoding/json"
    "fmt"
    "strconv"
    "os"
)

type GetKeyResponse struct {
    Key int `json:"key"`
    Value string `json:"value"`
}

type GetKeysResponse struct {
    Response []GetKeyResponse `json:"response"`
}

var storage map[int]string = make(map[int]string)

func putKey(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    key := p.ByName("key_id")
    value := p.ByName("value")
    keyInt, _ := strconv.Atoi(key)

    storage[keyInt] = value

    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
}

func getKey(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    key := p.ByName("key_id")
    keyInt, _ := strconv.Atoi(key)
    value := storage[keyInt]

    response := new(GetKeyResponse)
    response.Key = keyInt
    response.Value = value

    resJson, _ := json.Marshal(response)
    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
    fmt.Fprintf(rw, "%s", resJson)
}

func getKeys(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

    var allKeys []GetKeyResponse
    for key, value := range storage {
        var singleKey GetKeyResponse
        singleKey.Key = key
        singleKey.Value = value
        allKeys = append(allKeys, singleKey)
    }
    var allResponse GetKeysResponse
    allResponse.Response = allKeys
    resJson, _ := json.Marshal(allResponse)
    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
    fmt.Fprintf(rw, "%s", resJson)
}

func main() {
    mux := httprouter.New()
    mux.PUT("/keys/:key_id/:value", putKey)
    mux.GET("/keys/:key_id", getKey)
    mux.GET("/keys/", getKeys)
    fmt.Println("Server running at " + os.Args[1])
    address := "0.0.0.0:" + os.Args[1]
    server := http.Server{
        Addr:    address,
        Handler: mux,
    }
    server.ListenAndServe()
}