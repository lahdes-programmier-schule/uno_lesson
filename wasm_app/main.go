package main

import "programmierschule.de/uno/dom"
// import "programmierschule.de/uno/uno"
import "syscall/js"
import "fmt"

func test(this js.Value, inputs []js.Value) interface{} {
	doc := dom.GetWindow().Document()
	doc.GetElementById("content").SetInnerHTML("<h1>UNO!</h1>")
	return nil
}

func waitForEvents() {
	c := make(chan struct{}, 0)
	<-c
}

func main() {
    fmt.Println("starting up...")

	js.Global().Set("test", js.FuncOf(test))
    
    waitForEvents()
}
