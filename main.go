package main

//import "programmierschule.de/uno/uno"
import "programmierschule.de/uno/dom"
import "syscall/js"

func main() {
    doc := dom.GetWindow().Document()
	doc.GetElementById("content").SetInnerHTML("<h1>UNO!</h1>")
}
