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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type query struct {
	op int //0 query,1 add
	k  int
	x  int64
}

func solveCase(n int, a []int64, sets [][]int, qs []query) string {
	arr := make([]int64, n)
	copy(arr, a)
	var sb strings.Builder
	for _, q := range qs {
		if q.op == 0 {
			var sum int64
			for _, idx := range sets[q.k] {
				sum += arr[idx]
			}
			sb.WriteString(fmt.Sprintf("%d\n", sum))
		} else {
			for _, idx := range sets[q.k] {
				arr[idx] += q.x
			}
		}
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	qn := rng.Intn(8) + 1
	a := make([]int64, n)
	for i := range a {
		a[i] = int64(rng.Intn(21) - 10)
	}
	sets := make([][]int, m)
	for i := 0; i < m; i++ {
		sz := rng.Intn(n) + 1
		perm := rng.Perm(n)[:sz]
		sets[i] = make([]int, sz)
		copy(sets[i], perm)
	}
	qs := make([]query, qn)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, qn))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", a[i]))
	}
	sb.WriteString("\n")
	for i := 0; i < m; i++ {
		sb.WriteString(fmt.Sprintf("%d", len(sets[i])))
		for _, idx := range sets[i] {
			sb.WriteString(fmt.Sprintf(" %d", idx+1))
		}
		sb.WriteString("\n")
	}
	for i := 0; i < qn; i++ {
		if rng.Intn(2) == 0 {
			k := rng.Intn(m)
			qs[i] = query{op: 0, k: k}
			sb.WriteString(fmt.Sprintf("? %d\n", k+1))
		} else {
			k := rng.Intn(m)
			x := int64(rng.Intn(21) - 10)
			qs[i] = query{op: 1, k: k, x: x}
			sb.WriteString(fmt.Sprintf("+ %d %d\n", k+1, x))
		}
	}
	expect := solveCase(n, a, sets, qs)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
