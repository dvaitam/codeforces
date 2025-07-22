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

func expected(n int, edges [][2]int) int64 {
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
	}
	cnt := make([]int, n+1)
	var ans int64
	for a := 1; a <= n; a++ {
		for _, b := range adj[a] {
			for _, c := range adj[b] {
				cnt[c]++
			}
		}
		for c := 1; c <= n; c++ {
			k := cnt[c]
			if k > 1 {
				ans += int64(k*(k-1)) / 2
			}
			cnt[c] = 0
		}
	}
	return ans
}

func parseInts(fields []string) ([]int, error) {
	vals := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, err
		}
		vals[i] = v
	}
	return vals, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
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
		values, err := parseInts(fields)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad test line %d: %v\n", idx, err)
			os.Exit(1)
		}
		if len(values) < 2 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		n := values[0]
		m := values[1]
		if len(values) != 2+2*m {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		edges := make([][2]int, m)
		p := 2
		for i := 0; i < m; i++ {
			edges[i] = [2]int{values[p], values[p+1]}
			p += 2
		}
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < m; i++ {
			b.WriteString(fmt.Sprintf("%d %d\n", edges[i][0], edges[i][1]))
		}
		input := b.String()
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		tokens := strings.Fields(out)
		if len(tokens) != 1 {
			fmt.Fprintf(os.Stderr, "case %d: invalid output\n", idx)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(tokens[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: bad integer\n", idx)
			os.Exit(1)
		}
		want := expected(n, edges)
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
