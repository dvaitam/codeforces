package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const solution1312GSource = `package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Fast input reader
type FastReader struct {
   r *bufio.Reader
}

func NewReader() *FastReader {
   return &FastReader{r: bufio.NewReader(os.Stdin)}
}

func (fr *FastReader) ReadInt() (int, error) {
   var x int
   var c byte
   var err error
   // skip non-digit
   for {
       c, err = fr.r.ReadByte()
       if err != nil {
           return 0, err
       }
       if c == '-' || (c >= '0' && c <= '9') {
           break
       }
   }
   neg := false
   if c == '-' {
       neg = true
       c, _ = fr.r.ReadByte()
   }
   for ; err == nil && c >= '0' && c <= '9'; c, err = fr.r.ReadByte() {
       x = x*10 + int(c-'0')
   }
   if neg {
       x = -x
   }
   return x, nil
}

func (fr *FastReader) ReadToken() (byte, error) {
   var c byte
   var err error
   for {
       c, err = fr.r.ReadByte()
       if err != nil {
           return 0, err
       }
       if c >= 'a' && c <= 'z' {
           return c, nil
       }
   }
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   fr := NewReader()
   n, _ := fr.ReadInt()
   parent := make([]int, n+1)
   parent[0] = -1
   cArr := make([]byte, n+1)
   for i := 1; i <= n; i++ {
       p, _ := fr.ReadInt()
       c, _ := fr.ReadToken()
       parent[i] = p
       cArr[i] = c
   }
   k, _ := fr.ReadInt()
   a := make([]int, k)
   mark := make([]bool, n+1)
   for i := 0; i < k; i++ {
       ai, _ := fr.ReadInt()
       a[i] = ai
       mark[ai] = true
   }
   // build adjacency flat arrays
   deg := make([]int, n+1)
   for i := 1; i <= n; i++ {
       deg[parent[i]]++
   }
   pos := make([]int, n+1)
   for u := 1; u <= n; u++ {
       pos[u] = pos[u-1] + deg[u-1]
   }
   childTo := make([]int, n)
   childC := make([]byte, n)
   cur := make([]int, n+1)
   copy(cur, pos)
   for i := 1; i <= n; i++ {
       p := parent[i]
       idx := cur[p]
       childTo[idx] = i
       childC[idx] = cArr[i]
       cur[p]++
   }
   // sort children by char
   for u := 0; u <= n; u++ {
       st := pos[u]
       cnt := deg[u]
       if cnt > 1 {
           sort.Slice(childTo[st:st+cnt], func(i, j int) bool {
               return childC[st+i] < childC[st+j]
           })
       }
   }
   // first DFS: compute lex order and intervals
   const INF = int(1e9)
   l := make([]int, n+1)
   r := make([]int, n+1)
   for i := range l {
       l[i] = INF
   }
   order := make([]int, n+1)
   cntOrder := 0
   type Frame struct{ u, idx int }
   stack := make([]Frame, 0, n+1)
   stack = append(stack, Frame{u: 0, idx: 0})
   for len(stack) > 0 {
       fr := &stack[len(stack)-1]
       u := fr.u
       if fr.idx == 0 {
           if mark[u] {
               cntOrder++
               order[u] = cntOrder
               l[u], r[u] = cntOrder, cntOrder
           }
       }
       if fr.idx < deg[u] {
           v := childTo[pos[u]+fr.idx]
           fr.idx++
           stack = append(stack, Frame{u: v, idx: 0})
       } else {
           // post
           if u != 0 {
               p := parent[u]
               if l[u] <= r[u] {
                   if l[p] > l[u] {
                       l[p] = l[u]
                   }
                   if r[p] < r[u] {
                       r[p] = r[u]
                   }
               }
           }
           stack = stack[:len(stack)-1]
       }
   }
   // dp with second DFS
   dp := make([]int, n+1)
   type Frame2 struct{ u, idx int; act bool }
   stack2 := make([]Frame2, 0, n+1)
   valStack := make([]int, 0, n+1)
   minStack := make([]int, 0, n+1)
   // init root
   dp[0] = 0
   stack2 = append(stack2, Frame2{u: 0, idx: 0, act: true})
   // push root val
   rootVal := dp[0] - l[0] + 1
   valStack = append(valStack, rootVal)
   minStack = append(minStack, rootVal)
   for len(stack2) > 0 {
       fr2 := &stack2[len(stack2)-1]
       u := fr2.u
       if fr2.idx == 0 && u != 0 {
           p := parent[u]
           dp[u] = dp[p] + 1
           if mark[u] {
               // autocomplete
               best := minStack[len(minStack)-1]
               cand := order[u] + best
               if cand < dp[u] {
                   dp[u] = cand
               }
           }
           if l[u] <= r[u] {
               val := dp[u] - l[u] + 1
               stack2[len(stack2)-1].act = true
               valStack = append(valStack, val)
               // update minStack
               m := val
               if last := minStack[len(minStack)-1]; last < m {
                   m = last
               }
               minStack = append(minStack, m)
           }
       }
       if fr2.idx < deg[u] {
           v := childTo[pos[u]+fr2.idx]
           stack2[len(stack2)-1].idx++
           stack2 = append(stack2, Frame2{u: v, idx: 0, act: false})
           continue
       }
       // exit
       if fr2.act {
           valStack = valStack[:len(valStack)-1]
           minStack = minStack[:len(minStack)-1]
       }
       stack2 = stack2[:len(stack2)-1]
   }
   // output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i, ai := range a {
       if i > 0 {
           w.WriteByte(' ')
       }
       fmt.Fprintf(w, "%d", dp[ai])
   }
}
`

