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

func expectedAnswerD(n int, vals []int, parents []int) int {
	e := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		p := parents[i-2]
		e[p] = append(e[p], i)
	}
	var S func(u int) int
	S = func(u int) int {
		if len(e[u]) == 0 {
			return vals[u-1]
		}
		r := 1000000000
		for _, v := range e[u] {
			val := S(v)
			if val < r {
				r = val
			}
		}
		if vals[u-1] < r && u != 1 {
			return (r + vals[u-1]) / 2
		}
		return r
	}
	return vals[0] + S(1)
}

func generateCaseD(rng *rand.Rand) (int, []int, []int) {
	n := rng.Intn(6) + 2 // at least 2
	vals := make([]int, n)
	for i := range vals {
		vals[i] = rng.Intn(10) + 1
	}
	parents := make([]int, n-1)
	for i := 2; i <= n; i++ {
		parents[i-2] = rng.Intn(i-1) + 1
	}
	return n, vals, parents
}

func runCaseD(bin string, n int, vals []int, parents []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")
	for i, p := range parents {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(p))
	}
	sb.WriteString("\n")
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := fmt.Sprint(expectedAnswerD(n, vals, parents))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, vals, parents := generateCaseD(rng)
		if err := runCaseD(bin, n, vals, parents); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
