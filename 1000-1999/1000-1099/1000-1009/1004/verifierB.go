package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type seg struct{ l, r int }
	type test struct {
		n, m int
		segs []seg
	}
	var cases []test
	// deterministic edge cases
	cases = append(cases, test{1, 1, []seg{{1, 1}}})
	cases = append(cases, test{2, 1, []seg{{1, 2}}})
	// random cases
	for i := 0; i < 98; i++ {
		n := rng.Intn(20) + 1
		m := rng.Intn(20) + 1
		segs := make([]seg, m)
		for j := 0; j < m; j++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			segs[j] = seg{l, r}
		}
		cases = append(cases, test{n, m, segs})
	}
	for idx, tc := range cases {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.m)
		for _, s := range tc.segs {
			input += fmt.Sprintf("%d %d\n", s.l, s.r)
		}
		want0 := make([]byte, tc.n)
		want1 := make([]byte, tc.n)
		for i := 0; i < tc.n; i++ {
			if i%2 == 0 {
				want0[i] = '0'
				want1[i] = '1'
			} else {
				want0[i] = '1'
				want1[i] = '0'
			}
		}
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		out := strings.TrimSpace(got)
		if out != string(want0) && out != string(want1) {
			fmt.Fprintf(os.Stderr, "case %d failed: output %q not alternating\n", idx+1, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
