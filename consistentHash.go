package main

import(
    "github.com/julienschmidt/httprouter"
    "net/http"
    "hash/crc32"
    "encoding/json"
    "fmt"
    "sort"
)
var c *Circle = new(Circle)

type GetKeyResponse struct {
    Key int `json:"key"`
    Value string `json:"value"`
}

type Circle struct {
    Nodes ServerNodes
}

func putKey(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    key := p.ByName("key_id")
    value := p.ByName("value")
    port := c.GetNode(key)
    serverurl := "http://localhost:" + port + "keys/"+ key + "/"+ value
    req, err := http.NewRequest("PUT", serverurl, nil)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    fmt.Println("Put request sent at " + serverurl)

    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
}

func getKey(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    key := p.ByName("key_id")
    port := c.GetNode(key)
    serverurl := "http://localhost:" + port + "keys/"+ key
    response, err := http.Get(serverurl)

    if err != nil {
        fmt.Println("Error while getting response from server", err.Error())
    }

    result := GetKeyResponse{}
    json.NewDecoder(response.Body).Decode(&result)
    jsondata, _ := json.Marshal(result)
    fmt.Println("Get request received from " + serverurl)
    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
    fmt.Fprintf(rw, "%s", jsondata)
}

func NewCircle() *Circle{
    var circle *Circle = new(Circle)
    circle.Nodes = ServerNodes{}
    return circle
}

func (c *Circle) AddServerNode(port string) {
    n := CreateNewNode(port)
    c.Nodes = append(c.Nodes, n)
    sort.Sort(c.Nodes)
}

func (c *Circle) GetNode(keyid string) string{
    index := c.search(keyid)
    if index >= c.Nodes.Len() {
        index = 0
    }
    return c.Nodes[index].serverPort
}

func (c *Circle) search(keyid string) int{
    searchfunction := func(i int) bool {
        return c.Nodes[i].HashdId >= hashValue(keyid)
    }

    return sort.Search(c.Nodes.Len(), searchfunction)    
}

type ServerNode struct{
	serverPort string
	HashdId uint32
}

func CreateNewNode(serverPort string) *ServerNode {
	var serverNode *ServerNode = new(ServerNode)
	serverNode.serverPort = serverPort
	serverNode.HashdId = hashValue(serverPort)

	return serverNode
}

type ServerNodes []*ServerNode

func(n ServerNodes) Len() int {
	return len(n)
}

func (n ServerNodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n ServerNodes) Less(i, j int)bool {
	return n[i].HashdId < n[j].HashdId
}

func hashValue(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func main() {
    server1Port := "3000" + "/"
    server2Port := "3001" + "/"
    server3Port := "3002" + "/"

    c.AddServerNode(server1Port)
    c.AddServerNode(server2Port)
    c.AddServerNode(server3Port)
    mux := httprouter.New()
    mux.PUT("/keys/:key_id/:value", putKey)
    mux.GET("/keys/:key_id", getKey)
    server := http.Server{
        Addr:    "0.0.0.0:8080",
        Handler: mux,
    }
    server.ListenAndServe()
}