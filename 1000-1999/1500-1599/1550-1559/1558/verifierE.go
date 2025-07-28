package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func canClear(n int, a, b []int64, g [][]int, start int64) bool {
	visited := make([]bool, n+1)
	visited[1] = true
	power := start
	changed := true
	for changed {
		changed = false
		for u := 1; u <= n; u++ {
			if !visited[u] {
				continue
			}
			for _, v := range g[u] {
				if visited[v] {
					continue
				}
				if power > a[v] {
					power += b[v]
					visited[v] = true
					changed = true
				}
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !visited[i] {
			return false
		}
	}
	return true
}

func solve(n, m int, a, b []int64, g [][]int) int64 {
	low, high := int64(1), int64(1)
	for i := 2; i <= n; i++ {
		if a[i]+1 > high {
			high = a[i] + 1
		}
	}
	var sum int64
	for i := 2; i <= n; i++ {
		sum += b[i]
	}
	if high < sum+1 {
		high = sum + 1
	}
	for low < high {
		mid := (low + high) / 2
		if canClear(n, a, b, g, mid) {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return low
}

func runCase(bin string, n, m int, a, b []int64, edges [][2]int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 2; i <= n; i++ {
		if i > 2 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	for i := 2; i <= n; i++ {
		if i > 2 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", b[i]))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	g := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	exp := fmt.Sprintf("%d", solve(n, m, a, b, g))
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Println("could not open testcasesE.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scanner.Text())
	for caseNum := 0; caseNum < t; caseNum++ {
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		m, _ := strconv.Atoi(scanner.Text())
		a := make([]int64, n+1)
		b := make([]int64, n+1)
		for i := 2; i <= n; i++ {
			if !scanner.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			val, _ := strconv.ParseInt(scanner.Text(), 10, 64)
			a[i] = val
		}
		for i := 2; i <= n; i++ {
			if !scanner.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			val, _ := strconv.ParseInt(scanner.Text(), 10, 64)
			b[i] = val
		}
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			if !scanner.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			u, _ := strconv.Atoi(scanner.Text())
			if !scanner.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			v, _ := strconv.Atoi(scanner.Text())
			edges[i] = [2]int{u, v}
		}
		if err := runCase(bin, n, m, a, b, edges); err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