// Keep the embedded reference solution reachable.
var _ = solution1312GSource

type testCase struct {
	n      int
	parent []int
	char   []byte
	nodes  []int
}

const testcasesRaw = `5 0 b 1 a 0 a 3 c 2 b 1 3
4 0 c 1 c 2 c 0 a 3 1 2 3
1 0 b 1 1
1 0 c 1 1
1 0 c 1 1
2 0 c 1 a 2 2 1
5 0 a 0 a 2 c 2 b 1 b 4 2 5 3 1
5 0 a 0 a 2 b 2 b 3 c 4 3 2 5 4
1 0 b 1 1
4 0 a 1 b 1 a 3 b 3 4 1 3
3 0 b 0 a 0 c 1 2
5 0 a 1 c 1 b 2 b 3 c 4 3 4 5 2
4 0 c 1 b 1 c 3 b 3 2 3 4
5 0 b 1 b 2 a 2 c 0 c 4 1 2 5 3
3 0 b 0 a 0 b 3 1 2 3
1 0 a 1 1
3 0 c 0 c 1 a 1 3
1 0 a 1 1
1 0 c 1 1
2 0 c 1 b 1 1
2 0 b 1 b 2 2 1
3 0 b 1 a 0 c 2 2 3
4 0 b 1 a 2 a 1 c 2 1 3
2 0 b 0 a 1 2
3 0 a 0 b 2 a 2 3 2
1 0 a 1 1
2 0 a 1 c 1 2
2 0 c 0 b 1 2
2 0 c 1 a 2 2 1
1 0 c 1 1
5 0 b 0 a 0 c 1 b 0 a 1 1
4 0 a 0 c 1 c 0 a 2 4 1
5 0 a 1 a 2 c 0 a 3 a 1 2
1 0 c 1 1
3 0 a 0 c 2 c 1 2
5 0 b 1 c 0 c 0 b 2 c 1 2
1 0 c 1 1
1 0 a 1 1
3 0 b 0 b 1 a 3 3 2 1
2 0 b 0 a 2 2 1
3 0 a 1 a 0 b 1 1
5 0 a 0 a 0 c 3 b 3 c 5 4 2 3 5 1
4 0 b 1 b 1 a 2 b 3 2 1 4
2 0 b 0 c 1 1
4 0 a 0 c 0 c 3 b 3 4 2 1
3 0 a 0 c 2 c 3 1 3 2
1 0 a 1 1
2 0 a 0 a 1 2
4 0 c 0 b 0 c 2 c 2 3 4
3 0 b 0 a 2 a 2 3 2
3 0 a 1 c 0 b 3 2 1 3
2 0 b 1 c 2 1 2
1 0 b 1 1
5 0 b 0 b 0 b 1 c 1 b 3 3 5 2
5 0 b 0 c 2 a 0 c 3 c 5 3 2 1 4 5
2 0 a 0 b 1 2
3 0 a 0 c 0 a 2 3 1
3 0 b 0 b 1 c 3 2 1 3
4 0 c 0 c 1 b 1 c 1 4
5 0 c 1 c 1 b 3 c 2 a 5 4 2 5 1 3
5 0 a 0 a 1 b 0 a 2 b 1 1
3 0 b 1 b 2 c 3 1 3 2
4 0 c 1 b 0 a 0 c 4 2 1 4 3
1 0 a 1 1
4 0 b 0 a 2 b 1 a 2 1 4
2 0 b 1 c 2 2 1
5 0 a 0 c 0 a 1 b 1 c 3 4 3 5
5 0 c 0 c 0 b 1 b 4 c 1 1
3 0 b 1 b 1 b 2 3 2
5 0 b 0 b 1 c 3 a 0 b 5 2 3 4 5 1
3 0 b 1 b 0 a 3 1 2 3
1 0 c 1 1
2 0 b 0 b 1 1
3 0 c 1 c 2 b 2 2 1
5 0 b 1 a 1 b 0 a 4 a 1 1
1 0 c 1 1
4 0 b 0 a 1 a 0 a 1 2
3 0 a 1 b 1 a 2 3 2
2 0 c 0 a 2 1 2
1 0 c 1 1
5 0 c 0 c 2 c 3 b 1 b 2 1 2
3 0 c 1 a 2 a 1 2
1 0 c 1 1
4 0 a 1 c 2 a 1 c 2 3 2
3 0 c 0 b 0 a 2 1 2
3 0 c 1 b 2 b 2 2 3
5 0 b 0 c 2 c 2 c 1 c 4 5 2 3 1
2 0 a 1 b 1 1
5 0 a 0 b 1 a 2 b 1 b 5 1 3 5 2 4
1 0 a 1 1
3 0 a 0 b 0 c 2 2 1
5 0 b 0 a 0 c 1 c 1 b 2 3 2
5 0 b 1 c 0 a 0 a 4 c 4 1 5 4 3
3 0 b 1 b 2 a 2 3 2
3 0 b 0 a 2 a 2 3 2
1 0 b 1 1
1 0 a 1 1
2 0 b 0 a 1 2
1 0 b 1 1
4 0 a 0 b 2 c 2 a 2 4 2
`

