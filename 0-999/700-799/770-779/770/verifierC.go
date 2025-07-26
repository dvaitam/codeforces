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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

var (
	g     [][]int
	state []int
	order []int
	cycle bool
)

func dfs(v int) {
	state[v] = 2
	for _, to := range g[v] {
		if cycle {
			return
		}
		if state[to] == 0 {
			dfs(to)
		} else if state[to] == 2 {
			cycle = true
			return
		}
	}
	state[v] = 1
	order = append(order, v)
}

func expectedOutput(n, k int, mains []int, graph [][]int) string {
	g = graph
	state = make([]int, n+1)
	order = nil
	cycle = false
	for _, x := range mains {
		if state[x] == 0 {
			dfs(x)
		}
	}
	if cycle {
		return "-1"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(order)))
	for _, v := range order {
		sb.WriteString(fmt.Sprintf("%d ", v))
	}
	return strings.TrimSpace(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("bad test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		mains := make([]int, k)
		for i := 0; i < k; i++ {
			scan.Scan()
			mains[i], _ = strconv.Atoi(scan.Text())
		}
		graph := make([][]int, n+1)
		for i := 1; i <= n; i++ {
			scan.Scan()
			ti, _ := strconv.Atoi(scan.Text())
			if ti > 0 {
				graph[i] = make([]int, ti)
				for j := 0; j < ti; j++ {
					scan.Scan()
					v, _ := strconv.Atoi(scan.Text())
					graph[i][j] = v
				}
			}
		}
		expect := expectedOutput(n, k, mains, graph)
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", n, k)
		for i := 0; i < k; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", mains[i])
		}
		input.WriteByte('\n')
		for i := 1; i <= n; i++ {
			fmt.Fprintf(&input, "%d", len(graph[i]))
			for _, v := range graph[i] {
				fmt.Fprintf(&input, " %d", v)
			}
			input.WriteByte('\n')
		}
		out, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", caseNum, expect, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
