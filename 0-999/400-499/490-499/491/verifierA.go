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

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expected(a, b int) string {
	var sb strings.Builder
	for i := b + 1; i <= b+a+1; i++ {
		fmt.Fprintf(&sb, "%d ", i)
	}
	for i := b; i > 0; i-- {
		fmt.Fprintf(&sb, "%d ", i)
	}
	return strings.TrimSpace(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	const tests = 100
	for t := 0; t < tests; t++ {
		a := rand.Intn(1000)
		b := rand.Intn(1000)
		input := fmt.Sprintf("%d %d\n", a, b)
		want := expected(a, b)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nOutput:\n%s\n", t+1, err, got)
			return
		}
		if got != want {
			fmt.Printf("Test %d failed.\nInput: %sExpected: %s\nGot: %s\n", t+1, input, want, got)
			return
		}
	}
	fmt.Println("All tests passed.")
}
