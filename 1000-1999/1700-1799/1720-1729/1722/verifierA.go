package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCaseA struct {
	n int
	s string
}

func solveA(n int, s string) string {
	if n != 5 {
		return "NO"
	}
	letters := []rune(s)
	sort.Slice(letters, func(i, j int) bool { return letters[i] < letters[j] })
	if string(letters) == "Timru" {
		return "YES"
	}
	return "NO"
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateTests() []testCaseA {
	rand.Seed(42)
	tests := make([]testCaseA, 100)
	for i := range tests {
		n := rand.Intn(10) + 1
		s := make([]rune, n)
		if n == 5 && rand.Intn(2) == 0 {
			t := []rune("Timur")
			rand.Shuffle(5, func(i, j int) { t[i], t[j] = t[j], t[i] })
			copy(s, t)
		} else {
			for j := 0; j < n; j++ {
				s[j] = letters[rand.Intn(len(letters))]
			}
		}
		tests[i] = testCaseA{n: n, s: string(s)}
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d\n%s\n", tc.n, tc.s)
		expected := solveA(tc.n, tc.s)
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution error: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(output) != expected {
			fmt.Printf("test %d failed:\ninput:%sexpected %s got %s\n", i+1, input, expected, output)
			return
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
