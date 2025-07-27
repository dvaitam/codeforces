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

func solve(n int, p []int64, h []int64, edges [][2]int) bool {
	adj := make([][]int, n)
	for _, e := range edges {
		x, y := e[0], e[1]
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}
	parent := make([]int, n)
	for i := range parent {
		parent[i] = -2
	}
	order := make([]int, 0, n)
	stack := []int{0}
	parent[0] = -1
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			stack = append(stack, v)
		}
	}
	sumPop := make([]int64, n)
	good := make([]int64, n)
	ok := true
	for i := len(order) - 1; i >= 0 && ok; i-- {
		u := order[i]
		sum := p[u]
		sumGood := int64(0)
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			sum += sumPop[v]
			sumGood += good[v]
		}
		sumPop[u] = sum
		if (sum+h[u])&1 != 0 {
			ok = false
			break
		}
		g := (sum + h[u]) / 2
		if g < 0 || g > sum {
			ok = false
			break
		}
		if sumGood > g {
			ok = false
			break
		}
		if g-sumGood > p[u] {
			ok = false
			break
		}
		good[u] = g
	}
	return ok
}

func runCase(bin string, n int, m int64, p []int64, h []int64, edges [][2]int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(p[i], 10))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(h[i], 10))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
	}
	input := sb.String()
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	exp := "NO"
	if solve(n, p, h, edges) {
		exp = "YES"
	}
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Println("could not open testcasesC.txt:", err)
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
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		mVal, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		p := make([]int64, n)
		for i := 0; i < n; i++ {
			if !scanner.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			v, _ := strconv.ParseInt(scanner.Text(), 10, 64)
			p[i] = v
		}
		h := make([]int64, n)
		for i := 0; i < n; i++ {
			if !scanner.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			v, _ := strconv.ParseInt(scanner.Text(), 10, 64)
			h[i] = v
		}
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			if !scanner.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			x, _ := strconv.Atoi(scanner.Text())
			if !scanner.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			y, _ := strconv.Atoi(scanner.Text())
			edges[i] = [2]int{x - 1, y - 1}
		}
		if err := runCase(bin, n, mVal, p, h, edges); err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
