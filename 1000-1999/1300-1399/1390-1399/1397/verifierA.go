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
	n    int
	strs []string
}

func solve(tc testCase) string {
	freq := make([]int, 26)
	for _, s := range tc.strs {
		for _, ch := range s {
			if ch >= 'a' && ch <= 'z' {
				freq[ch-'a']++
			}
		}
	}
	for _, f := range freq {
		if f%tc.n != 0 {
			return "NO"
		}
	}
	return "YES"
}

func (tc testCase) input() string {
	var b strings.Builder
	b.WriteString("1\n")
	b.WriteString(fmt.Sprintf("%d\n", tc.n))
	for _, s := range tc.strs {
		b.WriteString(s)
		b.WriteByte('\n')
	}
	return b.String()
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 1
	strs := make([]string, n)
	for i := 0; i < n; i++ {
		l := rng.Intn(6) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			b[j] = byte('a' + rng.Intn(26))
		}
		strs[i] = string(b)
	}
	return testCase{n: n, strs: strs}
}

func runProgram(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCase) error {
	in := tc.input()
	expected := solve(tc)
	got, err := runProgram(bin, in)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{{n: 3, strs: []string{"abc", "bca", "cab"}}}
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
