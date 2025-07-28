package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func generateInput() []byte {
	rand.Seed(42)
	t := 100
	var b bytes.Buffer
	fmt.Fprintf(&b, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rand.Intn(10) + 1
		k := rand.Intn(n + 1)
		fmt.Fprintf(&b, "%d %d\n", n, k)
		for j := 0; j < n; j++ {
			if j > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", rand.Intn(10))
		}
		b.WriteByte('\n')
	}
	return b.Bytes()
}

func runBinary(path string, input []byte) ([]byte, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	return cmd.CombinedOutput()
}

func runReference(input []byte) ([]byte, error) {
	cmd := exec.Command("go", "run", "1891E.go")
	cmd.Stdin = bytes.NewReader(input)
	return cmd.CombinedOutput()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	input := generateInput()

	expOut, err := runReference(input)
	if err != nil {
		fmt.Println("error running reference:", err)
		os.Exit(1)
	}

	gotOut, err := runBinary(os.Args[1], input)
	if err != nil {
		fmt.Println("error running binary:", err)
		os.Exit(1)
	}

	expLines := strings.Fields(strings.TrimSpace(string(expOut)))
	gotLines := strings.Fields(strings.TrimSpace(string(gotOut)))
	if len(expLines) != len(gotLines) {
		fmt.Printf("mismatch line count: expected %d got %d\n", len(expLines), len(gotLines))
		os.Exit(1)
	}
	for i := range expLines {
		if expLines[i] != gotLines[i] {
			fmt.Printf("mismatch on line %d: expected %s got %s\n", i+1, expLines[i], gotLines[i])
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
