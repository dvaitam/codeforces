package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `2 5 2
6 7 8 3 2 2 1
4 9 5 1 4
5 9 6 5 3 2
3 4 1 5
3 4 3 5
3 6 2 6
6 7 9 4 3 4 8
3 2 9 5
1 5
5 5 9 4 7 7
5 5 7 8 3 4
3 5 1 2
1 8
6 5 9 9 8 6 3
6 4 2 7 4 8 5
2 6 7
6 6 9 4 6 2 1
6 4 5 4 2 6 3
3 8 1 1
3 2 5 6
1 6
3 6 3 7
5 2 5 4 8 5
2 5 7
5 3 6 1 6 1
4 3 6 6 5
5 2 8 4 7 4
1 1
1 1
6 3 3 1 9 8 4
3 1 2 9
3 7 4 8
2 4 8
4 8 1 4 7
4 4 7 4 8
2 1 1
3 5 4 9
2 4 7
3 3 6 1
3 2 7 1
4 7 2 7 4
5 3 6 5 8 6
4 9 4 5 6
4 8 2 5 4
1 7
5 3 5 1 3 8
5 8 7 7 4 1
2 3 1
5 5 2 7 7 4
5 1 4 3 6 9
4 9 8 1 2
1 2
4 9 5 3 1
3 2 9 1
3 6 2 2
5 8 7 4 5 7
2 8 7
1 2
1 6
5 7 7 8 2 4
6 5 8 7 2 9 3
3 3 3 3
3 8 6 5
5 1 3 1 5 2
5 2 8 8 9 2
5 4 7 5 6 4
2 1 1
5 6 9 8 5 9
4 8 7 3 5
5 6 6 3 7 2
5 3 3 5 6 4
5 6 2 2 7 3
3 6 6 3
3 1 1 9
1 6
1 3
2 8 2
1 3
6 8 4 5 7 4 8
6 4 5 6 4 6 9
6 9 8 7 9 7 6
3 8 7 1
3 3 9 8
6 9 6 7 7 1 3
5 2 8 6 5 2
3 8 4 8
5 2 3 4 2 5
2 1 3
4 9 9 5 1
5 4 3 2 4 1
3 2 2 5
1 5
5 3 3 7 3 2
5 6 1 9 3 9
4 3 4 5 8
5 2 7 3 3 5
5 7 9 5 7 6
2 7 1
4 1 5 1 5`

type pair struct{ u, v int }

func validate(arr []int, output string) error {
	n := len(arr)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	m, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("invalid m: %v", err)
	}
	if m != len(lines)-1 {
		return fmt.Errorf("declared m=%d but got %d pairs", m, len(lines)-1)
	}

	// Count original inversions.
	origInv := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if arr[i] > arr[j] {
				origInv++
			}
		}
	}
	if m != origInv {
		return fmt.Errorf("expected %d inversions, got %d", origInv, m)
	}

	pairs := make([]pair, m)
	pairSet := make(map[[2]int]bool)
	for i := 0; i < m; i++ {
		parts := strings.Fields(lines[i+1])
		if len(parts) != 2 {
			return fmt.Errorf("pair %d: expected 2 numbers", i+1)
		}
		u, e1 := strconv.Atoi(parts[0])
		v, e2 := strconv.Atoi(parts[1])
		if e1 != nil || e2 != nil || u < 1 || v > n || u >= v {
			return fmt.Errorf("pair %d: invalid indices %d %d", i+1, u, v)
		}
		if arr[u-1] <= arr[v-1] {
			return fmt.Errorf("pair %d: (%d,%d) is not an inversion in original array", i+1, u, v)
		}
		key := [2]int{u, v}
		if pairSet[key] {
			return fmt.Errorf("pair %d: duplicate inversion (%d,%d)", i+1, u, v)
		}
		pairSet[key] = true
		pairs[i] = pair{u - 1, v - 1}
	}

	// Simulate swaps and verify the array is sorted.
	sim := make([]int, n)
	copy(sim, arr)
	for _, p := range pairs {
		sim[p.u], sim[p.v] = sim[p.v], sim[p.u]
	}
	for i := 1; i < n; i++ {
		if sim[i] < sim[i-1] {
			return fmt.Errorf("array not sorted after swaps")
		}
	}
	return nil
}

func parseTestcases(raw string) ([][]int, error) {
	lines := strings.Split(raw, "\n")
	var tests [][]int
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid n", idx+1)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d: expected %d values, got %d", idx+1, n, len(fields)-1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[i+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid value", idx+1)
			}
			arr[i] = v
		}
		tests = append(tests, arr)
	}
	return tests, nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid test data:", err)
		os.Exit(1)
	}

	for idx, arr := range tests {
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", len(arr)))
		for i, v := range arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := validate(arr, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %s\n", idx+1, err, input.String())
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
