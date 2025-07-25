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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

func expected(m, r []int64) string {
	l := int64(1)
	for i := range m {
		l = lcm(l, m[i])
	}
	var cnt int64
	for x := int64(0); x < l; x++ {
		for j := range m {
			if x%m[j] == r[j] {
				cnt++
				break
			}
		}
	}
	res := float64(cnt) / float64(l)
	return fmt.Sprintf("%.10f", res)
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) ([]int64, []int64) {
	n := rng.Intn(5) + 1 // 1..5
	m := make([]int64, n)
	r := make([]int64, n)
	for i := 0; i < n; i++ {
		m[i] = int64(rng.Intn(16) + 1)
	}
	for i := 0; i < n; i++ {
		r[i] = int64(rng.Intn(int(m[i])))
	}
	return m, r
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	const cases = 100
	for i := 0; i < cases; i++ {
		m, r := genCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(m)))
		for j := range m {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", m[j]))
		}
		sb.WriteByte('\n')
		for j := range r {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", r[j]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		want := expected(m, r)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected %s\ngot %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
