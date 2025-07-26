package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input    string
	expected string
}

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveB(input string) string {
	fields := strings.Fields(input)
	idx := 0
	n := 0
	fmt.Sscan(fields[idx], &n)
	idx++
	var k int64
	fmt.Sscan(fields[idx], &k)
	idx++
	s1 := fields[idx]
	idx++
	s2 := fields[idx]
	cnt := int64(0)
	cur := int64(1)
	ans := int64(0)
	for i := 0; i < n; i++ {
		cur <<= 1
		if s1[i] == 'b' {
			cur--
		}
		if s2[i] == 'a' {
			cur--
		}
		if cur > k {
			ans = cnt + k*int64(n-i)
			return fmt.Sprintf("%d", ans)
		}
		cnt += cur
		ans = cnt + cur*int64(n-1-i)
	}
	return fmt.Sprintf("%d", ans)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(43))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(8) + 1
		k := rng.Int63n(10) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		b1 := make([]byte, n)
		b2 := make([]byte, n)
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				b1[i] = 'a'
			} else {
				b1[i] = 'b'
			}
			if rng.Intn(2) == 0 {
				b2[i] = 'a'
			} else {
				b2[i] = 'b'
			}
		}
		s1 := string(b1)
		s2 := string(b2)
		if s1 > s2 {
			s1, s2 = s2, s1
		}
		sb.WriteString(fmt.Sprintf("%s\n%s\n", s1, s2))
		input := sb.String()
		expected := solveB(strings.TrimSpace(fmt.Sprintf("%d %d %s %s", n, k, s1, s2)))
		tests = append(tests, test{input, expected})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
