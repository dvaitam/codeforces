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

// Embedded source for the reference solution (was 1186F.go).
const solutionSource = `package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// Item is an entry in the priority queue
type Item struct {
   id, d int
}

// PriorityQueue implements a min-heap of Items by d, then id
type PriorityQueue []Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
   if pq[i].d != pq[j].d {
       return pq[i].d < pq[j].d
   }
   return pq[i].id < pq[j].id
}
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
   *pq = append(*pq, x.(Item))
}
func (pq *PriorityQueue) Pop() interface{} {
   old := *pq
   n := len(old)
   it := old[n-1]
   *pq = old[0 : n-1]
   return it
}

// popValid pops until it finds a valid Item matching current d and >0, or returns false
func popValid(pq *PriorityQueue, d []int) (Item, bool) {
   for pq.Len() > 0 {
       it := heap.Pop(pq).(Item)
       if it.d != d[it.id] {
           continue
       }
       if it.d <= 0 {
           continue
       }
       return it, true
   }
   return Item{}, false
}

// EdgeKey for visited edges
type EdgeKey struct{ u, v int }

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)
   d := make([]int, n+1)
   adj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       d[x]++
       d[y]++
       adj[x] = append(adj[x], y)
       adj[y] = append(adj[y], x)
   }
   // initial ceil(deg/2)
   for i := 1; i <= n; i++ {
       d[i] = d[i]/2 + d[i]%2
   }
   // build initial heap
   pq := &PriorityQueue{}
   heap.Init(pq)
   for i := 1; i <= n; i++ {
       heap.Push(pq, Item{id: i, d: d[i]})
   }

   visited := make(map[EdgeKey]struct{})
   type Nd struct{ d, id int }
   var temp []Nd
   var res [][2]int

   for i := 1; i <= n; i++ {
       it, ok := popValid(pq, d)
       if !ok {
           break
       }
       x := it.id
       temp = temp[:0]
       for _, y := range adj[x] {
           // skip if reversed edge visited
           if _, seen := visited[EdgeKey{u: y, v: x}]; seen {
               continue
           }
           temp = append(temp, Nd{d: d[y], id: y})
       }
       sort.Slice(temp, func(i, j int) bool {
           if temp[i].d != temp[j].d {
               return temp[i].d > temp[j].d
           }
           return temp[i].id > temp[j].id
       })
       for _, nd := range temp {
           if d[x] <= 0 {
               break
           }
           y := nd.id
           res = append(res, [2]int{x, y})
           visited[EdgeKey{u: x, v: y}] = struct{}{}
           d[x]--
           d[y]--
           if d[y] >= 0 {
               heap.Push(pq, Item{id: y, d: d[y]})
           }
       }
   }
   // output
   fmt.Fprintln(out, len(res))
   for _, e := range res {
       fmt.Fprintln(out, e[0], e[1])
   }
}
`

