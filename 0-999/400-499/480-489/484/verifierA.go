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

func solve(l, r uint64) uint64 {
	ans := r
	for k := 0; k < 63; k++ {
		if (ans>>k)&1 == 1 {
			mask := (uint64(1) << (k + 1)) - 1
			tmp := (ans & ^mask) | (mask >> 1)
			if tmp >= l {
				ans = tmp
			}
		}
	}
	return ans
}

func generateTests() []test {
	rnd := rand.New(rand.NewSource(1))
	tests := make([]test, 0, 100)
	for len(tests) < 100 {
		n := rnd.Intn(4) + 1
		var in bytes.Buffer
		fmt.Fprintln(&in, n)
		var out bytes.Buffer
		for i := 0; i < n; i++ {
			r := rnd.Uint64() % 1000000000000000000
			var l uint64
			if r > 0 {
				l = rnd.Uint64() % (r + 1)
			}
			fmt.Fprintf(&in, "%d %d\n", l, r)
			fmt.Fprintf(&out, "%d\n", solve(l, r))
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
		fmt.Fprintln(os.Stderr, "usage: verifierA.go <binary>")
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
