package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const embeddedSolutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Fenwick struct {
   n    int
   tree []int
}

func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, tree: make([]int, n+1)}
}

func (f *Fenwick) Update(i, v int) {
   for x := i; x <= f.n; x += x & -x {
       f.tree[x] += v
   }
}

func (f *Fenwick) Query(i int) int {
   s := 0
   for x := i; x > 0; x -= x & -x {
       s += f.tree[x]
   }
   return s
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   p := make([]int, n)
   absv := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &p[i])
       if p[i] < 0 {
           absv[i] = -p[i]
       } else {
           absv[i] = p[i]
       }
   }
   vals := make([]int, n)
   copy(vals, absv)
   sort.Ints(vals)
   uniq := vals[:0]
   for _, v := range vals {
       if len(uniq) == 0 || uniq[len(uniq)-1] != v {
           uniq = append(uniq, v)
       }
   }
   m := len(uniq)
   rank := make(map[int]int, m)
   for i, v := range uniq {
       rank[v] = i + 1
   }
   a := make([]int64, n)
   bit := NewFenwick(m)
   for i := 0; i < n; i++ {
       r := rank[absv[i]]
       if r > 1 {
           a[i] = int64(bit.Query(r - 1))
       }
       bit.Update(r, 1)
   }
   b := make([]int64, n)
   bit = NewFenwick(m)
   for i := n - 1; i >= 0; i-- {
       r := rank[absv[i]]
       if r > 1 {
           b[i] = int64(bit.Query(r - 1))
       }
       bit.Update(r, 1)
   }
   groups := make(map[int][]int, m)
   for i := 0; i < n; i++ {
       r := rank[absv[i]]
       groups[r] = append(groups[r], i)
   }
   var ans int64
   const INF = int64(4e18)
   for _, idxs := range groups {
       dp := make([]int64, 1)
       dp[0] = 0
       for _, idx := range idxs {
           ci0 := a[idx]
           ci1 := b[idx]
           j := len(dp)
           dp2 := make([]int64, j+1)
           for k := range dp2 {
               dp2[k] = INF
           }
           for ones, cur := range dp {
               if cur >= INF {
                   continue
               }
               c0 := cur + ci0 + int64(ones)
               if c0 < dp2[ones] {
                   dp2[ones] = c0
               }
               c1 := cur + ci1
               if c1 < dp2[ones+1] {
                   dp2[ones+1] = c1
               }
           }
           dp = dp2
       }
       best := INF
       for _, v := range dp {
           if v < best {
               best = v
           }
       }
       ans += best
   }
   fmt.Println(ans)
}`

const testcasesRaw = `
4 -1 -14 5 10
3 -15 -16 -19
7 15 -2 -17 -6 13 14 3
5 -9 -14 -4 -7 -19
11 -4 -3 -8 -10 -1 -2 20 3 -15 18 1
11 4 12 -5 -9 -5 10 -3 -15 15 -1 -20
5 16 -1 12 -8 6
7 18 -2 7 8 -10 -6 -1
5 -18 -15 -18 9 20
5 13 14 10 1 -11
11 -8 -16 6 -8 20 20 8 -3 -9 2 7
12 17 0 20 15 -8 0 -14 -17 -6 -3 17 19
4 -13 1 -9 -2
8 -19 -18 2 -15 -2 0 -19 0
5 0 -11 6 19 -16
5 19 -8 8 -2 -12
5 4 18 -10 1 16
1 3
1 9
3 3 3 -2
10 -14 8 -7 7 -7 -13 -17 -17 -17 -10
10 -11 18 -18 14 11 17 -5 0 -18 -13
9 -2 6 -8 10 -8 -5 8 6 11
1 -6
7 8 -5 7 -7 11 -8 -18
1 -4
5 -5 13 -7 -6 6
5 -11 0 -17 0 16
2 16 5
11 -18 11 4 -15 7 -7 16 -10 1 -2 10
11 0 6 13 -7 -3 1 5 11 -16 -3 20
11 -8 -18 5 19 -12 -3 -17 -10 20 9 16
8 5 4 -7 -20 -7 -10 -20 19
5 -13 5 4 -6 15
1 -8
3 18 1 15
8 13 8 -19 -15 -18 18 -13 11
9 -4 18 -12 -18 3 -15 13 -20 -1
6 -16 -15 14 9 4 -7
5 4 -6 11 5 -14
2 -13 19
6 12 7 6 8 -16 20
4 20 -1 10 7
2 15 -10
6 -10 -9 -11 0 11 1
5 14 -20 -10 -20 -1
2 14 -13
8 18 10 13 -16 13 -5 6 -2
6 -6 -9 20 -20 -17 19
6 14 9 16 -1 12 8
10 19 8 5 -11 -4 18 3 1 -12 7
2 18 -11
11 19 -9 -2 3 -8 16 2 19 -15 -16 5
11 -9 1 3 0 -9 -1 -19 18 -19 13 -15
6 -14 -10 -9 17 11 16
2 -13 -9
11 10 -6 19 -1 5 18 -5 11 -6 -1 3
4 0 14 20 13
8 5 12 5 0 -2 8 6 17
1 -4
3 14 9 15
10 3 5 4 19 -19 -11 12 -16 9 20
6 18 -1 -15 -4 10 -6
11 10 18 -16 -11 -5 -16 -1 -12 -17 -10 5
10 19 15 14 -4 -19 13 -5 -10 -14 -7
1 0
2 -13 -3
1 -2
10 20 -9 -11 6 -12 -15 15 3 -20 15
3 12 6 -11
4 -1 10 12 -16
7 -10 -10 -4 12 5 14 -1
7 1 -10 4 -17 7 -19 -3
11 -19 -1 -11 -15 -10 -13 18 -20 -6 -6 15
1 10
12 14 -9 8 4 1 -9 13 17 -8 -14 10 19
6 -1 6 18 -18 20 -7
5 17 17 -1 10 20
6 -15 -5 0 -14 -18 18
6 13 10 2 -15 -9 -18
8 13 14 18 -5 -18 -8 -16 1
10 -12 -1 -13 12 12 -20 -16 -6 19 -3
10 -19 -18 -20 10 -11 -7 3 -5 2 -1
7 20 18 5 -18 -9 6 11
2 12 12
5 15 -2 -7 7 0
3 -7 -19 10
7 -3 1 6 20 6 0 19
4 -4 -3 12 -15
1 5
5 18 -2 -12 19 0
2 -10 19
2 -7 -4
11 -7 11 -13 1 -20 9 16 -7 -10 -15 -5
9 -14 -13 0 -2 16 17 -11 -16 -9
12 -19 -9 -7 -4 -18 17 0 20 1 18 3 -10
3 -7 -17 -17
3 -17 -4 11
2 -8 -12
7 10 18 14 11 14 0 3
`

var (
	_            = embeddedSolutionSource
	rawTestcases = func() []string {
		scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
		scanner.Buffer(make([]byte, 0, 1024), 1024*1024)
		var cases []string
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				cases = append(cases, line)
			}
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		return cases
	}()
)

type fenwick struct {
	n    int
	tree []int
}

func newFenwick(n int) *fenwick {
	return &fenwick{n: n, tree: make([]int, n+1)}
}

func (f *fenwick) update(i, v int) {
	for x := i; x <= f.n; x += x & -x {
		f.tree[x] += v
	}
}

func (f *fenwick) query(i int) int {
	sum := 0
	for x := i; x > 0; x -= x & -x {
		sum += f.tree[x]
	}
	return sum
}

func solveCase(n int, p []int) string {
	absv := make([]int, n)
	for i, v := range p {
		if v < 0 {
			absv[i] = -v
		} else {
			absv[i] = v
		}
	}
	vals := make([]int, n)
	copy(vals, absv)
	sort.Ints(vals)
	uniq := vals[:0]
	for _, v := range vals {
		if len(uniq) == 0 || uniq[len(uniq)-1] != v {
			uniq = append(uniq, v)
		}
	}
	rank := make(map[int]int, len(uniq))
	for i, v := range uniq {
		rank[v] = i + 1
	}

	a := make([]int64, n)
	b := make([]int64, n)
	bit := newFenwick(len(uniq))
	for i := 0; i < n; i++ {
		r := rank[absv[i]]
		if r > 1 {
			a[i] = int64(bit.query(r - 1))
		}
		bit.update(r, 1)
	}
	bit = newFenwick(len(uniq))
	for i := n - 1; i >= 0; i-- {
		r := rank[absv[i]]
		if r > 1 {
			b[i] = int64(bit.query(r - 1))
		}
		bit.update(r, 1)
	}

	groups := make(map[int][]int, len(uniq))
	for i := 0; i < n; i++ {
		r := rank[absv[i]]
		groups[r] = append(groups[r], i)
	}

	var ans int64
	const INF = int64(4e18)
	for _, idxs := range groups {
		dp := make([]int64, 1)
		dp[0] = 0
		for _, idx := range idxs {
			ci0, ci1 := a[idx], b[idx]
			dp2 := make([]int64, len(dp)+1)
			for i := range dp2 {
				dp2[i] = INF
			}
			for ones, cur := range dp {
				if cur >= INF {
					continue
				}
				c0 := cur + ci0 + int64(ones)
				if c0 < dp2[ones] {
					dp2[ones] = c0
				}
				c1 := cur + ci1
				if c1 < dp2[ones+1] {
					dp2[ones+1] = c1
				}
			}
			dp = dp2
		}
		best := INF
		for _, v := range dp {
			if v < best {
				best = v
			}
		}
		ans += best
	}
	return strconv.FormatInt(ans, 10)
}

func parseCase(line string) (int, []int, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) < 1 {
		return 0, nil, fmt.Errorf("invalid line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, nil, err
	}
	if len(fields) != 1+n {
		return 0, nil, fmt.Errorf("expected %d numbers got %d", 1+n, len(fields))
	}
	p := make([]int, n)
	for i := 0; i < n; i++ {
		val, _ := strconv.Atoi(fields[1+i])
		p[i] = val
	}
	return n, p, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, line := range rawTestcases {
		n, p, err := parseCase(line)
		if err != nil {
			fmt.Printf("case %d invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		expected := solveCase(n, p)
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i, v := range p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(rawTestcases))
}
