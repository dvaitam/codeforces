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
	cntA, cntB, cntC := 0, 0, 0
	for _, ch := range s {
		switch ch {
		case 'A':
			cntA++
		case 'B':
			cntB++
		case 'C':
			cntC++
		}
	}
	if cntB == cntA+cntC {
		return "YES\n"
	}
	return "NO\n"
}

func genCaseA(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	letters := []byte{'A', 'B', 'C'}
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(3)]
	}
	s := string(b)
	input := fmt.Sprintf("1\n%s\n", s)
	expect := solveA(s)
	return input, expect
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
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := genCaseA(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
