package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// solve mirrors the C++ reference solution:
//   mp1[v] = last occurrence index of v (1-based)
//   num    = count of distinct values seen so far (left to right)
//   ans   += num whenever position i is the last occurrence of a[i]
func solve(a []int) int64 {
	lastOcc := make(map[int]int)
	for i, v := range a {
		lastOcc[v] = i + 1 // 1-based
	}
	seen := make(map[int]bool)
	num := 0
	var ans int64
	for i, v := range a {
		if !seen[v] {
			seen[v] = true
			num++
		}
		if lastOcc[v] == i+1 {
			ans += int64(num)
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// genBatch builds t test cases into one input string and returns expected answers.
func genBatch(rng *rand.Rand, t int) (string, []int64) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	expected := make([]int64, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(15) + 1
		a := make([]int, n)
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(n) + 1 // values in [1, n] to produce repetitions
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(a[j]))
		}
		sb.WriteByte('\n')
		expected[i] = solve(a)
	}
	return sb.String(), expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	total := 0
	for batch := 0; batch < 20; batch++ {
		t := rng.Intn(10) + 1
		input, expected := genBatch(rng, t)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "batch %d failed: %v\ninput:\n%s", batch+1, err, input)
			os.Exit(1)
		}
		lines := strings.Fields(got)
		if len(lines) != t {
			fmt.Fprintf(os.Stderr, "batch %d: expected %d output lines, got %d\ninput:\n%s\noutput:\n%s\n",
				batch+1, t, len(lines), input, got)
			os.Exit(1)
		}
		for i, expVal := range expected {
			total++
			gotVal, err := strconv.ParseInt(lines[i], 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d: could not parse output %q: %v\n", total, lines[i], err)
				os.Exit(1)
			}
			if gotVal != expVal {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s\n",
					total, expVal, gotVal, input)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", total)
}
