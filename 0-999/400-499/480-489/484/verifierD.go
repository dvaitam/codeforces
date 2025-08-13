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

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func solve(a []int64) int64 {
	var ans, u, v int64
	const negInf int64 = -1e18
	u, v = negInf, negInf
	for _, x := range a {
		if ans+x > u {
			u = ans + x
		}
		if ans-x > v {
			v = ans - x
		}
		ans = max(u-x, v+x)
	}
	return ans
}

func generateTests() []test {
	rnd := rand.New(rand.NewSource(4))
	tests := make([]test, 0, 100)
	for len(tests) < 100 {
		n := rnd.Intn(20) + 1
		a := make([]int64, n)
		var in bytes.Buffer
		fmt.Fprintln(&in, n)
		for i := 0; i < n; i++ {
			v := rnd.Intn(101) - 50
			a[i] = int64(v)
			if i > 0 {
				fmt.Fprint(&in, " ")
			}
			fmt.Fprint(&in, v)
		}
		in.WriteByte('\n')
		out := fmt.Sprintf("%d\n", solve(a))
		tests = append(tests, test{in: in.String(), out: out})
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
		fmt.Fprintln(os.Stderr, "usage: verifierD.go <binary>")
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
