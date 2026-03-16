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

const testcasesBRaw = `2 5 1 3 1
7 4 2 1 4 4 5 3 2 5 1 3 1 5 2
7 4 5 2 3 4 3 1 4 5 1 1 3 5 2
4 3 5 4 5 1 4 2 4
6 2 3 5 3 1 4 5 1 2 5 4 3
7 1 4 1 3 5 4 5 2 1 2 2 4 5 3
9 3 4 3 5 5 1 4 5 2 5 4 1 3 1 2 1 3 2
4 1 5 2 3 3 1 3 5
2 4 3 2 3
1 3 5
3 5 4 1 2 1 4
2 1 2 4 5
10 4 5 2 5 4 2 5 1 3 4 1 3 5 3 2 1 1 4 2 3
9 2 3 4 2 3 1 4 5 3 5 2 1 2 5 5 1 4 1
6 3 1 4 1 2 5 1 5 1 2 4 2
1 4 2
10 1 4 4 5 3 5 3 4 3 1 2 3 5 2 1 5 1 2 4 2
3 3 4 4 2 3 2
0
1 5 3
2 5 2 3 5
5 2 4 5 1 2 3 3 4 5 2
8 1 5 1 3 1 2 4 2 3 4 2 5 1 4 4 5
0
0
10 5 2 3 2 3 5 3 4 4 1 4 2 3 1 5 1 2 1 4 5
10 3 5 3 1 1 4 4 3 5 4 4 2 5 1 5 2 3 2 2 1
6 3 5 3 2 5 2 3 1 1 5 3 4
8 5 1 2 3 5 3 2 4 5 4 2 1 3 4 2 5
7 5 2 5 3 4 5 2 4 1 2 3 1 4 3
7 1 2 5 1 4 5 2 4 2 3 4 3 1 4
2 5 1 5 4
4 2 1 4 3 5 1 5 2
0
4 1 4 1 2 2 3 2 4
6 3 2 2 1 4 5 1 3 5 2 5 1
4 2 1 2 4 1 5 4 5
5 1 3 1 2 5 1 4 2 1 4
0
4 3 2 3 4 1 3 1 4
3 2 3 3 5 4 5
7 1 2 4 5 1 3 3 4 5 3 1 4 3 2
9 3 5 3 4 5 1 2 3 1 2 4 5 5 2 3 1 4 2
6 4 3 5 2 4 1 2 1 5 4 3 1
8 1 3 5 3 5 4 4 3 2 5 2 3 1 4 5 1
5 3 1 1 4 3 4 3 2 4 2
0
2 2 5 3 4
4 5 3 4 3 4 2 4 1
1 3 2
7 1 4 2 1 5 1 5 2 5 4 3 1 2 4
0
3 1 4 3 5 1 5
7 1 2 4 5 5 3 3 4 5 2 4 1 1 3
4 1 5 4 3 1 2 5 3
4 2 4 3 2 4 5 5 1
8 2 4 4 3 3 5 3 2 5 2 5 4 1 5 3 1
2 2 3 2 1
10 5 1 5 2 1 2 3 4 3 1 3 5 5 4 3 2 4 1 2 4
1 3 1
10 5 3 1 5 3 4 4 1 5 2 2 4 4 5 2 1 1 3 2 3
5 2 4 4 3 4 5 1 5 1 4
6 1 2 5 3 5 1 3 1 2 5 4 1
8 2 4 3 2 1 5 5 2 1 2 1 4 1 3 4 3
6 2 5 3 5 4 5 2 4 4 3 3 2
4 5 3 2 5 1 3 3 2
4 3 4 3 5 4 1 3 2
3 1 5 2 5 4 2
9 2 5 2 1 4 2 5 1 2 3 4 5 3 1 5 3 3 4
4 5 1 1 3 4 5 4 1
3 4 5 2 5 3 1
1 2 5
6 4 5 5 2 3 1 4 3 1 2 1 4
0
7 2 4 5 3 2 1 2 5 4 5 3 2 5 1
5 1 4 1 5 3 4 5 4 3 5
5 2 5 3 1 5 3 5 1 4 1
8 4 3 1 2 5 2 5 1 5 4 3 5 3 1 2 3
6 1 2 3 5 3 2 3 4 1 5 5 2
0
7 3 1 4 2 5 1 2 5 3 4 2 1 3 5
5 2 1 1 4 3 5 5 2 1 5
3 2 4 3 2 4 5
5 3 4 1 2 5 2 3 5 4 5
7 3 1 4 1 1 5 5 4 5 3 4 2 2 3
10 5 3 5 4 3 1 5 1 3 2 2 4 4 1 1 2 3 4 2 5
0
9 2 4 5 2 1 3 4 5 3 5 4 3 5 1 4 1 2 1
1 3 1
1 2 1
10 1 4 1 5 2 4 2 3 3 5 2 5 1 2 1 3 4 5 4 3
8 3 4 5 1 3 2 3 5 2 5 4 2 1 4 1 2
10 3 1 2 4 5 1 1 2 5 3 5 2 5 4 2 3 4 1 4 3
2 1 3 1 5
3 2 1 4 2 5 3
2 5 1 2 3
7 5 1 1 3 4 5 3 2 4 1 5 3 4 3
6 1 3 1 2 3 2 5 2 5 4 5 3
9 3 2 3 1 4 3 2 4 1 2 1 5 5 4 3 5 2 5
7 4 1 4 5 2 5 5 3 1 3 3 4 2 4
`

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

func expected(m int, edges [][2]int) string {
	var g [6][6]bool
	for _, e := range edges {
		g[e[0]][e[1]] = true
		g[e[1]][e[0]] = true
	}
	for i := 1; i <= 5; i++ {
		for j := i + 1; j <= 5; j++ {
			for k := j + 1; k <= 5; k++ {
				cnt := 0
				if g[i][j] {
					cnt++
				}
				if g[i][k] {
					cnt++
				}
				if g[j][k] {
					cnt++
				}
				if cnt == 3 || cnt == 0 {
					return "WIN"
				}
			}
		}
	}
	return "FAIL"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 1 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		m, err := strconv.Atoi(parts[0])
		if err != nil || len(parts) != 1+2*m {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			a, _ := strconv.Atoi(parts[1+2*i])
			b, _ := strconv.Atoi(parts[2+2*i])
			edges[i] = [2]int{a, b}
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", m))
		for _, e := range edges {
			input.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		want := expected(m, edges)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
