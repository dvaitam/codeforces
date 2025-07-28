package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func properSubset(a, b []int) bool {
	if len(a) >= len(b) {
		return false
	}
	mb := make(map[int]bool, len(b))
	for _, v := range b {
		mb[v] = true
	}
	for _, v := range a {
		if !mb[v] {
			return false
		}
	}
	return true
}

func genCase(rng *rand.Rand) (string, []string) {
	n := rng.Intn(4) + 2 // 2..5
	sets := make([][]int, n)
	used := map[string]bool{}
	for i := 0; i < n; i++ {
		for {
			size := rng.Intn(n) + 1
			perm := rng.Perm(n)
			arr := make([]int, size)
			for j := 0; j < size; j++ {
				arr[j] = perm[j] + 1
			}
			sort.Ints(arr)
			key := fmt.Sprint(arr)
			if !used[key] {
				used[key] = true
				sets[i] = arr
				break
			}
		}
	}
	matrix := make([]string, n)
	for i := 0; i < n; i++ {
		var sb strings.Builder
		for j := 0; j < n; j++ {
			if properSubset(sets[i], sets[j]) {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
		}
		matrix[i] = sb.String()
	}
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		input.WriteString(matrix[i])
		input.WriteByte('\n')
	}
	return input.String(), matrix
}

func parseSets(out string, n int) ([][]int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != n {
		return nil, fmt.Errorf("expected %d lines, got %d", n, len(lines))
	}
	sets := make([][]int, n)
	seen := map[string]bool{}
	for i, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 1 {
			return nil, fmt.Errorf("line %d empty", i+1)
		}
		sz, err := strconv.Atoi(fields[0])
		if err != nil || sz <= 0 {
			return nil, fmt.Errorf("line %d invalid size", i+1)
		}
		if len(fields)-1 != sz {
			return nil, fmt.Errorf("line %d wrong number of elements", i+1)
		}
		set := make([]int, sz)
		mp := map[int]bool{}
		for j := 0; j < sz; j++ {
			v, err := strconv.Atoi(fields[j+1])
			if err != nil || v < 1 || v > n {
				return nil, fmt.Errorf("line %d invalid element", i+1)
			}
			if mp[v] {
				return nil, fmt.Errorf("line %d duplicate element", i+1)
			}
			mp[v] = true
			set[j] = v
		}
		sort.Ints(set)
		key := fmt.Sprint(set)
		if seen[key] {
			return nil, fmt.Errorf("duplicate set")
		}
		seen[key] = true
		sets[i] = set
	}
	return sets, nil
}

func verify(sets [][]int, matrix []string) error {
	n := len(sets)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			exp := properSubset(sets[i], sets[j])
			if (matrix[i][j] == '1') != exp {
				return fmt.Errorf("matrix mismatch for (%d,%d)", i+1, j+1)
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input, matrix := genCase(rng)
		out, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		sets, err := parseSets(out, len(matrix))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\noutput:\n%s\ninput:\n%s", t+1, err, out, input)
			os.Exit(1)
		}
		if err := verify(sets, matrix); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\noutput:\n%s\ninput:\n%s", t+1, err, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
