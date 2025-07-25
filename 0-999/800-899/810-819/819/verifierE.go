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

type node struct{ a, b, c, d int }

var ans []node

func add(n int) {
	if n&1 == 1 {
		ans = append(ans, node{n, 0, n + 1, -1})
		ans = append(ans, node{n, 0, n + 1, -1})
		for i := 1; i < n; i += 2 {
			ans = append(ans, node{n, i, n + 1, i + 1})
			ans = append(ans, node{n, i, n + 1, i + 1})
		}
	} else {
		ans = append(ans, node{n, 0, n + 1, -1})
		ans = append(ans, node{n, 1, n + 1, -1})
		ans = append(ans, node{n, 0, n + 1, 1})
		for i := 2; i < n; i += 2 {
			ans = append(ans, node{n, i, n + 1, i + 1})
			ans = append(ans, node{n, i, n + 1, i + 1})
		}
	}
}

func solveE(n int) string {
	ans = nil
	if n&1 == 1 {
		ans = append(ans, node{0, 1, 2, -1})
		ans = append(ans, node{0, 1, 2, -1})
		for now := 3; now < n; now += 2 {
			add(now)
		}
	} else {
		ans = append(ans, node{0, 1, 2, -1})
		ans = append(ans, node{1, 2, 3, -1})
		ans = append(ans, node{2, 3, 0, -1})
		ans = append(ans, node{3, 0, 1, -1})
		for now := 4; now < n; now += 2 {
			add(now)
		}
	}
	var buf bytes.Buffer
	fmt.Fprintln(&buf, len(ans))
	for _, i := range ans {
		if i.d == -1 {
			fmt.Fprintln(&buf, 3, i.a+1, i.b+1, i.c+1)
		} else {
			fmt.Fprintln(&buf, 4, i.a+1, i.b+1, i.c+1, i.d+1)
		}
	}
	return strings.TrimSpace(buf.String())
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

type testCaseE struct {
	n int
}

func generateTests() []testCaseE {
	rand.Seed(42)
	tests := make([]testCaseE, 100)
	for i := range tests {
		tests[i] = testCaseE{n: rand.Intn(6) + 3}
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("%d\n", tc.n)
		expected := solveE(tc.n)
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution error: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(output) != expected {
			fmt.Printf("test %d failed:\ninput:%sexpected:\n%s\ngot:\n%s\n", i+1, input, expected, output)
			return
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
