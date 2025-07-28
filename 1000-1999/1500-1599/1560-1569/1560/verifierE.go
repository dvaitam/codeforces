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
	t string
}

func solve(t string) (string, string) {
	seen := make(map[byte]bool)
	orderRev := make([]byte, 0)
	for i := len(t) - 1; i >= 0; i-- {
		ch := t[i]
		if !seen[ch] {
			seen[ch] = true
			orderRev = append(orderRev, ch)
		}
	}
	order := make([]byte, len(orderRev))
	for i := range orderRev {
		order[len(orderRev)-1-i] = orderRev[i]
	}

	freq := make(map[byte]int)
	for i := 0; i < len(t); i++ {
		freq[t[i]]++
	}

	origCount := make(map[byte]int)
	prefLen := 0
	for i, ch := range order {
		c := freq[ch]
		step := i + 1
		if c%step != 0 {
			return "-1", ""
		}
		origCount[ch] = c / step
		prefLen += origCount[ch]
	}
	if prefLen > len(t) {
		return "-1", ""
	}
	s := t[:prefLen]
	cur := []byte(s)
	built := make([]byte, 0, len(t))
	for _, ch := range order {
		built = append(built, cur...)
		filtered := make([]byte, 0, len(cur))
		for _, c := range cur {
			if c != ch {
				filtered = append(filtered, c)
			}
		}
		cur = filtered
	}
	if string(built) != t {
		return "-1", ""
	}
	return s, string(order)
}

func genTests() []testCase {
	r := rand.New(rand.NewSource(5))
	tests := make([]testCase, 100)
	letters := "abcde"
	for i := range tests {
		l := r.Intn(10) + 1
		var sb strings.Builder
		for j := 0; j < l; j++ {
			sb.WriteByte(letters[r.Intn(len(letters))])
		}
		tests[i].t = sb.String()
	}
	return tests
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%s\n", tc.t)
		s, ord := solve(tc.t)
		want := "-1"
		if s != "-1" {
			want = fmt.Sprintf("%s %s", s, ord)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, want, got)
			fmt.Printf("input:\n%s", input)
			return
		}
	}
	fmt.Println("All tests passed")
}
