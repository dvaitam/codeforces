package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded reference solution (1388D.go).
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

   var n int
   fmt.Fscan(reader, &n)

   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   b := make([]int, n)
   for i := 0; i < n; i++ {
       var bi int
       fmt.Fscan(reader, &bi)
       if bi == -1 {
           b[i] = -1
       } else {
           b[i] = bi - 1
       }
   }

   // f values
   f := make([]int64, n)
   copy(f, a)
   // indegree for initial graph
   indeg := make([]int, n)
   for i := 0; i < n; i++ {
       if b[i] >= 0 {
           indeg[b[i]]++
       }
   }
   // build edges
   graph := make([][]int, n)
   // indegree for DAG
   deg := make([]int, n)

   // initial queue
   q := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if indeg[i] == 0 {
           q = append(q, i)
       }
   }
   // process in top order
   for head := 0; head < len(q); head++ {
       x := q[head]
       if b[x] >= 0 {
           p := b[x]
           if f[x] < 0 {
               // process x after p: edge p->x
               graph[p] = append(graph[p], x)
               deg[x]++
           } else {
               // process x before p: edge x->p, and add f[x]
               f[p] += f[x]
               graph[x] = append(graph[x], p)
               deg[p]++
           }
           indeg[p]--
           if indeg[p] == 0 {
               q = append(q, p)
           }
       }
   }
   // compute answer
   var ans int64
   for i := 0; i < n; i++ {
       ans += f[i]
   }
   fmt.Fprintln(writer, ans)

   // final order topological sort
   q2 := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if deg[i] == 0 {
           q2 = append(q2, i)
       }
   }
   order := make([]int, 0, n)
   for head := 0; head < len(q2); head++ {
       x := q2[head]
       order = append(order, x)
       for _, y := range graph[x] {
           deg[y]--
           if deg[y] == 0 {
               q2 = append(q2, y)
           }
       }
   }
   // print order 1-based
   for i, v := range order {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v+1)
   }
   writer.WriteByte('\n')
}
`

const testcasesRaw = `1 4 -1
5 4 0 2 2 -1 -1 1 4 5 1
4 2 3 0 -4 2 1 4 3
2 0 -1 1 -1
2 -5 0 2 2
2 0 2 2 2
2 -1 0 2 -1
3 -5 0 4 -1 1 2
1 2 -1
1 -2 -1
1 -2 1
1 -5 1
4 -3 -2 2 1 1 2 2 1
3 1 1 -5 3 2 3
1 4 -1
4 1 -3 -5 3 1 4 4 1
1 0 -1
2 -2 -5 -1 2
5 -3 -4 1 4 -4 4 5 3 5 1
5 4 -5 -1 0 1 -1 5 -1 3 -1
3 -1 5 -3 3 1 2
2 1 0 1 1
4 -5 -1 3 -1 4 3 -1 4
5 3 -1 5 -5 2 3 5 -1 3 2
4 -5 -5 -1 -5 2 4 2 1
5 3 0 1 -1 -2 -1 4 2 1 4
5 5 0 -3 -3 0 5 -1 4 -1 4
2 0 0 1 2
3 0 2 1 3 1 -1
2 4 -5 1 -1
3 -5 2 5 2 1 -1
5 1 -1 -3 3 -3 -1 5 5 1 4
1 3 1
4 -1 -1 -1 -5 3 2 2 4
5 3 0 0 -2 1 1 -1 4 1 5
5 1 0 2 -5 3 3 5 4 1 -1
3 3 -3 5 1 2 3
1 -2 -1
3 -3 1 -4 3 1 3
5 -4 -2 -3 2 -4 3 5 3 2 2
4 0 4 0 -4 2 -1 3 -1
3 -2 1 1 3 3 -1
5 2 0 4 -3 4 5 2 2 -1 3
4 3 -3 -5 -4 4 2 2 -1
1 -2 -1
5 2 -5 0 -5 0 3 1 5 2 3
2 4 -3 1 1
5 -5 -3 -3 -3 2 5 5 5 5 -1
5 -5 3 5 1 -3 2 4 -1 -1 4
2 -1 -2 1 1
3 -1 3 2 1 3 1
1 5 -1
5 -5 0 -4 -2 5 3 4 1 3 4
2 0 5 -1 1
5 -5 3 -4 3 0 -1 -1 -1 3 4
3 4 2 0 3 2 -1
2 0 -5 -1 2
2 -5 0 2 -1
1 1 -1
3 1 -5 4 1 2 -1
4 -5 2 0 4 4 2 2 4
4 1 -3 -5 2 2 1 3 -1
3 -4 -4 -5 2 -1 1
4 4 5 -5 0 3 4 3 3
1 -5 1
3 -5 5 3 -1 -1 2
3 -4 2 -5 1 2 -1
1 1 1
2 3 -3 -1 -1
2 5 5 -1 2
3 -5 2 5 3 -1 1
2 5 -1 1 -1
2 0 -2 -1 1
4 0 4 -3 1 4 -1 1 4
1 5 1
2 -3 -5 -1 1
5 -5 -5 -5 -4 4 4 1 1 5 3
1 1 1
5 5 0 -2 -3 0 4 1 4 3 -1
2 1 4 2 1
1 1 1
2 2 1 -1 1
2 5 2 1 2
1 -1 1
5 0 3 -5 4 4 3 1 2 -1 4
3 1 5 -4 -1 1 3
1 1 1
4 -4 2 4 2 1 1 2 3
4 -3 4 -1 3 -1 2 2 1
3 2 5 3 -1 1 2
2 4 0 1 1
1 -1 1
1 5 1
3 1 -2 0 1 1 -1
5 0 -3 -5 1 4 3 3 -1 5 5
1 1 -1
2 -3 -2 -1 -1
1 3 -1
4 0 4 -1 0 -1 3 2 1
3 5 3 0 1 3 2
1 1 -1`

var _ = solutionSource

type testCase struct {
	n int
	a []int64
	b []int
}

func solve(n int, a []int64, b []int) (int64, []int) {
	f := make([]int64, n)
	copy(f, a)
	indeg := make([]int, n)
	for i := 0; i < n; i++ {
		if b[i] >= 0 {
			indeg[b[i]]++
		}
	}
	graph := make([][]int, n)
	deg := make([]int, n)
	q := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if indeg[i] == 0 {
			q = append(q, i)
		}
	}
	for head := 0; head < len(q); head++ {
		x := q[head]
		if b[x] >= 0 {
			p := b[x]
			if f[x] < 0 {
				graph[p] = append(graph[p], x)
				deg[x]++
			} else {
				f[p] += f[x]
				graph[x] = append(graph[x], p)
				deg[p]++
			}
			indeg[p]--
			if indeg[p] == 0 {
				q = append(q, p)
			}
		}
	}
	var ans int64
	for i := 0; i < n; i++ {
		ans += f[i]
	}
	q2 := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if deg[i] == 0 {
			q2 = append(q2, i)
		}
	}
	order := make([]int, 0, n)
	for head := 0; head < len(q2); head++ {
		x := q2[head]
		order = append(order, x)
		for _, y := range graph[x] {
			deg[y]--
			if deg[y] == 0 {
				q2 = append(q2, y)
			}
		}
	}
	for i := range order {
		order[i]++
	}
	return ans, order
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
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("bad n on line %d: %v", idx+1, err)
		}
		if len(fields) != 1+2*n {
			return nil, fmt.Errorf("line %d: expected %d tokens got %d", idx+1, 1+2*n, len(fields))
		}
		a := make([]int64, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(fields[1+i], 10, 64)
			a[i] = v
		}
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[1+n+i])
			if v == -1 {
				b[i] = -1
			} else {
				b[i] = v - 1
			}
		}
		tests = append(tests, testCase{n: n, a: a, b: b})
	}
	return tests, nil
}

func runCase(bin string, n int, a []int64, b []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(a[i], 10))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		if b[i] < 0 {
			sb.WriteString("-1")
		} else {
			sb.WriteString(strconv.Itoa(b[i] + 1))
		}
	}
	sb.WriteByte('\n')
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
	outFields := strings.Fields(strings.TrimSpace(out.String()))
	if len(outFields) < 1+n {
		return fmt.Errorf("not enough output")
	}
	ansGot, err := strconv.ParseInt(outFields[0], 10, 64)
	if err != nil {
		return fmt.Errorf("bad sum output")
	}
	orderGot := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(outFields[1+i])
		if err != nil {
			return fmt.Errorf("bad order value")
		}
		orderGot[i] = v
	}
	ansExp, orderExp := solve(n, a, b)
	if ansGot != ansExp {
		return fmt.Errorf("expected sum %d got %d", ansExp, ansGot)
	}
	for i := 0; i < n; i++ {
		if orderGot[i] != orderExp[i] {
			return fmt.Errorf("expected order %v got %v", orderExp, orderGot)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTests()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for idx, tc := range tests {
		if err := runCase(bin, tc.n, tc.a, tc.b); err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
