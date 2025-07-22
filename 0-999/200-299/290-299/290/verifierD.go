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

func shift(s string, k int) string {
	r := []rune(s)
	for i, ch := range r {
		switch {
		case ch >= 'a' && ch <= 'z':
			r[i] = 'a' + (ch-'a'+rune(k))%26
		case ch >= 'A' && ch <= 'Z':
			r[i] = 'A' + (ch-'A'+rune(k))%26
		}
	}
	return string(r)
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		s := randString(rand.Intn(20) + 1)
		k := rand.Intn(27)
		expected := shift(s, k)
		input := fmt.Sprintf("%s\n%d\n", s, k)

		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\noutput:\n%s\n", i+1, err, string(out))
			os.Exit(1)
		}
		outStr := strings.TrimSpace(string(out))
		if outStr != expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d expected %s got %s\n", i+1, expected, outStr)
			os.Exit(1)
		}
	}
	fmt.Println("Accepted")
}
