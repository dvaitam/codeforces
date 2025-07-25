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

func genCase(rng *rand.Rand) (string, string) {
	numWords := rng.Intn(5) + 1
	words := make([]string, numWords)
	maxCap := 0
	for i := 0; i < numWords; i++ {
		l := rng.Intn(10) + 1
		b := make([]byte, l)
		caps := 0
		for j := 0; j < l; j++ {
			if rng.Intn(2) == 0 {
				b[j] = byte('a' + rng.Intn(26))
			} else {
				b[j] = byte('A' + rng.Intn(26))
				caps++
			}
		}
		if caps > maxCap {
			maxCap = caps
		}
		words[i] = string(b)
	}
	text := strings.Join(words, " ")
	n := len(text)
	input := fmt.Sprintf("%d\n%s\n", n, text)
	expected := fmt.Sprintf("%d", maxCap)
	return input, expected
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
