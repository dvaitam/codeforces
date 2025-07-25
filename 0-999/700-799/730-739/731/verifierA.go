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

func solveA(s string) string {
	curr := 'a'
	total := 0
	for _, ch := range s {
		diff := int(ch - curr)
		if diff < 0 {
			diff = -diff
		}
		if diff > 26-diff {
			diff = 26 - diff
		}
		total += diff
		curr = ch
	}
	return fmt.Sprintf("%d\n", total)
}

func genCaseA(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(rng.Intn(26) + 'a')
	}
	s := string(b)
	return s + "\n", solveA(s)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, expect := genCaseA(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
