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
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := checkAnswer(strings.TrimSpace(got), tc); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nreference answer:\n%s\ncontestant answer:\n%s\n", i+1, err, tc.input, tc.expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func checkAnswer(out string, tc testCase) error {
	out = strings.TrimSpace(out)
	if len(out) == 0 {
		return fmt.Errorf("output is empty")
	}
	orig := strings.TrimSpace(tc.input)
	if len(out) != len(orig) {
		return fmt.Errorf("wrong length: got %d want %d", len(out), len(orig))
	}
	for _, ch := range out {
		if ch < 'a' || ch > 'z' {
			return fmt.Errorf("output contains non-lowercase characters")
		}
	}
	if !isPalindrome(out) {
		return fmt.Errorf("output is not a palindrome")
	}
	minChanges := minimalChanges(orig)
	if neededChanges(orig, out) != minChanges {
		return fmt.Errorf("does not use minimal number of changes; need %d", minChanges)
	}
	if out != tc.expect {
		return fmt.Errorf("palindrome is not lexicographically minimal among optimal ones")
	}
	return nil
}

func isPalindrome(s string) bool {
	i, j := 0, len(s)-1
	for i < j {
		if s[i] != s[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func minimalChanges(s string) int {
	var cnt [26]int
	for _, ch := range s {
		cnt[ch-'a']++
	}
	odd := 0
	for _, c := range cnt {
		if c%2 == 1 {
			odd++
		}
	}
	if odd == 0 {
		return 0
	}
	return odd / 2
}

func neededChanges(src, dst string) int {
	var a, b [26]int
	for _, ch := range src {
		a[ch-'a']++
	}
	for _, ch := range dst {
		b[ch-'a']++
	}
	diffPos := 0
	total := 0
	for i := 0; i < 26; i++ {
		diff := b[i] - a[i]
		total += diff
		if diff > 0 {
			diffPos += diff
		}
	}
	if total != 0 {
		return -1
	}
	return diffPos
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []testCase {
	rand.Seed(42)
	samples := []string{"abba", "abb", "abc", "a", "zzz", "zxy", "mikhailrubinchikkihcniburliahkim"}
	tests := make([]testCase, 0, len(samples)+200)
	for _, s := range samples {
		tests = append(tests, makeTest(s))
	}
	for i := 0; i < 200; i++ {
		n := rand.Intn(40) + 1
		var sb strings.Builder
		for j := 0; j < n; j++ {
			sb.WriteByte(byte('a' + rand.Intn(26)))
		}
		tests = append(tests, makeTest(sb.String()))
	}
	return tests
}

func makeTest(s string) testCase {
	input := s + "\n"
	return testCase{
		input:  input,
		expect: solveRef(s),
	}
}

func solveRef(s string) string {
	var cnt [26]int
	for _, ch := range s {
		cnt[ch-'a']++
	}
	i, j := 0, 25
	for i < j {
		for i < j && cnt[i]%2 == 0 {
			i++
		}
		for i < j && cnt[j]%2 == 0 {
			j--
		}
		if i >= j {
			break
		}
		cnt[i]++
		cnt[j]--
		i++
		j--
	}
	first := make([]byte, 0, len(s)/2)
	var mid byte
	for k := 0; k < 26; k++ {
		for cnt[k] >= 2 {
			first = append(first, byte('a'+k))
			cnt[k] -= 2
		}
		if cnt[k] == 1 {
			mid = byte('a' + k)
		}
	}
	res := make([]byte, 0, len(s))
	res = append(res, first...)
	if mid != 0 {
		res = append(res, mid)
	}
	for i := len(first) - 1; i >= 0; i-- {
		res = append(res, first[i])
	}
	return string(res)
}
