package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func pickUnique(rng *rand.Rand, max, count int) []int {
	m := make(map[int]struct{})
	for len(m) < count {
		v := rng.Intn(max) + 1
		m[v] = struct{}{}
	}
	res := make([]int, 0, count)
	for v := range m {
		res = append(res, v)
	}
	sort.Ints(res)
	return res
}

func expected(n, m int, rows, cols []int) string {
	emptyRow := make([]bool, n+2)
	for _, r := range rows {
		if r >= 1 && r <= n {
			emptyRow[r] = true
		}
	}
	rowSeg := 0
	for i := 1; i <= n; i++ {
		if !emptyRow[i] && (i == 1 || emptyRow[i-1]) {
			rowSeg++
		}
	}
	emptyCol := make([]bool, m+2)
	for _, c := range cols {
		if c >= 1 && c <= m {
			emptyCol[c] = true
		}
	}
	colSeg := 0
	for j := 1; j <= m; j++ {
		if !emptyCol[j] && (j == 1 || emptyCol[j-1]) {
			colSeg++
		}
	}
	return fmt.Sprintf("%d", rowSeg*colSeg)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	t := rng.Intn(n + 1)
	rows := pickUnique(rng, n, t)
	s := rng.Intn(m + 1)
	cols := pickUnique(rng, m, s)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	sb.WriteString(fmt.Sprintf("%d", t))
	for _, r := range rows {
		sb.WriteString(fmt.Sprintf(" %d", r))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%d", s))
	for _, c := range cols {
		sb.WriteString(fmt.Sprintf(" %d", c))
	}
	sb.WriteString("\n")
	input := sb.String()
	return input, expected(n, m, rows, cols)
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
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected '%s' got '%s'", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
