package main

import (
  "fmt"
   "github.com/CompileWithG/go-gin-auth/routes"
  "github.com/CompileWithG/go-gin-auth/config"
)

func main() {
  router:=routes.SetupRouter()
  router.Run("localhost:8080")
  //run database
  config.ConnectDB()
  fmt.Println("server running :)")

}

