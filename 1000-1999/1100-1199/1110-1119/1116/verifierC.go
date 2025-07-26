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

func alternating(s string) int {
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			return 0
		}
	}
	return 1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		n := rand.Intn(20) + 1
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				b[j] = '0'
			} else {
				b[j] = '1'
			}
		}
		input := fmt.Sprintf("%d\n%s\n", n, string(b))
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "run %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		valStr := strings.TrimSpace(out.String())
		if valStr != "0" && valStr != "1" {
			fmt.Fprintf(os.Stderr, "test %d: invalid output %q\n", i+1, valStr)
			os.Exit(1)
		}
		want := fmt.Sprintf("%d", alternating(string(b)))
		if valStr != want {
			fmt.Fprintf(os.Stderr, "test %d: got %s want %s input %s\n", i+1, valStr, want, string(b))
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
