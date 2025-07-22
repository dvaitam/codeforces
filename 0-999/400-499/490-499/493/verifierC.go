package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func solveC(r *bufio.Reader) string {
	var n, m int
	fmt.Fscan(r, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &a[i])
	}
	fmt.Fscan(r, &m)
	b := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(r, &b[i])
	}
	sort.Ints(a)
	sort.Ints(b)
	score1 := 3 * n
	score2 := 3 * m
	best1, best2 := score1, score2
	bestDiff := score1 - score2
	i, j := 0, 0
	for i < n || j < m {
		var d int
		if i < n && (j >= m || a[i] < b[j]) {
			d = a[i]
		} else if j < m && (i >= n || b[j] < a[i]) {
			d = b[j]
		} else {
			d = a[i]
		}
		k1 := 0
		for i < n && a[i] == d {
			k1++
			i++
		}
		k2 := 0
		for j < m && b[j] == d {
			k2++
			j++
		}
		if k1 > 0 {
			score1 -= k1
		}
		if k2 > 0 {
			score2 -= k2
		}
		diff := score1 - score2
		if diff > bestDiff || (diff == bestDiff && score1 > best1) {
			bestDiff = diff
			best1 = score1
			best2 = score2
		}
	}
	return fmt.Sprintf("%d:%d", best1, best2)
}

func generateCaseC(rng *rand.Rand) string {
	n := rng.Intn(6) + 1
	m := rng.Intn(6) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", rng.Intn(20)+1)
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", m)
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", rng.Intn(20)+1)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseC(rng)
		expect := solveC(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
