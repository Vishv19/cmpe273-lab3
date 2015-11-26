##**To use consistent hashing**

**go run consistentHash.go**

Here consistentHash.go is a "go client" which distributes keys among three servers

Start server at port 3000, 3001, 3002 using following commands
**go run server.go 3000**

**go run server.go 3001**

**go run server.go 3002**

#**Sample PUT request at client:**
http://localhost:8080/keys/1/a

PUT Response

200 response code

#**Sample GET Request at client:**

http://localhost:8080/keys/1

GET Response

{
  "key": 1,
  "value": "a"
}

#**Sample PUT request at server running at port 3000:**

http://localhost:3000/keys/1/a

PUT Response

200 response code

#**Sample GET Request at server running at port 3000:**

http://localhost:3000/keys/1

GET Response

{
  "key": 1,
  "value": "a"
}
