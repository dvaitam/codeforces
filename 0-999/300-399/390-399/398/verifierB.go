package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solveB(n int, pairs [][2]int) float64 {
	a := make([]bool, n+1)
	b := make([]bool, n+1)
	for _, p := range pairs {
		if p[0] <= n {
			a[p[0]] = true
		}
		if p[1] <= n {
			b[p[1]] = true
		}
	}
	A, B := 0, 0
	for i := 1; i <= n; i++ {
		if a[i] {
			A++
		}
		if b[i] {
			B++
		}
	}
	f := make([][]float64, n+2)
	for i := range f {
		f[i] = make([]float64, n+2)
	}
	nn := float64(n)
	nn2 := nn * nn
	for i := n; i >= A; i-- {
		for j := n; j >= B; j-- {
			if i < n || j < n {
				denom := nn2 - float64(i*j)
				term1 := float64(n-i) * float64(j) * f[i+1][j]
				term2 := float64(i) * float64(n-j) * f[i][j+1]
				term3 := float64(n-i) * float64(n-j) * f[i+1][j+1]
				f[i][j] = (term1 + term2 + term3 + nn2) / denom
			}
		}
	}
	return f[A][B]
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("could not open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("bad test case %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		var pairs [][2]int
		for i := 0; i < m; i++ {
			x, _ := strconv.Atoi(fields[2+2*i])
			y, _ := strconv.Atoi(fields[3+2*i])
			pairs = append(pairs, [2]int{x, y})
		}
		exp := solveB(n, pairs)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, p := range pairs {
			input.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(string(out))
		got, err := strconv.ParseFloat(gotStr, 64)
		if err != nil {
			fmt.Printf("Test %d: output not a float\n", idx)
			os.Exit(1)
		}
		if abs(got-exp) > 1e-6 {
			fmt.Printf("Test %d failed: expected %.9f got %.9f\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
