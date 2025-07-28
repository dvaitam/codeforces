package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input  string
	expect string
}

func solve(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var out strings.Builder
	for scanner.Scan() {
		out.WriteString("NO\n")
	}
	return strings.TrimSpace(out.String())
}

func genTests() []testCase {
	r := rand.New(rand.NewSource(1))
	tests := make([]testCase, 100)
	for i := range tests {
		tok := r.Intn(10) + 1
		var sb strings.Builder
		for j := 0; j < tok; j++ {
			if j > 0 {
				if r.Intn(2) == 0 {
					sb.WriteByte(' ')
				} else {
					sb.WriteByte('\n')
				}
			}
			l := r.Intn(5) + 1
			for k := 0; k < l; k++ {
				sb.WriteByte(byte('a' + r.Intn(26)))
			}
		}
		sb.WriteByte('\n')
		tests[i].input = sb.String()
		tests[i].expect = solve(tests[i].input)
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := genTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			fmt.Println("input:")
			fmt.Println(tc.input)
			return
		}
		if got != strings.TrimSpace(tc.expect) {
			fmt.Printf("case %d failed:\nexpected: %q\n   got: %q\n", i+1, tc.expect, got)
			fmt.Println("input:")
			fmt.Println(tc.input)
			return
		}
	}
	fmt.Println("All tests passed")
}
