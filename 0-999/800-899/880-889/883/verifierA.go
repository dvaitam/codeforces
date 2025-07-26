package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func runCmd(name string, args []string, input string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func runProg(prog string, input string) (string, error) {
	if strings.HasSuffix(prog, ".go") {
		return runCmd("go", []string{"run", prog}, input)
	}
	return runCmd(prog, nil, input)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	ref := filepath.Join(dir, "883A.go")
	rand.Seed(1)
	for tcase := 0; tcase < 100; tcase++ {
		n := rand.Int63n(5) + 1
		m := rand.Int63n(5) + 1
		a := rand.Int63n(5) + 1
		d := rand.Int63n(10) + 1
		clients := make([]int64, m)
		cur := rand.Int63n(5) + 1
		for i := int64(0); i < m; i++ {
			cur += rand.Int63n(5) + 1
			clients[i] = cur
		}
		input := fmt.Sprintf("%d %d %d %d\n", n, m, a, d)
		for i := int64(0); i < m; i++ {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprint(clients[i])
		}
		input += "\n"

		expOut, expErr := runCmd("go", []string{"run", ref}, input)
		if expErr != nil {
			fmt.Println("reference error:", expErr)
			os.Exit(1)
		}
		actOut, actErr := runProg(bin, input)
		if actErr != nil {
			fmt.Printf("runtime error on test %d: %v\n", tcase+1, actErr)
			fmt.Println("input:", input)
			os.Exit(1)
		}
		if strings.TrimSpace(expOut) != strings.TrimSpace(actOut) {
			fmt.Printf("wrong answer on test %d\n", tcase+1)
			fmt.Println("input:", input)
			fmt.Println("expected:", strings.TrimSpace(expOut))
			fmt.Println("got:", strings.TrimSpace(actOut))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
