package main

import (
	"fmt"
	"task-management/internal/api"
)


func main()  {
	addr := "0.0.0.0:8000"
	httpServer := api.NewHttpServer()
	err :=httpServer.ListenAndServe(addr)
	if err != nil {
		fmt.Println("err",err)
		panic("an error occurred")
	}

}