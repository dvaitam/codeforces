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

func expected(n int, values []int, edges [][2]int) string {
	degree := make([]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		degree[u]++
		degree[v]++
	}
	leaves := 0
	for i := 2; i <= n; i++ {
		if degree[i] == 1 {
			leaves++
		}
	}
	return fmt.Sprintf("%d", leaves)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		values := make([]int, n+1)
		for j := 1; j <= n; j++ {
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			values[j] = v
		}
		edges := make([][2]int, n-1)
		for j := 0; j < n-1; j++ {
			scan.Scan()
			u, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			edges[j] = [2]int{u, v}
		}
		exp := expected(n, values, edges)
		var sb bytes.Buffer
		fmt.Fprintln(&sb, n)
		for i := 1; i <= n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(values[i]))
		}
		sb.WriteByte('\n')
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = bytes.NewReader(sb.Bytes())
		res, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", caseIdx+1, err, res)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(res))
		if got != exp {
			fmt.Printf("case %d failed: expected %s got %s\n", caseIdx+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
