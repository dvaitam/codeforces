package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	n    int
	m    int
	name string
	ops  [][2]byte
}

func generateTests() []testCase {
	tests := []testCase{
		{1, 1, "a", [][2]byte{{'a', 'b'}}},
		{3, 2, "abc", [][2]byte{{'a', 'b'}, {'b', 'c'}}},
	}
	rng := rand.New(rand.NewSource(2))
	for len(tests) < 100 {
		n := rng.Intn(20) + 1
		m := rng.Intn(20) + 1
		nameBytes := make([]byte, n)
		for i := 0; i < n; i++ {
			nameBytes[i] = byte('a' + rng.Intn(26))
		}
		ops := make([][2]byte, m)
		for i := 0; i < m; i++ {
			ops[i][0] = byte('a' + rng.Intn(26))
			ops[i][1] = byte('a' + rng.Intn(26))
		}
		tests = append(tests, testCase{n, m, string(nameBytes), ops})
	}
	return tests
}

func solve(tc testCase) string {
	mapping := [26]byte{}
	for i := 0; i < 26; i++ {
		mapping[i] = byte('a' + i)
	}
	for _, op := range tc.ops {
		x := op[0]
		y := op[1]
		for j := 0; j < 26; j++ {
			if mapping[j] == x {
				mapping[j] = y
			} else if mapping[j] == y {
				mapping[j] = x
			}
		}
	}
	bytes := []byte(tc.name)
	for i := 0; i < len(bytes); i++ {
		bytes[i] = mapping[bytes[i]-'a']
	}
	return string(bytes)
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
		input.WriteString(tc.name)
		input.WriteByte('\n')
		for _, op := range tc.ops {
			fmt.Fprintf(&input, "%c %c\n", op[0], op[1])
		}
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		exp := solve(tc)
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input.String(), exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
