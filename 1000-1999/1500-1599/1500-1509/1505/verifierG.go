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

type testCase struct {
	input  string
	expect string
}

var enc = map[byte][5]int{
	'a': {1, 0, 0, 1, 0},
	'c': {2, 0, 0, 1, 1},
	'o': {1, 1, 1, 2, 1},
	'd': {2, 1, 0, 1, 2},
	'e': {1, 1, 0, 1, 1},
	'f': {2, 1, 0, 2, 1},
	'r': {1, 2, 1, 3, 1},
	'z': {1, 1, 2, 2, 2},
}

func decode(vals [5]int) (byte, bool) {
	for ch, v := range enc {
		if v == vals {
			return ch, true
		}
	}
	return '?', false
}

func solve(lines [][5]int) string {
	var sb strings.Builder
	for _, v := range lines {
		ch, ok := decode(v)
		if !ok {
			return "?"
		}
		sb.WriteByte(ch)
	}
	return sb.String()
}

func genTests() []testCase {
	r := rand.New(rand.NewSource(1))
	letters := []byte{'a', 'c', 'o', 'd', 'e', 'f', 'r', 'z'}
	tests := make([]testCase, 100)
	for i := range tests {
		n := r.Intn(10) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		seq := make([][5]int, n)
		for j := 0; j < n; j++ {
			ch := letters[r.Intn(len(letters))]
			pattern := enc[ch]
			seq[j] = pattern
			sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", pattern[0], pattern[1], pattern[2], pattern[3], pattern[4]))
		}
		tests[i].input = sb.String()
		tests[i].expect = solve(seq)
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
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := genTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			fmt.Print(tc.input)
			return
		}
		if got != tc.expect {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, tc.expect, got)
			fmt.Print(tc.input)
			return
		}
	}
	fmt.Println("All tests passed")
}
