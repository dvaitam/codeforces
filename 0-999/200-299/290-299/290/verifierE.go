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

const charset = "HQ9abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func expected(s string) string {
	for _, ch := range s {
		if ch == 'H' || ch == 'Q' || ch == '9' {
			return "Yes"
		}
	}
	return "No"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		s := randStr(rand.Intn(40) + 1)
		exp := expected(s)

		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(s + "\n")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\noutput:\n%s\n", i+1, err, string(out))
			os.Exit(1)
		}
		outStr := strings.TrimSpace(string(out))
		if outStr != exp {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d expected %s got %s\n", i+1, exp, outStr)
			os.Exit(1)
		}
	}
	fmt.Println("Accepted")
}
