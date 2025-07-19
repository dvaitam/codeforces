package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Query represents an interval with probability
type Query struct {
   l, r int
   p    float64
}

var (
   n, m    int
   a       []int
   queries []Query
   children [][]int
   b, lim  []int
   dp      [][]float64
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

// dfs processes node x, builds dp[x], b[x], lim[x]
func dfs(x int) {
   // initialize dp[x] = [1]
   dp[x] = make([]float64, 1)
   dp[x][0] = 1.0
   lim[x] = 0
   // compute base max b[x] over gaps
   r := queries[x].r
   for _, v := range children[x] {
       for j := queries[v].r + 1; j <= r; j++ {
           if a[j] > b[x] {
               b[x] = a[j]
           }
       }
       r = queries[v].l - 1
   }
   for j := queries[x].l; j <= r; j++ {
       if a[j] > b[x] {
           b[x] = a[j]
       }
   }
   // combine with children
   for _, v := range children[x] {
       dfs(v)
       // convolution of dp[x] and dp[v]
       rb := max(b[x], b[v])
       // old lengths
       lx := lim[x]
       lv := lim[v]
       // new limit before own p
       rl1 := b[x] + lx - rb
       rl2 := b[v] + lv - rb
       rl := max(rl1, rl2)
       nd := make([]float64, rl+1)
       for j := 0; j <= lv; j++ {
           for k := 0; k <= lx; k++ {
               // max of two values
               v1 := b[v] + j
               v2 := b[x] + k
               mx := v1
               if v2 > mx {
                   mx = v2
               }
               ind := mx - rb
               nd[ind] += dp[v][j] * dp[x][k]
           }
       }
       // update x
       lim[x] = rl
       b[x] = rb
       dp[x] = make([]float64, rl+1)
       copy(dp[x], nd)
   }
   // apply this node's probability
   // extend dp by 1
   old := dp[x]
   oldLen := len(old)
   newdp := make([]float64, oldLen+1)
   p := queries[x].p
   for i := 0; i < oldLen+1; i++ {
       if i < oldLen {
           newdp[i] += old[i] * (1 - p)
       }
       if i > 0 {
           newdp[i] += old[i-1] * p
       }
   }
   dp[x] = newdp
   lim[x] = len(newdp) - 1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a = make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // read queries
   queries = make([]Query, m+1)
   // root covers full range
   queries[0] = Query{l: 1, r: n, p: 0}
   for i := 1; i <= m; i++ {
       fmt.Fscan(reader, &queries[i].l, &queries[i].r, &queries[i].p)
   }
   // sort intervals: l asc, r desc
   sort.Slice(queries, func(i, j int) bool {
       if queries[i].l != queries[j].l {
           return queries[i].l < queries[j].l
       }
       return queries[i].r > queries[j].r
   })
   total := m + 1
   children = make([][]int, total)
   b = make([]int, total)
   lim = make([]int, total)
   dp = make([][]float64, total)
   // build tree using stack
   stk := make([]int, 0, total)
   stk = append(stk, 0)
   for i := 1; i < total; i++ {
       // pop until current fits
       for len(stk) > 0 && queries[i].r > queries[stk[len(stk)-1]].r {
           stk = stk[:len(stk)-1]
       }
       par := stk[len(stk)-1]
       children[par] = append(children[par], i)
       // chain equal intervals
       for i+1 < total && queries[i+1].l == queries[i].l && queries[i+1].r == queries[i].r {
           children[i] = append(children[i], i+1)
           i++
       }
       stk = append(stk, i)
   }
   // dfs from root
   dfs(0)
   // compute answer
   ans := 0.0
   for i := 0; i <= lim[0]; i++ {
       ans += dp[0][i] * float64(b[0]+i)
   }
   // output with 6 decimal places
   fmt.Printf("%.6f\n", ans)
}
