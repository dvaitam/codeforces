package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// DSU for expected calculation
type dsu struct{ p, sz []int }

func newDSU(n int) *dsu {
	d := &dsu{p: make([]int, n+1), sz: make([]int, n+1)}
	for i := 0; i <= n; i++ {
		d.p[i] = i
		d.sz[i] = 1
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) union(a, b int) bool {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return false
	}
	if d.sz[ra] < d.sz[rb] {
		ra, rb = rb, ra
	}
	d.p[rb] = ra
	d.sz[ra] += d.sz[rb]
	return true
}

func expectedFromLines(lines []string) (string, error) {
	if len(lines) == 0 {
		return "", fmt.Errorf("empty case")
	}
	h := strings.Fields(lines[0])
	if len(h) < 2 {
		return "", fmt.Errorf("bad header: %q", lines[0])
	}
	n, err := strconv.Atoi(h[0])
	if err != nil {
		return "", err
	}
	m, err := strconv.Atoi(h[1])
	if err != nil {
		return "", err
	}
	if len(lines) < 1+m {
		return "", fmt.Errorf("need %d edges, got %d", m, len(lines)-1)
	}
	d := newDSU(n)
	const MOD int64 = 1000000009
	pow := int64(1)
	var sb strings.Builder
	for i := 0; i < m; i++ {
		parts := strings.Fields(lines[1+i])
		if len(parts) < 2 {
			return "", fmt.Errorf("bad edge line: %q", lines[1+i])
		}
		u, _ := strconv.Atoi(parts[0])
		v, _ := strconv.Atoi(parts[1])
		if !d.union(u, v) {
			pow = (pow * 2) % MOD
		}
		ans := (pow - 1 + MOD) % MOD
		sb.WriteString(strconv.FormatInt(ans, 10))
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n"), nil
}

func runProg(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	var lines []string
	process := func() {
		if len(lines) == 0 {
			return
		}
		idx++
		input := strings.Join(lines, "\n") + "\n"
		exp, err := expectedFromLines(lines)
		if err != nil {
			fmt.Printf("verifier error on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d mismatch\nexpected: %s\n got: %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			process()
			lines = nil
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	process()
	fmt.Printf("All %d tests passed\n", idx)
}
