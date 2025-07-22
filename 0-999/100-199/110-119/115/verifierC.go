package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const MOD = 1000003

func solveC(n, m int, grid []string) int {
	total := 1
	for i := 0; i < n; i++ {
		ways := 0
		for p := 0; p < 2; p++ {
			ok := true
			for j := 0; j < m; j++ {
				c := grid[i][j]
				var exp byte
				if (j%2 == 0) == (p == 0) {
					exp = 'L'
				} else {
					exp = 'R'
				}
				if c == '1' || c == '3' {
					if exp != 'L' {
						ok = false
						break
					}
				} else if c == '2' || c == '4' {
					if exp != 'R' {
						ok = false
						break
					}
				}
			}
			if ok {
				ways++
			}
		}
		if ways == 0 {
			return 0
		}
		total = (total * ways) % MOD
	}
	for j := 0; j < m; j++ {
		ways := 0
		for p := 0; p < 2; p++ {
			ok := true
			for i := 0; i < n; i++ {
				c := grid[i][j]
				var exp byte
				if (i%2 == 0) == (p == 0) {
					exp = 'T'
				} else {
					exp = 'B'
				}
				if c == '1' || c == '2' {
					if exp != 'T' {
						ok = false
						break
					}
				} else if c == '3' || c == '4' {
					if exp != 'B' {
						ok = false
						break
					}
				}
			}
			if ok {
				ways++
			}
		}
		if ways == 0 {
			return 0
		}
		total = (total * ways) % MOD
	}
	return total
}

func genCase() (string, int) {
	n := rand.Intn(4) + 1
	m := rand.Intn(4) + 1
	grid := make([]string, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	chars := []byte{'1', '2', '3', '4', '.'}
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			row[j] = chars[rand.Intn(len(chars))]
		}
		grid[i] = string(row)
		sb.WriteString(grid[i])
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	expect := solveC(n, m, grid)
	return sb.String(), expect
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(42)
	for t := 1; t <= 100; t++ {
		in, expect := genCase()
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			fmt.Println(in)
			return
		}
		if strings.TrimSpace(got) != fmt.Sprint(expect) {
			fmt.Printf("Test %d failed\nInput:\n%s\nExpected: %d\nGot: %s\n", t, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