const testcasesRaw = `4 6 2 4 1 2 3 4 1 4 2 3 1 3
5 4 4 5 1 2 1 4 1 5
3 1 1 3
5 6 1 2 3 4 1 4 2 3 1 3 3 5
1 0
2 0
1 0
1 0
6 3 1 3 1 4 1 5
2 0
6 1 2 4
6 0
6 0
5 6 2 4 1 2 1 5 1 4 4 5 1 3
2 0
1 0
2 1 1 2
5 7 1 2 3 4 1 5 1 4 2 3 2 5 1 3
2 1 1 2
6 4 4 6 2 5 2 6 3 6
4 3 2 3 1 2 1 4
4 4 2 4 1 2 3 4 1 4
1 0
1 0
6 8 1 2 1 5 1 4 2 3 5 6 3 6 1 6 1 3
3 1 1 2
2 1 1 2
6 7 2 4 1 2 3 4 4 6 2 3 2 5 3 5
3 2 2 3 1 2
5 0
1 0
2 0
5 4 2 4 2 5 1 3 3 4
1 0
2 0
3 3 2 3 1 2 1 3
6 2 4 5 2 6
6 1 2 6
2 0
4 0
1 0
5 8 2 4 1 2 3 4 2 3 4 5 2 5 1 3 3 5
2 0
5 7 2 4 1 2 3 4 4 5 2 5 1 3 3 5
2 1 1 2
1 0
3 1 1 3
4 2 2 4 1 4
3 3 2 3 1 2 1 3
4 6 2 4 1 2 3 4 1 4 2 3 1 3
6 1 3 5
2 1 1 2
5 3 2 4 3 4 1 4
3 0
1 0
4 2 1 2 3 4
1 0
4 4 2 4 1 2 1 3 1 4
2 1 1 2
3 0
5 3 2 3 3 4 1 4
1 0
3 2 2 3 1 3
1 0
5 4 2 3 2 4 2 5 1 4
2 0
4 0
6 1 2 4
3 3 2 3 1 2 1 3
2 1 1 2
4 3 2 3 3 4 1 4
3 3 2 3 1 2 1 3
2 1 1 2
4 6 2 4 1 2 3 4 1 4 2 3 1 3
2 0
4 6 2 4 1 2 3 4 1 4 2 3 1 3
1 0
3 2 1 2 1 3
1 0
4 5 2 4 1 2 3 4 1 4 2 3
2 0
2 1 1 2
5 1 3 5
3 3 2 3 1 2 1 3
2 0
4 5 2 4 1 2 3 4 2 3 1 3
3 0
2 1 1 2
1 0
5 1 1 4
4 5 2 4 1 2 3 4 1 4 1 3
5 8 2 4 1 5 1 4 2 3 4 5 2 5 1 3 3 5
6 4 1 2 2 6 3 5 3 6
2 1 1 2
3 0
2 1 1 2
5 3 2 4 2 5 3 5
6 7 1 5 4 6 2 6 1 6 2 5 1 3 3 5
3 1 2 3
6 3 1 2 1 4 3 6`

var _ = solutionSource

type edge struct{ u, v int }

func norm(u, v int) edge {
	if u > v {
		u, v = v, u
	}
	return edge{u, v}
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func validate(n, m int, edges []edge, out string) error {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid k")
	}
	if k < 0 || len(fields) != 1+2*k {
		return fmt.Errorf("bad edge count")
	}
	orig := make(map[edge]bool)
	degOrig := make([]int, n+1)
	for _, e := range edges {
		orig[norm(e.u, e.v)] = true
		degOrig[e.u]++
		degOrig[e.v]++
	}
	deg := make([]int, n+1)
	used := make(map[edge]bool)
	for i := 0; i < k; i++ {
		u, err1 := strconv.Atoi(fields[1+2*i])
		v, err2 := strconv.Atoi(fields[1+2*i+1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid edge value")
		}
		if u < 1 || u > n || v < 1 || v > n {
			return fmt.Errorf("edge out of range")
		}
		e := norm(u, v)
		if !orig[e] {
			return fmt.Errorf("edge %d %d not in graph", u, v)
		}
		if used[e] {
			return fmt.Errorf("duplicate edge")
		}
		used[e] = true
		deg[u]++
		deg[v]++
	}
	for i := 1; i <= n; i++ {
		req := (degOrig[i] + 1) / 2
		if deg[i] < req {
			return fmt.Errorf("vertex %d degree %d < required %d", i, deg[i], req)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("bad testcase %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		if len(fields) != 2+2*m {
			fmt.Printf("bad testcase %d\n", idx)
			os.Exit(1)
		}
		edges := make([]edge, m)
		for i := 0; i < m; i++ {
			u, _ := strconv.Atoi(fields[2+2*i])
			v, _ := strconv.Atoi(fields[2+2*i+1])
			edges[i] = edge{u, v}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
		}
		input := sb.String()
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if err := validate(n, m, edges, got); err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
