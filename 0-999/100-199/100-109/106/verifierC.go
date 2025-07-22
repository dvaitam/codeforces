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

type group struct{ a, b, c, d int }

func expected(n, m, c0, d0 int, g []group) int {
	dp := make([]int, n+1)
	for i := c0; i <= n; i++ {
		if dp[i-c0]+d0 > dp[i] {
			dp[i] = dp[i-c0] + d0
		}
	}
	for _, gr := range g {
		count := gr.a / gr.b
		for t := 0; t < count; t++ {
			for j := n; j >= gr.c; j-- {
				if dp[j-gr.c]+gr.d > dp[j] {
					dp[j] = dp[j-gr.c] + gr.d
				}
			}
		}
	}
	return dp[n]
}

func run(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		c0, _ := strconv.Atoi(parts[2])
		d0, _ := strconv.Atoi(parts[3])
		if len(parts) != 4+4*m {
			fmt.Fprintf(os.Stderr, "case %d: bad line\n", idx)
			os.Exit(1)
		}
		g := make([]group, m)
		p := 4
		for i := 0; i < m; i++ {
			a, _ := strconv.Atoi(parts[p])
			b, _ := strconv.Atoi(parts[p+1])
			c, _ := strconv.Atoi(parts[p+2])
			d, _ := strconv.Atoi(parts[p+3])
			g[i] = group{a, b, c, d}
			p += 4
		}
		expect := expected(n, m, c0, d0, g)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, m, c0, d0))
		for i := 0; i < m; i++ {
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", g[i].a, g[i].b, g[i].c, g[i].d))
		}
		gotStr, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		got, _ := strconv.Atoi(strings.Fields(gotStr)[0])
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