func parseTestcases() []testCase {
	fields := strings.Fields(testcasesRaw)
	var res []testCase
	for i := 0; i < len(fields); {
		n, _ := strconv.Atoi(fields[i])
		i++
		parent := make([]int, n+1)
		char := make([]byte, n+1)
		for v := 1; v <= n; v++ {
			if i+1 >= len(fields) {
				return res
			}
			p, _ := strconv.Atoi(fields[i])
			parent[v] = p
			i++
			char[v] = fields[i][0]
			i++
		}
		if i >= len(fields) {
			return res
		}
		k, _ := strconv.Atoi(fields[i])
		i++
		nodes := make([]int, k)
		for j := 0; j < k; j++ {
			if i >= len(fields) {
				return res
			}
			nodes[j], _ = strconv.Atoi(fields[i])
			i++
		}
		res = append(res, testCase{n: n, parent: parent, char: char, nodes: nodes})
	}
	return res
}

func solveExpected(tc testCase) []int {
	n := tc.n
	parent := tc.parent
	cArr := tc.char
	mark := make([]bool, n+1)
	for _, v := range tc.nodes {
		mark[v] = true
	}
	deg := make([]int, n+1)
	for i := 1; i <= n; i++ {
		deg[parent[i]]++
	}
	pos := make([]int, n+1)
	for u := 1; u <= n; u++ {
		pos[u] = pos[u-1] + deg[u-1]
	}
	childTo := make([]int, n)
	childC := make([]byte, n)
	cur := make([]int, n+1)
	copy(cur, pos)
	for i := 1; i <= n; i++ {
		p := parent[i]
		idx := cur[p]
		childTo[idx] = i
		childC[idx] = cArr[i]
		cur[p]++
	}
	for u := 0; u <= n; u++ {
		st := pos[u]
		cnt := deg[u]
		if cnt > 1 {
			sort.Slice(childTo[st:st+cnt], func(i, j int) bool {
				return childC[st+i] < childC[st+j]
			})
		}
	}
	const INF = int(1e9)
	l := make([]int, n+1)
	r := make([]int, n+1)
	for i := range l {
		l[i] = INF
	}
	order := make([]int, n+1)
	cntOrder := 0
	type frame struct{ u, idx int }
	stack := []frame{{u: 0, idx: 0}}
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		u := top.u
		if top.idx == 0 && mark[u] {
			cntOrder++
			order[u] = cntOrder
			l[u], r[u] = cntOrder, cntOrder
		}
		if top.idx < deg[u] {
			v := childTo[pos[u]+top.idx]
			top.idx++
			stack = append(stack, frame{u: v, idx: 0})
		} else {
			if u != 0 && l[u] <= r[u] {
				p := parent[u]
				if l[p] > l[u] {
					l[p] = l[u]
				}
				if r[p] < r[u] {
					r[p] = r[u]
				}
			}
			stack = stack[:len(stack)-1]
		}
	}
	dp := make([]int, n+1)
	type frame2 struct {
		u, idx int
		act    bool
	}
	stack2 := []frame2{{u: 0, idx: 0, act: true}}
	valStack := []int{dp[0] - l[0] + 1}
	minStack := []int{valStack[0]}
	for len(stack2) > 0 {
		top := &stack2[len(stack2)-1]
		u := top.u
		if top.idx == 0 && u != 0 {
			p := parent[u]
			dp[u] = dp[p] + 1
			if mark[u] {
				best := minStack[len(minStack)-1]
				if cand := order[u] + best; cand < dp[u] {
					dp[u] = cand
				}
			}
			if l[u] <= r[u] {
				val := dp[u] - l[u] + 1
				stack2[len(stack2)-1].act = true
				valStack = append(valStack, val)
				m := val
				if last := minStack[len(minStack)-1]; last < m {
					m = last
				}
				minStack = append(minStack, m)
			}
		}
		if top.idx < deg[u] {
			v := childTo[pos[u]+top.idx]
			top.idx++
			stack2 = append(stack2, frame2{u: v, idx: 0, act: false})
			continue
		}
		if top.act {
			valStack = valStack[:len(valStack)-1]
			minStack = minStack[:len(minStack)-1]
		}
		stack2 = stack2[:len(stack2)-1]
	}
	res := make([]int, len(tc.nodes))
	for i, v := range tc.nodes {
		res[i] = dp[v]
	}
	return res
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for i := 1; i <= tc.n; i++ {
		fmt.Fprintf(&sb, "%d %c\n", tc.parent[i], tc.char[i])
	}
	fmt.Fprintf(&sb, "%d\n", len(tc.nodes))
	for i, v := range tc.nodes {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, idx int, tc testCase) error {
	expectArr := solveExpected(tc)
	expect := make([]string, len(expectArr))
	for i, v := range expectArr {
		expect[i] = strconv.Itoa(v)
	}
	expectStr := strings.Join(expect, " ")

	input := buildInput(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("case %d failed: %v\nstderr: %s\ninput:\n%s", idx, err, string(out), input)
	}
	got := strings.TrimSpace(string(out))
	if got != expectStr {
		return fmt.Errorf("case %d failed\nexpected: %s\ngot: %s\ninput:\n%s", idx, expectStr, got, input)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases := parseTestcases()
	for i, tc := range testcases {
		if err := runCase(bin, i+1, tc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
