package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randString(r *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + r.Intn(3))
	}
	return string(b)
}

func generateTests() []string {
    r := rand.New(rand.NewSource(5))
    tests := make([]string, 0, 100)
    for t := 0; t < 100; t++ {
        n := r.Intn(6) + 1
        s := randString(r, n)
        tstr := randString(r, n)
        // Ensure S and T are distinct as per problem statement
        for tstr == s {
            // tweak one character deterministically within small alphabet
            b := []byte(tstr)
            pos := r.Intn(n)
            // rotate among 'a','b','c'
            if b[pos] == 'a' { b[pos] = 'b' } else if b[pos] == 'b' { b[pos] = 'c' } else { b[pos] = 'a' }
            tstr = string(b)
        }
        tests = append(tests, fmt.Sprintf("%d\n%s\n%s\n", n, s, tstr))
    }
    return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	cand := os.Args[1]
	official := "./officialE"
	if err := exec.Command("go", "build", "-o", official, "532E.go").Run(); err != nil {
		fmt.Println("failed to build official solution:", err)
		os.Exit(1)
	}
	defer os.Remove(official)
	tests := generateTests()
	for i, tc := range tests {
		exp, eerr := runBinary(official, tc)
		got, gerr := runBinary(cand, tc)
		if eerr != nil {
			fmt.Printf("official solution failed on test %d: %v\n", i+1, eerr)
			os.Exit(1)
		}
		if gerr != nil {
			fmt.Printf("candidate failed on test %d: %v\n", i+1, gerr)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
