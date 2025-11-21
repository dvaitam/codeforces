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

func solveCase(s string) string {
	cnt := make([]int64, 11)
	next := make([]int64, 11)
	var res int64
	for i := 0; i < len(s); i++ {
		d := int(s[i] - '0')
		for j := 0; j < 11; j++ {
			next[j] = 0
		}
		for j := d + 1; j < 11; j++ {
			nj := int((int64(j)*(int64(j)-1)/2 + int64(d) + 10) % 11)
			next[nj] += cnt[j]
		}
		cnt, next = next, cnt
		if d != 0 {
			cnt[d]++
		}
		for j := 0; j < 11; j++ {
			res += cnt[j]
		}
	}
	return fmt.Sprintf("%d", res)
}

func formatInput(s string) string {
	return s + "\n"
}

func fixedTests() []test {
	cases := []string{
		"1",
		"12",
		"4021",
		"10",
		"99999",
		"101010",
	}
	var tests []test
	for _, s := range cases {
		tests = append(tests, test{
			input:    formatInput(s),
			expected: solveCase(s),
		})
	}
	return tests
}

func randomString(rng *rand.Rand, length int) string {
	if length == 0 {
		return ""
	}
	sb := make([]byte, length)
	sb[0] = byte(rng.Intn(9) + '1')
	for i := 1; i < length; i++ {
		sb[i] = byte(rng.Intn(10) + '0')
	}
	return string(sb)
}

func randomTests(rng *rand.Rand, count, minLen, maxLen int) []test {
	tests := make([]test, 0, count)
	for len(tests) < count {
		length := rng.Intn(maxLen-minLen+1) + minLen
		s := randomString(rng, length)
		tests = append(tests, test{
			input:    formatInput(s),
			expected: solveCase(s),
		})
	}
	return tests
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(1142))
	tests := fixedTests()
	tests = append(tests, randomTests(rng, 40, 1, 20)...)
	tests = append(tests, randomTests(rng, 25, 21, 60)...)
	tests = append(tests, randomTests(rng, 10, 61, 200)...)
	return tests
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\nInput:%s\n", i+1, err, t.input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
