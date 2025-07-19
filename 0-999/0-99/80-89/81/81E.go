package main

import (
   "bufio"
   "fmt"
   "os"
)

// Pair stores a pair of integers for DP values
type Pair struct{ first, second int }

// add returns the sum of two pairs
func (a Pair) add(b Pair) Pair { return Pair{a.first + b.first, a.second + b.second} }
// sub returns the difference of two pairs
func (a Pair) sub(b Pair) Pair { return Pair{a.first - b.first, a.second - b.second} }
// less returns true if a < b in lexicographical order
func (a Pair) less(b Pair) bool {
   if a.first != b.first {
       return a.first < b.first
   }
   return a.second < b.second
}

var (
   n       int
   f, s    []int
   tag, pre []int
   c       [][]int
   dp, cdp [][]Pair
   cpre    []int
   ans     Pair
   todo    []Pair
)

// gao performs DP on tree rooted at v
func gao(v int) {
   if tag[v] == 0 {
       tag[v] = 1
   }
   dp[v][0] = Pair{0, 0}
   dp[v][1] = Pair{0, 0}
   pre[v] = -1
   for _, w := range c[v] {
       if tag[w] == -1 {
           continue
       }
       gao(w)
       // not pair v->w
       dp[v][0] = dp[v][0].add(dp[w][1])
       // consider pairing v->w
       diff := dp[w][1].sub(dp[w][0])
       base := Pair{1, 0}
       if s[v] != s[w] {
           base.second = 1
       }
       tmp := base.sub(diff)
       if dp[v][1].less(tmp) {
           dp[v][1] = tmp
           pre[v] = w
       }
   }
   dp[v][1] = dp[v][1].add(dp[v][0])
}

// dump collects the chosen pairs in subtree v
func dump(v int, flag bool, ret *[]Pair) {
   for _, w := range c[v] {
       if tag[w] == -1 {
           continue
       }
       if !flag || w != pre[v] {
           dump(w, true, ret)
       } else {
           *ret = append(*ret, Pair{v, w})
           dump(w, false, ret)
       }
   }
}

// solve processes one connected component starting at v
func solve(v0 int) {
   v := v0
   // find cycle
   for tag[v] == 0 {
       tag[v] = 1
       v = f[v]
   }
   // collect cycle nodes
   circ := []int{v}
   for u := f[v]; u != v; u = f[u] {
       circ = append(circ, u)
   }
   // mark cycle nodes
   for _, u := range circ {
       tag[u] = -1
   }
   // run gao on cycle nodes
   for _, u := range circ {
       gao(u)
   }
   // DP on cycle: case 1, no edge between last and first
   m := len(circ)
   // initialize
   cdp = cdp
   cpre = cpre
   // first node
   cdp[circ[0]][0] = dp[circ[0]][0]
   cdp[circ[0]][1] = dp[circ[0]][1]
   cpre[circ[0]] = -1
   for i := 1; i < m; i++ {
       u := circ[i]
       p := circ[i-1]
       cdp[u][0] = cdp[p][1].add(dp[u][0])
       cdp[u][1] = cdp[p][1].add(dp[u][1])
       cpre[u] = -1
       // consider pairing p->u: cdp[p][0] + (1, delta) + dp[u][0]
       base := Pair{1, 0}
       if s[p] != s[u] {
           base.second = 1
       }
       tmp := cdp[p][0].add(base).add(dp[u][0])
       if cdp[u][1].less(tmp) {
           cdp[u][1] = tmp
           cpre[u] = 1
       }
   }
   best := cdp[circ[m-1]][1]
   var how []Pair
   // backtrack case 1
   for i := m - 1; i >= 0; {
       u := circ[i]
       if cpre[u] == -1 {
           dump(u, true, &how)
           i--
       } else {
           // paired with previous
           prev := circ[i-1]
           dump(u, false, &how)
           dump(prev, false, &how)
           how = append(how, Pair{u, prev})
           i -= 2
       }
   }
   // case 2: allow last->first, rotate cycle
   // rotate circ
   circ = append(circ[1:], circ[0])
   // reinit dp arrays for rotated
   cdp[circ[0]][0] = dp[circ[0]][0]
   cdp[circ[0]][1] = dp[circ[0]][1]
   cpre[circ[0]] = -1
   for i := 1; i < m; i++ {
       u := circ[i]
       p := circ[i-1]
       cdp[u][0] = cdp[p][1].add(dp[u][0])
       cdp[u][1] = cdp[p][1].add(dp[u][1])
       cpre[u] = -1
       base := Pair{1, 0}
       if s[p] != s[u] {
           base.second = 1
       }
       tmp := cdp[p][0].add(base).add(dp[u][0])
       if cdp[u][1].less(tmp) {
           cdp[u][1] = tmp
           cpre[u] = 1
       }
   }
   // check case 2
   cand := cdp[circ[m-1]][1]
   if best.less(cand) {
       best = cand
       how = how[:0]
       for i := m - 1; i >= 0; {
           u := circ[i]
           if cpre[u] == -1 {
               dump(u, true, &how)
               i--
           } else {
               prev := circ[i-1]
               dump(u, false, &how)
               dump(prev, false, &how)
               how = append(how, Pair{u, prev})
               i -= 2
           }
       }
   }
   ans = ans.add(best)
   todo = append(todo, how...)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n)
   f = make([]int, n+1)
   s = make([]int, n+1)
   tag = make([]int, n+1)
   pre = make([]int, n+1)
   c = make([][]int, n+1)
   dp = make([][]Pair, n+1)
   cdp = make([][]Pair, n+1)
   cpre = make([]int, n+1)
   for i := 1; i <= n; i++ {
       dp[i] = make([]Pair, 2)
       cdp[i] = make([]Pair, 2)
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &f[i], &s[i])
       c[f[i]] = append(c[f[i]], i)
   }
   ans = Pair{0, 0}
   for i := 1; i <= n; i++ {
       if tag[i] == 0 {
           solve(i)
       }
   }
   fmt.Fprintf(writer, "%d %d\n", ans.first, ans.second)
   for _, p := range todo {
       fmt.Fprintf(writer, "%d %d\n", p.first, p.second)
   }
}
