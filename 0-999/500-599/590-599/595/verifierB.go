package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const MOD int64 = 1000000007

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

type testCase struct {
	input  string
	expect int64
}

func pow10(k int) int64 {
	p := int64(1)
	for i := 0; i < k; i++ {
		p *= 10
	}
	return p
}

func countMultiples(l, r, a int64) int64 {
	if l > r {
		return 0
	}
	start := ((l + a - 1) / a) * a
	if start > r {
		return 0
	}
	return (r-start)/a + 1
}

func expectedAnswer(a, b []int64, k int) int64 {
	p10 := pow10(k)
	base := p10 / 10
	p10m1 := p10 - 1
	ans := int64(1)
	for i := range a {
		total := p10m1/a[i] + 1
		l := int64(b[i]) * base
		r := l + base - 1
		bad := countMultiples(l, r, a[i])
		good := (total - bad) % MOD
		if good < 0 {
			good += MOD
		}
		ans = (ans * good) % MOD
	}
	return ans % MOD
}

func generateCase(rng *rand.Rand) testCase {
	k := rng.Intn(9) + 1
	m := rng.Intn(5) + 1
	n := k * m
	a := make([]int64, m)
	b := make([]int64, m)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	p10 := pow10(k)
	for i := 0; i < m; i++ {
		a[i] = int64(rng.Intn(int(p10-1)) + 1)
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		b[i] = int64(rng.Intn(10))
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", b[i]))
	}
	sb.WriteByte('\n')
	expect := expectedAnswer(a, b, k)
	return testCase{input: sb.String(), expect: expect}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	// simple deterministic case
	cases = append(cases, testCase{input: "1 1\n1\n0\n", expect: 9})
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != 1 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected single integer output, got %q\ninput:\n%s", i+1, out, tc.input)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid integer output\ninput:\n%soutput:\n%s", i+1, tc.input, out)
			os.Exit(1)
		}
		if val%MOD != tc.expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d, got %d\ninput:\n%s", i+1, tc.expect, val%MOD, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
