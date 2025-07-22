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
	in  string
	out string
}

func applySort(s []byte, k, d int) []byte {
	n := len(s)
	tmp := make([]byte, n)
	for i := 0; i <= n-k; i++ {
		p := 0
		for g := 0; g < d; g++ {
			for j := i + g; j < i+k; j += d {
				tmp[p] = s[j]
				p++
			}
		}
		for t := 0; t < p; t++ {
			s[i+t] = tmp[t]
		}
	}
	return s
}

func generateTests() []test {
	rnd := rand.New(rand.NewSource(3))
	tests := make([]test, 0, 100)
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	for len(tests) < 100 {
		n := rnd.Intn(10) + 1
		s := make([]byte, n)
		for i := 0; i < n; i++ {
			s[i] = letters[rnd.Intn(len(letters))]
		}
		m := rnd.Intn(3) + 1
		var in bytes.Buffer
		fmt.Fprintln(&in, string(s))
		fmt.Fprintln(&in, m)
		var out bytes.Buffer
		S := make([]byte, n)
		copy(S, s)
		for q := 0; q < m; q++ {
			k := rnd.Intn(n) + 1
			d := rnd.Intn(k) + 1
			fmt.Fprintf(&in, "%d %d\n", k, d)
			S = applySort(S, k, d)
			fmt.Fprintln(&out, string(S))
		}
		tests = append(tests, test{in: in.String(), out: out.String()})
	}
	return tests
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierC.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		g := strings.TrimSpace(got)
		e := strings.TrimSpace(t.out)
		if g != e {
			fmt.Fprintf(os.Stderr, "Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, t.in, e, g)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
