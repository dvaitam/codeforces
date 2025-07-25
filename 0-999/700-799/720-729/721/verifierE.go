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

type segment struct {
	l, r int64
}

type testCase struct {
	L    int64
	n    int
	p, t int64
	segs []segment
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

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.L, tc.n, tc.p, tc.t))
	for _, s := range tc.segs {
		sb.WriteString(fmt.Sprintf("%d %d\n", s.l, s.r))
	}
	return sb.String()
}

func expected(tc testCase) string {
	nextAvail := int64(0)
	ans := int64(0)
	for _, s := range tc.segs {
		y0 := s.l
		if nextAvail > y0 {
			y0 = nextAvail
		}
		if s.r-y0 >= tc.p {
			avail := s.r - y0
			k := avail / tc.p
			ans += k
			lastEnd := y0 + k*tc.p
			nextAvail = lastEnd + tc.t
		}
	}
	return fmt.Sprint(ans)
}

func randomCase(rng *rand.Rand) testCase {
	L := int64(rng.Intn(100) + 1)
	n := rng.Intn(5)
	p := int64(rng.Intn(10) + 1)
	t := int64(rng.Intn(5) + 1)
	segs := make([]segment, 0, n)
	cur := int64(0)
	for i := 0; i < n; i++ {
		gap := int64(rng.Intn(3) + 1)
		start := cur + gap
		if start >= L {
			break
		}
		length := int64(rng.Intn(10) + 1)
		end := start + length
		if end > L {
			end = L
		}
		segs = append(segs, segment{start, end})
		cur = end
	}
	return testCase{L: L, n: len(segs), p: p, t: t, segs: segs}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{}
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
