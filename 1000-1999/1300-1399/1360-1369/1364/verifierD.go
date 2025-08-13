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

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func fail(idx int, line, msg string) {
	fmt.Printf("case %d failed\ninput: %s\n%s\n", idx, line, msg)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Println("could not open testcasesD.txt:", err)
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
		input := line + "\n"

		tokens := strings.Fields(line)
		if len(tokens) < 3 {
			fail(idx, line, "invalid testcase header")
		}
		n, _ := strconv.Atoi(tokens[0])
		m, _ := strconv.Atoi(tokens[1])
		k, _ := strconv.Atoi(tokens[2])
		if len(tokens) != 3+2*m {
			fail(idx, line, "invalid number of tokens in testcase")
		}

		adj := make([]map[int]bool, n+1)
		for i := 0; i < m; i++ {
			u, _ := strconv.Atoi(tokens[3+2*i])
			v, _ := strconv.Atoi(tokens[4+2*i])
			if adj[u] == nil {
				adj[u] = map[int]bool{}
			}
			if adj[v] == nil {
				adj[v] = map[int]bool{}
			}
			adj[u][v] = true
			adj[v][u] = true
		}

		out, err := run(bin, input)
		if err != nil {
			fail(idx, line, fmt.Sprintf("candidate runtime error: %v", err))
		}
		fields := strings.Fields(out)
		if len(fields) == 0 {
			fail(idx, line, "empty output")
		}

		typ := fields[0]
		if typ == "1" {
			need := (k + 1) / 2
			if len(fields) != 1+need {
				fail(idx, line, fmt.Sprintf("expected %d vertices, got %d", need, len(fields)-1))
			}
			seen := map[int]bool{}
			verts := make([]int, need)
			for i := 0; i < need; i++ {
				v, err := strconv.Atoi(fields[1+i])
				if err != nil || v < 1 || v > n {
					fail(idx, line, "invalid vertex in independent set")
				}
				if seen[v] {
					fail(idx, line, "duplicate vertex in independent set")
				}
				seen[v] = true
				verts[i] = v
			}
			for i := 0; i < need; i++ {
				for j := i + 1; j < need; j++ {
					if adj[verts[i]][verts[j]] {
						fail(idx, line, "set is not independent")
					}
				}
			}
		} else if typ == "2" {
			if len(fields) < 2 {
				fail(idx, line, "missing cycle length")
			}
			c, err := strconv.Atoi(fields[1])
			if err != nil {
				fail(idx, line, "invalid cycle length")
			}
			if c > k {
				fail(idx, line, "cycle length exceeds k")
			}
			if len(fields) != 2+c {
				fail(idx, line, "wrong number of vertices in cycle")
			}
			cycle := make([]int, c)
			used := map[int]bool{}
			for i := 0; i < c; i++ {
				v, err := strconv.Atoi(fields[2+i])
				if err != nil || v < 1 || v > n {
					fail(idx, line, "invalid vertex in cycle")
				}
				if used[v] {
					fail(idx, line, "cycle contains duplicate vertex")
				}
				used[v] = true
				cycle[i] = v
			}
			for i := 0; i < c; i++ {
				u := cycle[i]
				v := cycle[(i+1)%c]
				if adj[u] == nil || !adj[u][v] {
					fail(idx, line, "cycle uses non-existent edge")
				}
			}
		} else {
			fail(idx, line, "first token must be 1 or 2")
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
