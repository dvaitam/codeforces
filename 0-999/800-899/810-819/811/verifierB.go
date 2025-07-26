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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(n, m int, p []int, queries [][3]int) []string {
	res := make([]string, m)
	for i, q := range queries {
		l := q[0] - 1
		r := q[1] - 1
		x := q[2] - 1
		target := p[x]
		cnt := 0
		for j := l; j <= r; j++ {
			if p[j] < target {
				cnt++
			}
		}
		if l+cnt == x {
			res[i] = "Yes"
		} else {
			res[i] = "No"
		}
	}
	return res
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	perm := rng.Perm(n)
	for i := range perm {
		perm[i]++
	}
	queries := make([][3]int, m)
	for i := 0; i < m; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		x := rng.Intn(r-l+1) + l
		queries[i] = [3]int{l, r, x}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i, v := range perm {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for _, q := range queries {
		fmt.Fprintf(&sb, "%d %d %d\n", q[0], q[1], q[2])
	}
	expected := strings.Join(solveCase(n, m, perm, queries), "\n")
	return sb.String(), expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
