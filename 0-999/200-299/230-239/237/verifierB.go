package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input   string
	c       []int
	initial [][]int
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	c := make([]int, n)
	c[0] = rng.Intn(4) + 1
	for i := 1; i < n; i++ {
		c[i] = rng.Intn(c[i-1]) + 1
	}
	s := 0
	for _, v := range c {
		s += v
	}
	nums := rng.Perm(s)
	arr := make([][]int, n)
	idx := 0
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", c[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		arr[i] = make([]int, c[i])
		for j := 0; j < c[i]; j++ {
			val := nums[idx] + 1
			idx++
			arr[i][j] = val
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", val))
		}
		sb.WriteByte('\n')
	}
	return testCase{input: sb.String(), c: c, initial: arr}
}

func runCase(bin string, tc testCase) error {
	// make copy
	arr := make([][]int, len(tc.initial))
	for i := range arr {
		arr[i] = append([]int(nil), tc.initial[i]...)
	}
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("missing number of swaps")
	}
	m, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return fmt.Errorf("invalid number of swaps")
	}
	s := 0
	for _, v := range tc.c {
		s += v
	}
	if m < 0 || m > s {
		return fmt.Errorf("invalid swaps count")
	}
	for i := 0; i < m; i++ {
		vals := make([]int, 4)
		for j := 0; j < 4; j++ {
			if !scanner.Scan() {
				return fmt.Errorf("not enough values for swap %d", i+1)
			}
			v, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return fmt.Errorf("invalid integer in swap")
			}
			vals[j] = v
		}
		x, y, p, q := vals[0]-1, vals[1]-1, vals[2]-1, vals[3]-1
		if x < 0 || x >= len(arr) || p < 0 || p >= len(arr) ||
			y < 0 || y >= len(arr[x]) || q < 0 || q >= len(arr[p]) {
			return fmt.Errorf("indices out of range")
		}
		arr[x][y], arr[p][q] = arr[p][q], arr[x][y]
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output")
	}
	// verify numbers
	seen := make(map[int]bool)
	total := 0
	for i := range arr {
		for j := range arr[i] {
			val := arr[i][j]
			if val <= 0 || val > s {
				return fmt.Errorf("invalid value")
			}
			if seen[val] {
				return fmt.Errorf("duplicate value")
			}
			seen[val] = true
			total++
		}
	}
	if total != s {
		return fmt.Errorf("missing values")
	}
	// row condition
	for i := range arr {
		for j := 1; j < len(arr[i]); j++ {
			if arr[i][j] <= arr[i][j-1] {
				return fmt.Errorf("row not strictly increasing")
			}
		}
	}
	// column condition
	for i := 1; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			if arr[i][j] <= arr[i-1][j] {
				return fmt.Errorf("column not strictly increasing")
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	// simple deterministic case
	cases = append(cases, generateCase(rand.New(rand.NewSource(1))))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
