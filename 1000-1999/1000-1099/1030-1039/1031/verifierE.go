package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solve(a []int) (bool, [][3]int) {
	pos := make([]int, 0, len(a))
	for i, v := range a {
		if v == 1 {
			pos = append(pos, i)
		}
	}
	m := len(pos)
	if m%2 != 0 {
		return false, nil
	}
	ops := make([][3]int, 0, m)
	for i := 0; i < m/2; i++ {
		l := pos[i]
		r := pos[m-1-i]
		d := r - l
		if d < 2 {
			return false, nil
		}
		if d%2 == 0 {
			mid := (l + r) / 2
			ops = append(ops, [3]int{l + 1, mid + 1, r + 1})
		} else {
			mid1 := (l + r - 1) / 2
			mid2 := mid1 + 1
			if mid2 > r || mid1 <= l {
				return false, nil
			}
			ops = append(ops, [3]int{l + 1, mid1 + 1, mid2 + 1})
			ops = append(ops, [3]int{mid1 + 1, mid2 + 1, r + 1})
		}
	}
	return true, ops
}

func generateCase(rng *rand.Rand) (string, bool, [][3]int) {
	n := rng.Intn(20) + 3
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(2)
	}
	ok, ops := solve(a)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(a[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), ok, ops
}

func parseOutput(out string) (bool, [][3]int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return false, nil, fmt.Errorf("empty output")
	}
	t := strings.ToUpper(strings.TrimSpace(lines[0]))
	if t == "NO" {
		return false, nil, nil
	}
	if t != "YES" {
		return false, nil, fmt.Errorf("expected YES/NO")
	}
	if len(lines) < 2 {
		return false, nil, fmt.Errorf("missing count line")
	}
	m, err := strconv.Atoi(strings.TrimSpace(lines[1]))
	if err != nil {
		return false, nil, fmt.Errorf("bad count")
	}
	ops := make([][3]int, m)
	if len(lines) != 2+m {
		return false, nil, fmt.Errorf("wrong number of lines")
	}
	for i := 0; i < m; i++ {
		fields := strings.Fields(lines[2+i])
		if len(fields) != 3 {
			return false, nil, fmt.Errorf("bad op line")
		}
		x, _ := strconv.Atoi(fields[0])
		y, _ := strconv.Atoi(fields[1])
		z, _ := strconv.Atoi(fields[2])
		ops[i] = [3]int{x, y, z}
	}
	return true, ops, nil
}

func runCase(bin, input string, expectOk bool, expectOps [][3]int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, errBuf.String())
	}
	ok, ops, err := parseOutput(out.String())
	if err != nil {
		return err
	}
	if ok != expectOk {
		return fmt.Errorf("expected %v but got %v", expectOk, ok)
	}
	if !ok {
		return nil
	}
	if len(ops) != len(expectOps) {
		return fmt.Errorf("expected %d ops got %d", len(expectOps), len(ops))
	}
	for i := range ops {
		if ops[i] != expectOps[i] {
			return fmt.Errorf("op %d mismatch", i+1)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		in, ok, ops := generateCase(rng)
		if err := runCase(bin, in, ok, ops); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
