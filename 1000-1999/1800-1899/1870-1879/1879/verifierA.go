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

type testCase struct {
	n int
	s []int64
	e []int64
}

func expected(tc testCase) string {
	s1 := tc.s[0]
	e1 := tc.e[0]
	maxS := int64(-1)
	for i := 1; i < tc.n; i++ {
		if tc.e[i] >= e1 && tc.s[i] > maxS {
			maxS = tc.s[i]
		}
	}
	if s1 > maxS {
		return fmt.Sprintf("%d", s1)
	}
	return "-1"
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 102)
	cases = append(cases, testCase{n: 2, s: []int64{1, 1}, e: []int64{1, 1}})
	cases = append(cases, testCase{n: 2, s: []int64{10, 1}, e: []int64{5, 4}})
	for len(cases) < 102 {
		n := rng.Intn(9) + 2
		s := make([]int64, n)
		e := make([]int64, n)
		for i := 0; i < n; i++ {
			s[i] = rng.Int63n(100) + 1
			e[i] = rng.Int63n(50) + 1
		}
		cases = append(cases, testCase{n: n, s: s, e: e})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genCases()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for j := 0; j < tc.n; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", tc.s[j], tc.e[j]))
		}
		input := sb.String()
		want := expected(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
