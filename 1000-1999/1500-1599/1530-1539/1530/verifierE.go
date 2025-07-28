package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testE struct{ s string }

func genTestsE() []testE {
	rand.Seed(1530005)
	tests := make([]testE, 100)
	for i := range tests {
		n := rand.Intn(20) + 1
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			b[j] = byte('a' + rand.Intn(26))
		}
		tests[i].s = string(b)
	}
	return tests
}

func solveE(tc testE) string {
	s := tc.s
	freq := make([]int, 26)
	for i := 0; i < len(s); i++ {
		freq[s[i]-'a']++
	}
	distinct := 0
	for i := 0; i < 26; i++ {
		if freq[i] > 0 {
			distinct++
		}
	}
	if distinct == 1 {
		return s
	}
	unique := -1
	for i := 0; i < 26; i++ {
		if freq[i] == 1 {
			unique = i
			break
		}
	}
	if unique != -1 {
		var b strings.Builder
		b.WriteByte(byte('a' + unique))
		freq[unique]--
		for i := 0; i < 26; i++ {
			for freq[i] > 0 {
				b.WriteByte(byte('a' + i))
				freq[i]--
			}
		}
		return b.String()
	}
	first := 0
	for freq[first] == 0 {
		first++
	}
	n := len(s)
	cntFirst := freq[first]
	if cntFirst-2 <= n-cntFirst {
		var b strings.Builder
		b.WriteByte(byte('a' + first))
		b.WriteByte(byte('a' + first))
		freq[first] -= 2
		others := make([]byte, 0, n-cntFirst)
		for i := first + 1; i < 26; i++ {
			for j := 0; j < freq[i]; j++ {
				others = append(others, byte('a'+i))
			}
		}
		sort.Slice(others, func(i, j int) bool { return others[i] < others[j] })
		pos := 0
		for pos < len(others) {
			b.WriteByte(others[pos])
			pos++
			if freq[first] > 0 {
				b.WriteByte(byte('a' + first))
				freq[first]--
			}
		}
		for freq[first] > 0 {
			b.WriteByte(byte('a' + first))
			freq[first]--
		}
		return b.String()
	}
	second := first + 1
	for freq[second] == 0 {
		second++
	}
	if distinct == 2 {
		var b strings.Builder
		b.WriteByte(byte('a' + first))
		for i := 0; i < freq[second]; i++ {
			b.WriteByte(byte('a' + second))
		}
		for i := 0; i < freq[first]-1; i++ {
			b.WriteByte(byte('a' + first))
		}
		return b.String()
	}
	third := second + 1
	for freq[third] == 0 {
		third++
	}
	var b strings.Builder
	b.WriteByte(byte('a' + first))
	b.WriteByte(byte('a' + second))
	freq[first]--
	freq[second]--
	for freq[first] > 0 {
		b.WriteByte(byte('a' + first))
		freq[first]--
	}
	b.WriteByte(byte('a' + third))
	freq[third]--
	for i := 0; i < 26; i++ {
		for freq[i] > 0 {
			b.WriteByte(byte('a' + i))
			freq[i]--
		}
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsE()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.s)
	}

	expected := make([]string, len(tests))
	for i, tc := range tests {
		expected[i] = solveE(tc)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n%s\n", err, stderr.String())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	for i, exp := range expected {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
			os.Exit(1)
		}
		if scanner.Text() != exp {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
