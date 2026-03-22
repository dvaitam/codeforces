package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func buildRef() (string, error) {
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1735F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return ref, nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	a := rng.Int63n(1000000000)
	b := rng.Int63n(1000000000)
	p := make([]int64, n)
	q := make([]int64, n)
	for i := 0; i < n; i++ {
		p[i] = rng.Int63n(1000000000) + 1
		q[i] = rng.Int63n(1000000000) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d %d\n", n, a, b)
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range q {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func genCases() []string {
	rng := rand.New(rand.NewSource(1735))
	cases := make([]string, 0, 100)
	cases = append(cases, "1\n1 0 0\n1\n1\n")
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}
	return cases
}

func floatsMatch(exp, got string) bool {
	expTokens := strings.Fields(exp)
	gotTokens := strings.Fields(got)
	if len(expTokens) != len(gotTokens) {
		return false
	}
	for i := range expTokens {
		a, errA := strconv.ParseFloat(expTokens[i], 64)
		b, errB := strconv.ParseFloat(gotTokens[i], 64)
		if errA != nil || errB != nil {
			if expTokens[i] != gotTokens[i] {
				return false
			}
			continue
		}
		diff := math.Abs(a - b)
		denom := math.Max(1.0, math.Abs(b))
		if diff/denom > 1e-6 {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	cases := genCases()
	for i, tc := range cases {
		exp, err := runBinary(ref, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if !floatsMatch(exp, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on case %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed!\n", len(cases))
}
