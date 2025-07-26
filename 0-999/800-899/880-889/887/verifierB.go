package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveCase(cubes [][10]bool) int {
	n := len(cubes)
	isPossible := func(num int) bool {
		d1 := num % 10
		num /= 10
		if num == 0 {
			for i := 0; i < n; i++ {
				if cubes[i][d1] {
					return true
				}
			}
			return false
		}
		d2 := num % 10
		num /= 10
		if num == 0 {
			for i := 0; i < n; i++ {
				if !cubes[i][d2] {
					continue
				}
				for j := 0; j < n; j++ {
					if i == j {
						continue
					}
					if cubes[j][d1] {
						return true
					}
				}
			}
			return false
		}
		d3 := num % 10
		num /= 10
		if num > 0 || n < 3 {
			return false
		}
		for i := 0; i < n; i++ {
			if !cubes[i][d3] {
				continue
			}
			for j := 0; j < n; j++ {
				if j == i || !cubes[j][d2] {
					continue
				}
				for k := 0; k < n; k++ {
					if k == i || k == j {
						continue
					}
					if cubes[k][d1] {
						return true
					}
				}
			}
		}
		return false
	}

	ans := 0
	for i := 1; i <= 999; i++ {
		if isPossible(i) {
			ans = i
		} else {
			break
		}
	}
	return ans
}

func runCase(bin string, digits [][]int) error {
	n := len(digits)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	cubes := make([][10]bool, n)
	for i := 0; i < n; i++ {
		for j := 0; j < 6; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", digits[i][j])
			cubes[i][digits[i][j]] = true
		}
		sb.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	var got int
	fmt.Sscanf(gotStr, "%d", &got)
	exp := solveCase(cubes)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(2)
	cases := make([][][]int, 100)
	for i := range cases {
		n := rand.Intn(3) + 1
		digits := make([][]int, n)
		for j := 0; j < n; j++ {
			digits[j] = make([]int, 6)
			for k := 0; k < 6; k++ {
				digits[j][k] = rand.Intn(10)
			}
		}
		cases[i] = digits
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
