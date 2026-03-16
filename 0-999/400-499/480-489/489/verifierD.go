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
const testcasesDRaw = `
3 4 3 1 1 3 3 2 2 3
5 20 3 4 4 3 3 1 5 4 5 1 2 5 1 3 4 2 4 5 5 3 2 4 1 2 2 1 1 5 3 2 4 1 3 5 5 2 1 4 2 3
3 0 
5 19 3 4 4 3 3 1 5 4 5 1 2 5 1 3 4 2 4 5 5 3 1 2 2 1 1 5 3 2 4 1 3 5 5 2 1 4 2 3
5 3 5 4 1 3 1 5
2 1 2 1
6 2 6 4 5 1
3 6 1 2 2 1 3 1 2 3 3 2 1 3
6 5 2 4 1 2 4 3 2 6 5 6
4 9 1 3 2 4 3 4 4 3 1 4 4 2 2 3 3 2 4 1
3 3 3 1 3 2 1 2
5 2 4 1 4 2
3 5 1 2 2 1 3 1 3 2 1 3
2 0 
4 5 2 4 4 3 3 1 4 2 4 1
2 2 1 2 2 1
6 23 3 4 3 1 5 4 4 6 5 1 1 6 2 5 6 2 6 5 4 5 3 6 5 3 1 2 2 1 1 5 6 1 6 4 3 2 4 1 3 5 5 2 2 3 2 6
5 4 2 3 2 4 2 5 5 2
4 12 1 3 2 4 1 2 3 4 2 1 4 3 3 1 4 2 1 4 2 3 3 2 4 1
4 2 1 2 4 2
6 12 1 3 1 2 2 1 4 3 1 5 6 1 1 4 4 2 1 6 4 1 3 5 5 2
6 8 1 2 4 1 3 1 1 4 2 5 1 6 6 3 1 3
6 28 3 4 4 3 3 1 5 4 4 6 5 1 2 5 1 3 6 2 6 5 4 2 4 5 5 6 5 3 2 4 1 2 2 1 1 5 6 1 6 4 3 2 4 1 3 5 5 2 1 4 2 3 2 6 6 3
5 17 1 3 2 4 1 2 2 1 3 4 4 3 3 1 1 5 5 1 1 4 2 3 4 5 3 2 2 5 4 1 3 5 5 2
4 8 2 4 1 2 2 1 4 3 3 1 4 2 2 3 4 1
4 12 2 4 1 2 3 4 2 1 4 1 3 1 4 3 1 4 4 2 2 3 3 2 1 3
4 8 3 4 2 1 4 1 4 3 4 2 2 3 3 2 1 3
6 1 4 5
3 0 
6 0 
3 6 1 2 2 1 3 1 2 3 3 2 1 3
3 6 1 2 2 1 3 1 2 3 3 2 1 3
4 8 2 4 2 1 3 4 4 3 3 1 4 2 1 4 4 1
3 3 2 3 3 1 2 1
2 1 2 1
2 1 1 2
5 18 2 4 1 2 2 1 3 4 4 3 3 1 5 4 5 1 1 4 4 2 2 3 4 5 5 3 3 2 2 5 4 1 3 5 5 2
3 5 1 2 2 1 2 3 3 2 1 3
6 28 3 4 4 3 3 1 5 4 4 6 5 1 1 6 2 5 1 3 6 2 6 5 4 2 4 5 5 6 3 6 5 3 2 4 2 1 1 5 6 1 6 4 3 2 3 5 5 2 1 4 2 3 2 6 6 3
4 4 1 3 2 1 1 4 4 3
3 0 
6 10 6 2 2 1 6 5 1 5 6 1 5 4 6 4 5 6 1 6 3 2
2 1 2 1
2 0 
4 8 2 4 2 1 3 4 3 1 1 4 2 3 3 2 4 1
3 4 2 3 3 2 1 2 3 1
6 5 6 4 4 2 2 3 5 3 4 1
2 0 
2 1 2 1
2 1 1 2
6 2 6 2 4 6
3 0 
3 5 2 1 3 1 2 3 3 2 1 3
6 25 3 4 4 3 3 1 5 4 4 6 5 1 1 6 2 5 1 3 6 2 4 2 4 5 5 6 3 6 5 3 1 2 2 1 1 5 6 4 3 2 3 5 5 2 2 3 2 6 6 3
6 27 3 4 4 3 3 1 4 6 5 1 1 3 6 2 6 5 4 2 4 5 5 6 3 6 5 3 2 4 1 2 2 1 1 5 6 1 6 4 3 2 4 1 3 5 5 2 1 4 2 3 2 6 6 3
4 0 
3 0 
5 1 1 5
3 2 3 1 3 2
3 5 1 2 2 1 2 3 3 2 1 3
4 1 2 4
4 9 2 4 1 2 2 1 4 3 4 1 3 1 4 2 2 3 1 3
4 0 
3 0 
5 5 3 1 5 1 2 3 4 5 5 2
2 2 1 2 2 1
6 2 2 1 3 4
5 1 2 1
5 1 2 4
4 6 2 4 1 2 3 4 4 3 3 1 3 2
6 12 6 2 1 2 6 5 4 6 6 4 4 2 5 1 2 3 2 6 2 5 5 3 6 3
6 8 1 5 6 1 5 4 4 2 2 6 5 3 1 6 5 2
4 2 3 2 2 4
2 0 
4 0 
6 21 3 4 3 1 5 4 5 1 1 6 2 5 6 2 6 5 4 2 5 6 3 6 2 4 1 2 2 1 1 5 6 1 3 2 4 1 3 5 2 6 6 3
5 16 1 3 2 4 1 2 3 4 4 3 1 5 3 1 5 4 5 1 1 4 4 2 5 3 3 2 2 5 4 1 3 5
4 1 2 4
6 30 3 4 4 3 3 1 5 4 4 6 5 1 1 6 2 5 1 3 6 2 6 5 4 2 4 5 5 6 3 6 5 3 2 4 1 2 2 1 1 5 6 1 6 4 3 2 4 1 3 5 5 2 1 4 2 3 2 6 6 3
6 29 3 4 4 3 3 1 5 4 4 6 5 1 1 6 2 5 1 3 6 2 6 5 4 2 4 5 5 6 3 6 5 3 1 2 2 1 1 5 6 1 6 4 3 2 4 1 3 5 5 2 1 4 2 3 2 6 6 3
4 9 2 4 3 4 4 3 3 1 4 1 1 4 2 3 3 2 1 3
6 13 1 2 2 1 3 4 5 4 5 1 6 4 5 6 3 6 3 2 2 5 1 3 3 5 5 2
6 10 6 2 3 4 4 6 6 4 1 4 2 3 4 5 3 6 5 3 6 3
6 16 2 1 3 4 4 3 3 1 6 5 5 4 4 6 6 4 4 2 2 3 2 6 5 6 5 3 6 3 1 3 5 2
4 10 2 4 1 2 2 1 4 1 3 1 4 3 1 4 4 2 3 2 1 3
6 2 4 5 6 4
2 0 
4 5 4 1 4 2 2 3 3 2 1 3
3 2 3 1 3 2
6 19 3 4 3 1 5 4 4 6 5 1 1 6 2 5 1 3 4 5 5 6 3 6 5 3 2 4 2 1 1 5 6 4 3 2 3 5 2 3
4 12 1 3 2 4 1 2 2 1 3 4 4 3 3 1 4 2 1 4 2 3 3 2 4 1
2 0 
3 4 2 3 1 2 1 3 2 1
5 10 2 4 1 5 4 3 3 1 1 4 4 5 5 3 3 2 2 5 1 3
6 20 3 4 4 3 3 1 5 1 1 6 2 5 1 3 6 2 6 5 4 5 5 6 2 4 1 2 1 5 4 1 3 5 5 2 1 4 2 3 6 3
2 1 2 1
3 1 3 1
3 2 3 1 2 1
3 3 1 2 1 3 2 1
3 6 1 2 2 1 3 1 2 3 3 2 1 3
`


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
	scanner := bufio.NewScanner(strings.NewReader(testcasesDRaw))
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
