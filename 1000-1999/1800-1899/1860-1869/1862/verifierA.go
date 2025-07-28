package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveA(n, m int, grid []string) string {
	letters := "vika"
	pos := 0
	for col := 0; col < m && pos < 4; col++ {
		for row := 0; row < n; row++ {
			if grid[row][col] == letters[pos] {
				pos++
				break
			}
		}
	}
	if pos == 4 {
		return "YES"
	}
	return "NO"
}

func genCases() []string {
	rand.Seed(1)
	cases := make([]string, 100)
	for idx := 0; idx < 100; idx++ {
		n := rand.Intn(5) + 1
		m := rand.Intn(8) + 4
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			row := make([]byte, m)
			for j := 0; j < m; j++ {
				row[j] = byte('a' + rand.Intn(26))
			}
			grid[i] = string(row)
		}
		sb := strings.Builder{}
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			sb.WriteString(grid[i])
			sb.WriteByte('\n')
		}
		cases[idx] = sb.String()
	}
	return cases
}

func runCase(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		lines := strings.Split(strings.TrimSpace(tc), "\n")
		var n, m int
		fmt.Sscan(lines[1], &n, &m)
		grid := lines[2 : 2+n]
		want := solveA(n, m, grid)
		got, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(strings.ToUpper(got)) != want {
			fmt.Fprintf(os.Stderr, "Wrong answer on case %d\nInput:\n%sExpected: %s Got: %s\n", i+1, tc, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
