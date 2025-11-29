package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded reference solution (1364D.go).
const solutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m, k int
   G        [][]int
   st       []int
   dfn      []int
   col      []int
   col1     []int
   col2     []int
)

func dfs2(u, fa int) {
   st = append(st, u)
   dfn[u] = len(st)
   for _, v := range G[u] {
       if v == fa {
           continue
       }
       if dfn[v] == 0 {
           dfs2(v, u)
       } else {
           length := dfn[u] - dfn[v] + 1
           if length <= k {
               // found cycle
               fmt.Println(2)
               fmt.Println(length)
               for i := dfn[v] - 1; i < dfn[u]; i++ {
                   fmt.Printf("%d ", st[i])
               }
               fmt.Println()
               os.Exit(0)
           }
       }
   }
   st = st[:len(st)-1]
}

func dfs1(u, fa int) {
   if col[u] == 1 {
       col1 = append(col1, u)
   } else {
       col2 = append(col2, u)
   }
   need := (k + 1) / 2
   if len(col1) >= need {
       fmt.Println(1)
       for i := 0; i < need; i++ {
           fmt.Printf("%d ", col1[i])
       }
       fmt.Println()
       os.Exit(0)
   }
   if len(col2) >= need {
       fmt.Println(1)
       for i := 0; i < need; i++ {
           fmt.Printf("%d ", col2[i])
       }
       fmt.Println()
       os.Exit(0)
   }
   for _, v := range G[u] {
       if v == fa || col[v] != 0 {
           continue
       }
       col[v] = 3 - col[u]
       dfs1(v, u)
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   // fast input
   _, _ = fmt.Fscan(reader, &n, &m, &k)
   G = make([][]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       _, _ = fmt.Fscan(reader, &u, &v)
       G[u] = append(G[u], v)
       G[v] = append(G[v], u)
   }
   st = make([]int, 0, n)
   dfn = make([]int, n+1)
   // try find short cycle
   dfs2(1, 0)
   // no short cycle, output independent set
   fmt.Println(1)
   col = make([]int, n+1)
   col[1] = 1
   dfs1(1, 0)
}
`

const testcasesRaw = `5 8 5 1 2 3 4 1 5 2 3 4 5 2 5 1 3 3 5
2 1 1 1 2
1 0 1
1 0 1
2 1 1 1 2
2 1 1 1 2
1 0 1
5 0 2
5 1 4 2 5
1 0 1
3 2 2 2 3 1 3
4 6 2 2 4 1 2 3 4 1 4 2 3 1 3
4 3 1 2 4 1 2 3 4
3 2 1 2 3 1 3
1 0 1
3 1 1 1 3
1 0 1
3 1 2 1 3
1 0 1
3 3 2 2 3 1 2 1 3
4 3 4 2 3 2 4 1 2
1 0 1
1 0 1
1 0 1
3 3 2 2 3 1 2 1 3
3 3 3 2 3 1 2 1 3
3 1 1 2 3
5 10 2 2 4 1 2 3 4 1 5 1 4 2 3 4 5 2 5 1 3 3 5
1 0 1
5 3 1 4 5 2 5 1 3
5 8 2 2 4 1 2 3 4 1 5 2 3 4 5 2 5 1 3
3 1 1 1 3
4 1 2 1 4
5 5 2 2 4 1 2 3 4 2 3 3 5
4 0 2
4 0 1
5 2 4 2 4 3 5
2 0 1
3 3 2 2 3 1 2 1 3
5 9 2 2 4 1 2 3 4 1 5 1 4 2 3 4 5 2 5 1 3
1 0 1
1 0 1
2 1 2 1 2
3 1 3 1 3
4 2 2 2 3 2 4
1 0 1
4 6 4 2 4 1 2 3 4 1 4 2 3 1 3
2 0 1
2 0 2
3 0 2
4 4 2 2 3 1 3 3 4 1 4
2 1 1 1 2
5 9 4 2 4 1 2 3 4 1 5 1 4 2 3 4 5 1 3 3 5
4 5 4 2 4 1 2 1 4 2 3 1 3
1 0 1
4 6 2 2 4 1 2 3 4 1 4 2 3 1 3
2 1 1 1 2
3 1 2 1 3
4 0 1
1 0 1
4 1 3 2 4
4 5 3 2 4 3 4 1 4 2 3 1 3
4 4 3 1 2 1 3 3 4 1 4
2 0 1
4 1 1 3 4
2 0 2
1 0 1
1 0 1
4 3 3 2 4 1 3 3 4
4 4 1 2 3 2 4 3 4 1 4
4 3 3 2 3 1 2 1 4
5 5 2 2 4 1 2 3 4 4 5 2 5
4 2 1 2 4 1 2
1 0 1
3 2 2 1 2 1 3
1 0 1
4 5 3 2 4 1 2 1 4 2 3 1 3
4 4 2 2 3 1 2 1 3 3 4
1 0 1
3 1 3 1 2
4 5 2 1 2 3 4 1 4 2 3 1 3
4 4 3 2 3 2 4 1 3 1 4
3 0 3
3 0 2
4 0 3
2 0 2
3 0 1
4 1 1 3 4
1 0 1
1 0 1
2 1 1 1 2
5 10 4 2 4 1 2 3 4 1 5 1 4 2 3 4 5 2 5 1 3 3 5
5 5 2 2 4 1 2 3 4 4 5 3 5
2 0 2
5 4 4 1 2 2 5 1 3 1 4
4 0 1
2 1 1 1 2
3 2 1 2 3 1 2
5 5 2 2 4 1 2 3 4 1 4 2 3
1 0 1`

var _ = solutionSource

type testCase struct {
	line  string
	n     int
	m     int
	k     int
	edges [][2]int
}

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

func parseTests() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	tests := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tokens := strings.Fields(line)
		if len(tokens) < 3 {
			return nil, fmt.Errorf("invalid header on line %d", idx+1)
		}
		n, err := strconv.Atoi(tokens[0])
		if err != nil {
			return nil, fmt.Errorf("invalid n on line %d: %v", idx+1, err)
		}
		m, err := strconv.Atoi(tokens[1])
		if err != nil {
			return nil, fmt.Errorf("invalid m on line %d: %v", idx+1, err)
		}
		k, err := strconv.Atoi(tokens[2])
		if err != nil {
			return nil, fmt.Errorf("invalid k on line %d: %v", idx+1, err)
		}
		if len(tokens) != 3+2*m {
			return nil, fmt.Errorf("line %d: token count mismatch", idx+1)
		}
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			u, _ := strconv.Atoi(tokens[3+2*i])
			v, _ := strconv.Atoi(tokens[4+2*i])
			edges[i] = [2]int{u, v}
		}
		tests = append(tests, testCase{line: line, n: n, m: m, k: k, edges: edges})
	}
	return tests, nil
}

func normalize(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func solveExpected(tc testCase) string {
	G := make([][]int, tc.n+1)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		G[u] = append(G[u], v)
		G[v] = append(G[v], u)
	}
	st := make([]int, 0, tc.n)
	dfn := make([]int, tc.n+1)
	var res string

	var dfs2 func(u, fa int)
	dfs2 = func(u, fa int) {
		if res != "" {
			return
		}
		st = append(st, u)
		dfn[u] = len(st)
		for _, v := range G[u] {
			if v == fa {
				continue
			}
			if dfn[v] == 0 {
				dfs2(v, u)
				if res != "" {
					return
				}
			} else {
				length := dfn[u] - dfn[v] + 1
				if length <= tc.k {
					var sb strings.Builder
					sb.WriteString("2\n")
					sb.WriteString(fmt.Sprintf("%d\n", length))
					for i := dfn[v] - 1; i < dfn[u]; i++ {
						sb.WriteString(fmt.Sprintf("%d ", st[i]))
					}
					sb.WriteString("\n")
					res = sb.String()
					return
				}
			}
		}
		st = st[:len(st)-1]
	}

	dfs2(1, 0)
	if res != "" {
		return normalize(res)
	}

	col := make([]int, tc.n+1)
	col1 := make([]int, 0)
	col2 := make([]int, 0)
	var out strings.Builder
	out.WriteString("1\n")

	var dfs1 func(u, fa int)
	dfs1 = func(u, fa int) {
		if res != "" {
			return
		}
		if col[u] == 1 {
			col1 = append(col1, u)
		} else {
			col2 = append(col2, u)
		}
		need := (tc.k + 1) / 2
		if len(col1) >= need {
			out.WriteString("1\n")
			for i := 0; i < need; i++ {
				out.WriteString(fmt.Sprintf("%d ", col1[i]))
			}
			out.WriteString("\n")
			res = out.String()
			return
		}
		if len(col2) >= need {
			out.WriteString("1\n")
			for i := 0; i < need; i++ {
				out.WriteString(fmt.Sprintf("%d ", col2[i]))
			}
			out.WriteString("\n")
			res = out.String()
			return
		}
		for _, v := range G[u] {
			if v == fa || col[v] != 0 {
				continue
			}
			col[v] = 3 - col[u]
			dfs1(v, u)
		}
	}

	col[1] = 1
	dfs1(1, 0)
	if res != "" {
		return normalize(res)
	}
	return normalize(out.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for idx, tc := range tests {
		input := tc.line + "\n"
		out, err := run(bin, input)
		if err != nil {
			fail(idx+1, tc.line, fmt.Sprintf("candidate runtime error: %v", err))
		}
		got := normalize(out)
		want := normalize(solveExpected(tc))
		if got != want {
			fail(idx+1, tc.line, fmt.Sprintf("mismatch:\nexpected: %s\ngot: %s", want, got))
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
