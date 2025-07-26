package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func expected(name, value string) string {
	prefix := 'i'
	if strings.Contains(value, ".") {
		prefix = 'f'
	}
	runes := []rune(name)
	if len(runes) > 0 {
		runes[0] = unicode.ToUpper(runes[0])
	}
	return fmt.Sprintf("%c%s", prefix, string(runes))
}

func randName() string {
	n := rand.Intn(10) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rand.Intn(26))
	}
	return string(b)
}

func randValue() string {
	if rand.Intn(2) == 0 {
		// integer
		return fmt.Sprintf("%d", rand.Intn(100000))
	}
	// real
	return fmt.Sprintf("%d.%d", rand.Intn(1000), rand.Intn(1000))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(42)
	for t := 0; t < 100; t++ {
		name := randName()
		value := randValue()
		in := fmt.Sprintf("%s\n%s\n", name, value)
		want := expected(name, value)
		out, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != want {
			fmt.Printf("test %d failed: expected %q got %q\n", t+1, want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
