package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded reference solution (129B.go).
const solutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   adj := make([][]int, n)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       a--
       b--
       adj[a] = append(adj[a], b)
       adj[b] = append(adj[b], a)
   }
   deg := make([]int, n)
   removed := make([]bool, n)
   for i := 0; i < n; i++ {
       deg[i] = len(adj[i])
   }
   rounds := 0
   for {
       var toRemove []int
       for i := 0; i < n; i++ {
           if !removed[i] && deg[i] == 1 {
               toRemove = append(toRemove, i)
           }
       }
       if len(toRemove) == 0 {
           break
       }
       rounds++
       // mark removed
       for _, u := range toRemove {
           removed[u] = true
       }
       // update degrees
       for _, u := range toRemove {
           for _, v := range adj[u] {
               if !removed[v] {
                   deg[v]--
               }
           }
       }
   }
   fmt.Fprintln(writer, rounds)
}
`

const testcasesRaw = `1 0
7 12 2 7 3 7 4 6 1 4 5 7 2 3 4 5 2 6 3 6 1 6 2 5 3 5
8 23 3 4 3 7 4 6 5 7 1 6 2 5 1 3 2 8 6 8 4 5 5 6 4 8 3 6 2 4 1 5 6 7 4 7 3 5 3 8 5 8 1 4 1 7 2 6
6 4 2 3 3 4 3 5 1 4
1 0
2 1 1 2
4 2 1 3 1 4
10 3 2 3 2 9 6 9
10 9 4 10 3 10 2 3 8 9 5 6 3 6 1 6 6 9 4 7
2 0
8 16 2 4 3 4 5 8 1 5 4 6 5 7 1 4 2 3 6 7 1 7 4 5 2 6 5 6 1 6 1 3 2 8
8 0
3 0
3 0
2 1 1 2
5 4 4 5 1 2 3 5 1 5
2 0
10 35 3 4 4 9 3 7 4 6 5 7 5 10 1 6 2 5 1 9 2 8 7 10 4 5 3 9 5 6 3 6 5 9 9 10 1 2 2 7 1 5 2 10 1 8 7 9 6 7 3 5 4 10 5 8 8 10 2 3 2 9 1 7 2 6 1 10 6 9 7 8
1 0
5 6 1 2 3 4 4 5 2 5 1 3 3 5
1 0
6 1 2 5
2 0
1 0
3 1 2 3
9 31 3 4 4 9 3 7 5 7 8 9 1 6 2 5 1 3 2 8 6 8 4 5 3 9 5 6 4 8 3 6 2 4 1 2 2 7 1 5 6 7 4 7 3 5 3 8 5 8 1 4 2 3 2 9 1 7 2 6 6 9 7 8
3 0
1 0
7 11 2 4 1 2 3 7 5 7 4 5 1 7 5 6 3 6 1 6 1 3 3 5
2 0
5 3 2 5 3 4 1 5
2 1 1 2
10 13 2 4 1 2 3 4 4 9 2 10 1 8 5 7 5 10 8 9 2 6 1 10 3 6 7 8
7 7 2 4 4 6 1 4 2 3 1 7 3 6 2 5
5 7 2 4 1 2 3 4 4 5 2 5 1 3 3 5
8 27 3 4 3 7 4 6 5 7 1 6 2 5 1 3 2 8 4 5 5 6 4 8 3 6 2 4 1 2 2 7 1 5 1 8 6 7 4 7 3 5 3 8 5 8 1 4 2 3 1 7 2 6 7 8
10 12 3 4 5 8 8 10 7 9 1 7 8 9 2 6 4 5 5 9 1 6 1 3 7 8
8 5 2 4 6 7 4 5 1 6 2 5
10 37 4 9 3 7 3 10 5 7 5 10 8 9 1 6 1 3 1 9 2 8 6 8 4 5 3 9 5 6 4 8 3 6 5 9 9 10 2 7 1 5 2 10 1 8 7 9 6 10 4 7 3 5 4 10 3 8 8 10 1 4 2 3 2 9 1 7 2 6 1 10 6 9 7 8
10 36 3 4 3 10 5 7 5 10 8 9 2 5 1 3 1 9 2 8 6 8 4 5 3 9 5 6 4 8 3 6 5 9 9 10 2 4 1 2 2 7 1 5 2 10 1 8 7 9 6 10 4 7 3 5 4 10 3 8 5 8 8 10 1 4 1 7 2 6 1 10 7 8
1 0
8 10 2 7 1 8 5 7 1 4 2 3 4 5 2 6 2 5 1 3 3 5
7 17 1 2 3 4 2 7 1 5 3 7 1 4 5 7 2 3 6 7 1 7 5 6 3 6 1 6 2 5 1 3 4 7 3 5
1 0
4 3 2 3 2 4 1 4
10 20 3 4 3 7 1 6 1 3 2 8 6 8 4 5 5 6 9 10 2 4 2 7 2 10 7 9 6 7 6 10 4 7 3 8 2 6 1 10 6 9
9 14 2 4 1 2 3 8 2 7 3 7 4 6 1 8 1 4 5 7 3 9 1 6 6 9 1 9 2 8
1 0
10 13 9 10 3 8 5 8 3 10 1 4 2 3 4 7 6 7 8 9 6 10 1 10 1 9 7 8
8 7 3 7 4 6 1 4 2 6 2 5 4 7 3 5
7 4 6 7 4 5 1 3 4 7
7 11 2 4 1 2 3 4 3 7 4 6 4 5 2 6 3 6 2 5 1 3 4 7
7 20 3 4 3 7 4 6 5 7 1 6 2 5 1 3 4 5 5 6 3 6 2 4 2 7 1 5 6 7 4 7 3 5 1 4 2 3 1 7 2 6
9 3 1 6 3 8 5 6
9 1 5 7
8 8 3 4 3 7 4 6 2 3 4 5 5 6 4 8 4 7
8 22 4 6 5 7 1 6 2 5 1 3 2 8 6 8 4 5 4 8 2 4 1 2 1 5 6 7 4 7 3 5 3 8 5 8 1 4 2 3 1 7 2 6 7 8
6 0
6 10 2 4 1 2 3 4 1 5 1 4 4 5 2 6 5 6 1 6 2 5
3 3 2 3 1 2 1 3
5 9 2 4 1 2 3 4 1 5 1 4 2 3 4 5 2 5 1 3
7 10 2 4 1 2 2 7 1 5 5 7 6 7 4 5 1 7 1 6 1 3
6 8 3 4 1 5 1 4 2 3 2 6 5 6 3 6 1 6
3 3 2 3 1 2 1 3
3 0
4 5 2 4 1 2 3 4 2 3 1 3
4 3 1 2 3 4 1 4
7 20 3 4 3 7 4 6 5 7 1 6 2 5 1 3 4 5 5 6 3 6 2 4 1 2 1 5 6 7 4 7 3 5 1 4 2 3 1 7 2 6
7 1 1 7
5 9 2 4 1 2 3 4 1 5 2 3 4 5 2 5 1 3 3 5
10 16 4 10 2 4 3 4 5 8 1 5 2 10 6 8 8 10 7 9 4 5 8 9 2 6 5 6 6 9 4 7 7 8
10 30 4 9 3 7 4 6 3 10 5 7 5 10 8 9 2 5 1 9 2 8 7 10 6 8 4 5 3 9 5 6 3 6 5 9 2 4 1 2 2 7 1 5 1 8 6 7 4 10 3 8 5 8 8 10 1 4 2 6 6 9
8 9 3 8 3 4 1 5 4 6 1 8 1 4 4 5 5 6 7 8
7 19 3 7 4 6 5 7 1 6 2 5 1 3 4 5 5 6 3 6 2 4 1 2 2 7 1 5 6 7 4 7 3 5 1 4 1 7 2 6
5 3 1 3 3 5 1 5
2 0
4 3 2 4 1 2 1 3
6 15 2 4 1 2 3 4 1 5 4 6 1 4 2 3 4 5 2 6 5 6 3 6 1 6 2 5 1 3 3 5
5 9 2 4 1 2 1 5 1 4 2 3 4 5 2 5 1 3 3 5
7 16 2 4 1 2 4 6 1 4 2 3 6 7 1 7 4 5 2 6 5 6 3 6 1 6 2 5 1 3 4 7 3 5
1 0
3 1 2 3
1 0
9 33 3 4 4 9 3 7 4 6 5 7 8 9 1 6 1 3 1 9 2 8 6 8 4 5 3 9 5 6 4 8 5 9 2 4 1 2 2 7 1 5 1 8 7 9 6 7 4 7 3 5 3 8 5 8 1 4 2 3 2 9 1 7 2 6 6 9
7 21 3 4 3 7 4 6 5 7 1 6 2 5 1 3 4 5 5 6 3 6 2 4 1 2 2 7 1 5 6 7 4 7 3 5 1 4 2 3 1 7 2 6
5 7 2 4 1 2 3 4 1 5 1 4 2 3 1 3
10 45 3 4 4 9 3 7 4 6 3 10 5 7 5 10 8 9 1 6 2 5 1 3 1 9 2 8 7 10 6 8 4 5 3 9 5 6 4 8 3 6 5 9 9 10 2 4 1 2 2 7 1 5 2 10 1 8 7 9 6 7 6 10 4 7 3 5 4 10 3 8 5 8 8 10 1 4 2 3 2 9 1 7 2 6 1 10 6 9 7 8
3 3 2 3 1 2 1 3
4 3 1 3 3 4 1 4
1 0
3 1 2 3
1 0
4 4 2 3 2 4 1 2 3 4
4 3 2 4 1 2 1 4
4 2 2 4 1 2
1 0
2 0
8 19 3 4 4 6 5 7 1 3 6 8 4 5 5 6 4 8 3 6 2 4 1 5 1 8 6 7 4 7 3 5 5 8 1 4 2 3 2 6
1 0
5 5 1 2 3 4 2 3 4 5 2 5`

var _ = solutionSource

type testCase struct {
	n     int
	m     int
	edges [][2]int
}

func runProgram(bin, input string) (string, error) {
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

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func parseTests() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	tests := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("invalid test line %d", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("bad n on line %d: %v", idx+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("bad m on line %d: %v", idx+1, err)
		}
		if len(fields) != 2+2*m {
			return nil, fmt.Errorf("line %d edge count mismatch", idx+1)
		}
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			a := atoi(fields[2+2*i])
			b := atoi(fields[2+2*i+1])
			edges[i] = [2]int{a, b}
		}
		tests = append(tests, testCase{n: n, m: m, edges: edges})
	}
	return tests, nil
}

func expected(tc testCase) int {
	adj := make([][]int, tc.n)
	for _, e := range tc.edges {
		a := e[0] - 1
		b := e[1] - 1
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	deg := make([]int, tc.n)
	removed := make([]bool, tc.n)
	for i := 0; i < tc.n; i++ {
		deg[i] = len(adj[i])
	}
	rounds := 0
	for {
		var toRemove []int
		for i := 0; i < tc.n; i++ {
			if !removed[i] && deg[i] == 1 {
				toRemove = append(toRemove, i)
			}
		}
		if len(toRemove) == 0 {
			break
		}
		rounds++
		for _, u := range toRemove {
			removed[u] = true
		}
		for _, u := range toRemove {
			for _, v := range adj[u] {
				if !removed[v] {
					deg[v]--
				}
			}
		}
	}
	return rounds
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTests()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range tests {
		var input strings.Builder
		input.WriteString(strconv.Itoa(tc.n))
		input.WriteByte(' ')
		input.WriteString(strconv.Itoa(tc.m))
		input.WriteByte('\n')
		for _, e := range tc.edges {
			input.WriteString(strconv.Itoa(e[0]))
			input.WriteByte(' ')
			input.WriteString(strconv.Itoa(e[1]))
			input.WriteByte('\n')
		}
		want := strconv.Itoa(expected(tc))
		got, err := runProgram(bin, input.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
