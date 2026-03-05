package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func maxSum(n int) int64 {
	var s int64
	for i := 1; i <= n; i++ {
		s += int64(2*i-1) * int64(i)
	}
	return s
}

// checkOutput verifies the solution output for a single test case with n.
// It reads from the scanner and returns an error description or "".
func checkOutput(sc *bufio.Scanner, n int) string {
	readLine := func() (string, bool) {
		if sc.Scan() {
			return strings.TrimSpace(sc.Text()), true
		}
		return "", false
	}

	line, ok := readLine()
	if !ok {
		return "unexpected EOF reading s m"
	}
	var s int64
	var m int
	if _, err := fmt.Sscan(line, &s, &m); err != nil {
		return fmt.Sprintf("bad s m line %q: %v", line, err)
	}
	expected := maxSum(n)
	if s != expected {
		return fmt.Sprintf("wrong sum: got %d, want %d", s, expected)
	}
	if m < 0 || m > 2*n {
		return fmt.Sprintf("m=%d out of range [0, %d]", m, 2*n)
	}

	// simulate operations on n×n matrix
	mat := make([][]int, n)
	for i := range mat {
		mat[i] = make([]int, n)
	}
	for op := 0; op < m; op++ {
		line, ok = readLine()
		if !ok {
			return fmt.Sprintf("unexpected EOF reading operation %d", op+1)
		}
		fields := strings.Fields(line)
		if len(fields) != 2+n {
			return fmt.Sprintf("operation %d: expected %d fields, got %d", op+1, 2+n, len(fields))
		}
		var c, idx int
		if _, err := fmt.Sscan(fields[0], &c); err != nil || (c != 1 && c != 2) {
			return fmt.Sprintf("operation %d: bad type %q", op+1, fields[0])
		}
		if _, err := fmt.Sscan(fields[1], &idx); err != nil || idx < 1 || idx > n {
			return fmt.Sprintf("operation %d: bad index %q", op+1, fields[1])
		}
		perm := make([]int, n)
		seen := make([]bool, n+1)
		for j := 0; j < n; j++ {
			if _, err := fmt.Sscan(fields[2+j], &perm[j]); err != nil || perm[j] < 1 || perm[j] > n || seen[perm[j]] {
				return fmt.Sprintf("operation %d: invalid permutation", op+1)
			}
			seen[perm[j]] = true
		}
		if c == 1 {
			for j := 0; j < n; j++ {
				mat[idx-1][j] = perm[j]
			}
		} else {
			for j := 0; j < n; j++ {
				mat[j][idx-1] = perm[j]
			}
		}
	}

	var total int64
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			total += int64(mat[i][j])
		}
	}
	if total != expected {
		return fmt.Sprintf("simulated matrix sum=%d, want %d", total, expected)
	}
	return ""
}

func genInput(rng *rand.Rand) string {
	t := rng.Intn(3) + 1
	var in bytes.Buffer
	fmt.Fprintf(&in, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(5) + 1
		fmt.Fprintf(&in, "%d\n", n)
	}
	return in.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		input := genInput(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s\n", i+1, err, got)
			os.Exit(1)
		}

		// parse t from input to know how many cases to check
		var t int
		fmt.Sscan(input, &t)
		ns := make([]int, t)
		sc := bufio.NewScanner(strings.NewReader(input))
		sc.Scan() // skip t line
		for j := 0; j < t; j++ {
			sc.Scan()
			fmt.Sscan(sc.Text(), &ns[j])
		}

		outSc := bufio.NewScanner(strings.NewReader(got))
		outSc.Buffer(make([]byte, 1<<20), 1<<20)
		for j := 0; j < t; j++ {
			if msg := checkOutput(outSc, ns[j]); msg != "" {
				fmt.Fprintf(os.Stderr, "test %d case %d (n=%d) failed: %s\ninput:\n%s\noutput:\n%s\n",
					i+1, j+1, ns[j], msg, input, got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
