package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func validateTournament(n int, out string) bool {
	out = strings.TrimSpace(out)
	if n < 3 {
		return out == "-1"
	}
	fields := strings.Fields(out)
	if len(fields) != n*n {
		return false
	}
	a := make([][]int, n)
	idx := 0
	for i := 0; i < n; i++ {
		row := make([]int, n)
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(fields[idx])
			if err != nil || (v != 0 && v != 1) {
				return false
			}
			row[j] = v
			idx++
		}
		a[i] = row
	}
	for i := 0; i < n; i++ {
		if a[i][i] != 0 {
			return false
		}
		for j := i + 1; j < n; j++ {
			if a[i][j]+a[j][i] != 1 {
				return false
			}
		}
	}
	for s := 0; s < n; s++ {
		for t := 0; t < n; t++ {
			if s == t {
				continue
			}
			if a[s][t] == 1 {
				continue
			}
			ok := false
			for k := 0; k < n && !ok; k++ {
				if a[s][k] == 1 && a[k][t] == 1 {
					ok = true
				}
			}
			if !ok {
				return false
			}
		}
	}
	return true
}

func runCase(bin string, n int) error {
	in := fmt.Sprintf("%d\n", n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if !validateTournament(n, buf.String()) {
		return fmt.Errorf("wrong answer")
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for n := 1; n <= 100; n++ {
		if err := runCase(bin, n); err != nil {
			fmt.Printf("test %d failed: %v\n", n, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
