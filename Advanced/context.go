package main

import (
	"context"
	"fmt"
)


func main() {
	contextTodo := context.TODO()
	contextBc := context.Background()
	
	
	ctx := context.WithValue(contextTodo, "name", "Jhon" )
	fmt.Println(ctx)
	fmt.Println(ctx.Value("name"))
	
	ctxBC := context.WithValue(contextBc, "city", "Parana" )
	fmt.Println(ctxBC)
	fmt.Println(ctxBC.Value("city"))
}