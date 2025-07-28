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

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	m := rng.Intn(6) + 2
	C := rng.Intn(10) + 1
	q := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, m, C, q))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rng.Intn(m)+1))
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func brute(input string) (string, error) {
	fields := strings.Fields(input)
	idx := 0
	n := atoi(fields[idx])
	idx++
	_ = atoi(fields[idx])
	idx++
	C := atoi(fields[idx])
	idx++
	q := atoi(fields[idx])
	idx++
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = atoi(fields[idx])
		idx++
	}
	outputs := make([]string, q)
	for i := 0; i < q; i++ {
		l := atoi(fields[idx])
		idx++
		r := atoi(fields[idx])
		idx++
		prod := 1
		for j := l - 1; j < r; j++ {
			prod *= a[j]
		}
		count := 0
		for x := 1; x <= C; x++ {
			if gcd(x, prod) == 1 {
				count++
			}
		}
		outputs[i] = fmt.Sprintf("%d", count)
	}
	return strings.Join(outputs, "\n"), nil
}

func atoi(s string) int {
	v := 0
	for i := 0; i < len(s); i++ {
		v = v*10 + int(s[i]-'0')
	}
	return v
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect, err := brute(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "brute error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
