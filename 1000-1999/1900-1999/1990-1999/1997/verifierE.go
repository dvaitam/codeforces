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

const N_E = 200001

func expectedAnswerE(n int, arr []int, queries [][2]int) []string {
	a := make([]int, n+1)
	copy(a[1:], arr)
	t := make([]int, N_E)
	req := make([]int, n+1)
	add := func(x int) {
		for x < N_E {
			t[x]++
			x += x & -x
		}
	}
	for i := 1; i <= n; i++ {
		x, y := 0, 0
		for j := 17; j >= 0; j-- {
			nxt := x | (1 << j)
			if nxt < N_E && int64(a[i])*int64(nxt) <= int64(y+t[nxt]) {
				x = nxt
				y += t[nxt]
			}
		}
		x++
		add(x)
		req[i] = x
	}
	res := make([]string, len(queries))
	for i, q := range queries {
		x, y := q[0], q[1]
		if y < req[x] {
			res[i] = "NO"
		} else {
			res[i] = "YES"
		}
	}
	return res
}

func generateCaseE(rng *rand.Rand) (int, []int, [][2]int) {
	n := rng.Intn(5) + 1
	q := rng.Intn(5) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(5) + 1
	}
	queries := make([][2]int, q)
	for i := range queries {
		queries[i][0] = rng.Intn(n) + 1
		queries[i][1] = rng.Intn(10) + 1
	}
	return n, arr, queries
}

func runCaseE(bin string, n int, arr []int, queries [][2]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(queries)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")
	for _, q := range queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	expected := expectedAnswerE(n, arr, queries)
	if len(gotLines) != len(expected) {
		return fmt.Errorf("expected %d lines got %d", len(expected), len(gotLines))
	}
	for i := range expected {
		if strings.TrimSpace(gotLines[i]) != expected[i] {
			return fmt.Errorf("line %d expected %s got %s", i+1, expected[i], gotLines[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, arr, queries := generateCaseE(rng)
		if err := runCaseE(bin, n, arr, queries); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
