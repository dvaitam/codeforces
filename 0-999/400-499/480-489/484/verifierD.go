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

type node struct {
	val  int64
	idx  int
	best int64
}

func solve(a []int64) int64 {
	n := len(a)
	dp := make([]int64, n+1)
	stMax := make([]node, 0, n)
	stMin := make([]node, 0, n)
	for i := 1; i <= n; i++ {
		x := a[i-1]
		dpi := dp[i-1]
		pre := i
		for len(stMax) > 0 && stMax[len(stMax)-1].val <= x {
			top := stMax[len(stMax)-1]
			stMax = stMax[:len(stMax)-1]
			if top.best+x-top.val > dpi {
				dpi = top.best + x - top.val
			}
			pre = top.idx
		}
		stMax = append(stMax, node{val: x, idx: pre, best: dp[pre-1]})
		pre2 := i
		for len(stMin) > 0 && stMin[len(stMin)-1].val >= x {
			top := stMin[len(stMin)-1]
			stMin = stMin[:len(stMin)-1]
			if top.best+top.val-x > dpi {
				dpi = top.best + top.val - x
			}
			pre2 = top.idx
		}
		stMin = append(stMin, node{val: x, idx: pre2, best: dp[pre2-1]})
		dp[i] = dpi
	}
	return dp[n]
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
