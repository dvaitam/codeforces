package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solve(a []int64) int64 {
	var sum int64
	for _, v := range a {
		sum += v
	}
	avg := sum / int64(len(a))
	var ans int64
	var pref int64
	for i := 0; i < len(a)-1; i++ {
		pref += a[i] - avg
		if pref < 0 {
			ans -= pref
		} else {
			ans += pref
		}
	}
	return ans
}

type testCase struct {
	in  string
	out string
}

func generate() []testCase {
	const T = 100
	rand.Seed(2)
	cases := make([]testCase, T)
	for t := 0; t < T; t++ {
		n := rand.Intn(50) + 1
		a := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			a[i] = int64(rand.Intn(1000))
			sum += a[i]
		}
		mod := sum % int64(n)
		if mod != 0 {
			a[0] += int64(n) - mod
			sum += int64(n) - mod
		}
		var in strings.Builder
		fmt.Fprintf(&in, "%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				in.WriteByte(' ')
			}
			fmt.Fprintf(&in, "%d", a[i])
		}
		in.WriteByte('\n')
		cases[t] = testCase{
			in:  in.String(),
			out: fmt.Sprintf("%d\n", solve(a)),
		}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generate()
	for idx, tc := range cases {
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(tc.in)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != strings.TrimSpace(tc.out) {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, tc.out, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
