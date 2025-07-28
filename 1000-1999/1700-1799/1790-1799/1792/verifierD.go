package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

func compose(a, b []int) []int {
	m := len(a)
	r := make([]int, m)
	for i := 0; i < m; i++ {
		r[i] = b[a[i]-1]
	}
	return r
}

func beauty(p []int) int {
	k := 0
	for k < len(p) && p[k] == k+1 {
		k++
	}
	return k
}

func solveCaseD(perms [][]int) []int {
	n := len(perms)
	m := len(perms[0])
	res := make([]int, n)
	for i := 0; i < n; i++ {
		best := 0
		for j := 0; j < n; j++ {
			r := compose(perms[i], perms[j])
			b := beauty(r)
			if b > best {
				best = b
			}
		}
		res[i] = best
	}
	_ = m
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s /path/to/binary\n", os.Args[0])
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(45)
	const t = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	var expected [][]int
	for i := 0; i < t; i++ {
		n := rand.Intn(4) + 1
		m := rand.Intn(5) + 1
		fmt.Fprintf(&input, "%d %d\n", n, m)
		perms := make([][]int, n)
		for j := 0; j < n; j++ {
			perm := rand.Perm(m)
			for k := 0; k < m; k++ {
				perm[k]++
			}
			perms[j] = perm
			for k := 0; k < m; k++ {
				if k > 0 {
					input.WriteByte(' ')
				}
				fmt.Fprint(&input, perm[k])
			}
			input.WriteByte('\n')
		}
		expected = append(expected, solveCaseD(perms))
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	outBytes, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "binary execution failed:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(outBytes))
	scanner.Split(bufio.ScanWords)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		exp := expected[caseIdx]
		for j := 0; j < len(exp); j++ {
			if !scanner.Scan() {
				fmt.Printf("not enough output on case %d value %d\n", caseIdx+1, j+1)
				os.Exit(1)
			}
			got, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Printf("invalid integer on case %d value %d: %v\n", caseIdx+1, j+1, err)
				os.Exit(1)
			}
			if got != exp[j] {
				fmt.Printf("mismatch on case %d value %d: expected %d got %d\n", caseIdx+1, j+1, exp[j], got)
				os.Exit(1)
			}
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed.")
}
