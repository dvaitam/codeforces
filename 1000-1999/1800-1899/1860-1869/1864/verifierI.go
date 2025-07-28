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

func expectedI(n, q int, queries [][2]int) string {
	idx := make([][]int, n)
	cnt := 0
	for i := 0; i < n; i++ {
		idx[i] = make([]int, n)
		if i%2 == 0 {
			for j := 0; j < n; j++ {
				cnt++
				idx[i][j] = cnt
			}
		} else {
			for j := n - 1; j >= 0; j-- {
				cnt++
				idx[i][j] = cnt
			}
		}
	}
	prevAns := 0
	var sb strings.Builder
	for i, qu := range queries {
		x := (qu[0] ^ prevAns) - 1
		y := (qu[1] ^ prevAns) - 1
		ans := n*n - idx[x][y] + 1
		prevAns = ans
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", ans))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	q := rng.Intn(5) + 1
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		x := rng.Intn(n) + 1
		y := rng.Intn(n) + 1
		queries[i] = [2]int{x, y}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for _, qu := range queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", qu[0], qu[1]))
	}
	expect := expectedI(n, q, queries)
	return sb.String(), expect
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(expected)
	if outStr != expStr {
		return fmt.Errorf("expected %q got %q", expStr, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
