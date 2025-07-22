package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildRef() (string, error) {
	ref := "refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "459E.go")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return "./" + ref, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Intn(8) + 2
		maxEdges := n * (n - 1)
		m := rand.Intn(maxEdges) + 1
		input := fmt.Sprintf("%d %d\n", n, m)
		for i := 0; i < m; i++ {
			u := rand.Intn(n) + 1
			v := rand.Intn(n) + 1
			w := rand.Intn(20) + 1
			input += fmt.Sprintf("%d %d %d\n", u, v, w)
		}

		expect, err := runBinary(ref, input)
		if err != nil {
			fmt.Println("reference runtime error:", err)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(expect) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:%sexpected:%s\n got:%s\n", t+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
