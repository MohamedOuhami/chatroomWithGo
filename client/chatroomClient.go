package main

import (
	"fmt"
	"net"
)

func main() {
  
  // Connect to the server
  conn,err := net.Dial("tcp","localhost:8000")

  if err != nil {
    fmt.Println("There was an error",err)
    return
  }

  data := []byte("Hello, Server")
  _,err = conn.Write(data)

  if err != nil {
    fmt.Println("Error:",err)
    return
  }

  defer conn.Close()
}
