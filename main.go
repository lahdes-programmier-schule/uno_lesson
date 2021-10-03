package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func getGoEnv() []string {
	cmd := exec.Command("go", "env")
	output, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	ret := []string{}
	for _, e := range strings.Split(string(output), "\n") {
		cleaned := strings.Replace(e, "\"", "", -1)
		cleaned = strings.Replace(cleaned, "set ", "", 1)
		ret = append(ret, cleaned)
	}
	return ret
}

func getGoRoot() string {
	goEnv := getGoEnv()
	for _, env := range goEnv {
		if strings.Contains(env, "GOROOT") {
			return strings.Split(env, "=")[1]
		}
	}

	log.Fatal("Unable to find GOROOT")
	return ""
}

func copyGoWasmJsFile() {
	goRoot := getGoRoot()
	goWasmFile := filepath.Join(goRoot, "misc", "wasm", "wasm_exec.js")
	target := filepath.Join("wasm_app", "wasm_exec.js")
	if err := Copy(goWasmFile, target); err != nil {
		log.Fatal(err)
	}
}

var wasmBuildArgs = []string{
	"GOOS=js",
	"GOARCH=wasm",
}

func build(debug bool) {
	env := getGoEnv()
	for _, arg := range wasmBuildArgs {
		env = append(env, arg)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	buildTmp := filepath.Join(cwd, "tmp")
	err = os.MkdirAll(buildTmp, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	env = append(env, fmt.Sprintf("GOTMPDIR=%s", buildTmp))

	if debug {
		for _, ee := range env {
			fmt.Println(ee)
		}
	}

	cmd := exec.Command("go", "build", "-o", "lib.wasm")
	cmd.Dir = "wasm_app"
	cmd.Env = env

	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		log.Fatal(err)
	}
}

func startServer() {
	http.Handle("/", http.FileServer(http.Dir("wasm_app")))

	port := "8000"
	log.Printf("Serving HTTP on port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	argsWithoutProg := os.Args[1:]
	debug := false
	if len(argsWithoutProg) > 0 {
		dd := argsWithoutProg[0]
		if dd == "--debug" {
			debug = true
		}
	}

	log.Println("Building..")
	copyGoWasmJsFile()
	build(debug)

	log.Println("Starting server..")
	startServer()
}
