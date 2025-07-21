package main

import (
   "bufio"
   "fmt"
   "os"
)

// UnionFind for n elements
type UF struct {
   p, r, sz []int
}

func NewUF(n int) *UF {
   p := make([]int, n)
   r := make([]int, n)
   sz := make([]int, n)
   for i := 0; i < n; i++ {
       p[i] = i
       sz[i] = 1
   }
   return &UF{p: p, r: r, sz: sz}
}

func (u *UF) Find(x int) int {
   if u.p[x] != x {
       u.p[x] = u.Find(u.p[x])
   }
   return u.p[x]
}

func (u *UF) Union(a, b int) {
   x := u.Find(a)
   y := u.Find(b)
   if x == y {
       return
   }
   if u.r[x] < u.r[y] {
       x, y = y, x
   }
   u.p[y] = x
   u.sz[x] += u.sz[y]
   if u.r[x] == u.r[y] {
       u.r[x]++
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(in, &n, &m)
   uf := NewUF(n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--
       v--
       uf.Union(u, v)
   }
   // count component sizes
   comp := make(map[int]int)
   for i := 0; i < n; i++ {
       r := uf.Find(i)
       comp[r] = uf.sz[r]
   }
   // freq of sizes
   freq := make(map[int]int)
   for _, sz := range comp {
       freq[sz]++
   }
   // generate lucky numbers <= n
   luckies := genLuckies(n)
   // check zero roads
   for _, L := range luckies {
       if freq[L] > 0 {
           fmt.Println(0)
           return
       }
   }
   const INF = 1e9
   dp := make([]int, n+1)
   for i := 1; i <= n; i++ {
       dp[i] = INF
   }
   // bounded knapsack with cost = components count
   for s, f := range freq {
       old := make([]int, n+1)
       copy(old, dp)
       // process remainder classes
       for r0 := 0; r0 < s; r0++ {
           // deque of pairs (k, mval)
           type pair struct{ k, v int }
           dq := make([]pair, 0)
           // iterate k: w = r0 + k*s
           for k, w := 0, r0; w <= n; k, w = k+1, w+s {
               mv := old[w] - k
               // push mv
               for len(dq) > 0 && dq[len(dq)-1].v >= mv {
                   dq = dq[:len(dq)-1]
               }
               dq = append(dq, pair{k, mv})
               // pop old
               if dq[0].k < k-f {
                   dq = dq[1:]
               }
               // update dp
               dp[w] = min(dp[w], dq[0].v + k)
           }
       }
   }
   // find answer
   ans := INF
   for _, L := range luckies {
       if L <= n && dp[L] < INF {
           // roads = components -1
           ans = min(ans, dp[L]-1)
       }
   }
   if ans == INF {
       fmt.Println(-1)
   } else {
       fmt.Println(ans)
   }
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

// generate lucky numbers <= limit
func genLuckies(limit int) []int {
   var res []int
   var dfs func(int)
   dfs = func(x int) {
       if x > limit {
           return
       }
       if x > 0 {
           res = append(res, x)
       }
       dfs(x*10 + 4)
       dfs(x*10 + 7)
   }
   dfs(0)
   return res
}
