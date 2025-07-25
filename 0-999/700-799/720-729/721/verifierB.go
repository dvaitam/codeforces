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
	n         int
	k         int
	passwords []string
	target    string
}

func run(bin, input string) (string, error) {
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

func expected(tc testCase) string {
	L := len(tc.target)
	less, equal := 0, 0
	for _, p := range tc.passwords {
		if len(p) < L {
			less++
		} else if len(p) == L {
			equal++
		}
	}
	bestPos := less + 1
	worstPos := less + equal
	best := bestPos + (bestPos-1)/tc.k*5
	worst := worstPos + (worstPos-1)/tc.k*5
	return fmt.Sprintf("%d %d", best, worst)
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for _, p := range tc.passwords {
		sb.WriteString(p)
		sb.WriteByte('\n')
	}
	sb.WriteString(tc.target)
	sb.WriteByte('\n')
	return sb.String()
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func randStr(rng *rand.Rand, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(100) + 1
	k := rng.Intn(100) + 1
	passwords := make([]string, n)
	for i := 0; i < n; i++ {
		passwords[i] = randStr(rng, rng.Intn(10)+1)
	}
	target := passwords[rng.Intn(n)]
	return testCase{n: n, k: k, passwords: passwords, target: target}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		{1, 1, []string{"a"}, "a"},
		{2, 1, []string{"ab", "c"}, "ab"},
		{3, 2, []string{"a", "bb", "ccc"}, "bb"},
	}
	for len(cases) < 105 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		input := buildInput(tc)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		exp := expected(tc)
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
