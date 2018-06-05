package main

import "fmt"


type abc struct{
	a int
	b int
}


func main(){
	a := abc{a:10,b:20}
	// a.a := 10

	b := &a

	b.a = 11
	fmt.Println(a)
}