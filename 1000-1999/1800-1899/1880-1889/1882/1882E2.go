package main

import (
   "bufio"
   "fmt"
   "os"
)

// calc computes the number of operations for array a of size n with shift w
func calc(n int, a []int, w int) int {
   p := make([]int, n+1)
   vis := make([]bool, n+1)
   for i := 0; i <= n; i++ {
       p[i] = (a[i] + w) % (n + 1)
   }
   ans := n
   work(p, vis, 0, true)
   for i := 1; i <= n; i++ {
       if p[i] == i {
           ans--
       } else if !vis[i] {
           work(p, vis, i, false)
           ans++
       }
   }
   return ans
}

// work marks all elements in the cycle of x in permutation p
func work(p []int, vis []bool, x int, reset bool) {
   if reset {
       for i := range vis {
           vis[i] = false
       }
   }
   vis[x] = true
   for y := p[x]; y != x; y = p[y] {
       vis[y] = true
   }
}

// dfs records the sequence of shifts to bring p[x] to y
func dfs(p []int, v *[]int, x, y int) {
   if x == y {
       return
   }
   dfs(p, v, p[x], y)
   *v = append(*v, x-p[x])
}

// solve constructs the sequence of shifts achieving minimal operations
func solve(n int, a []int, w int) []int {
   p := make([]int, n+1)
   vis := make([]bool, n+1)
   v := make([]int, 0)
   for i := 0; i <= n; i++ {
       p[i] = (a[i] + w) % (n + 1)
   }
   cur := 0
   work(p, vis, 0, true)
   for i := 1; i <= n; i++ {
       if p[i] != i && !vis[i] {
           v = append(v, i-cur)
           p[cur], p[i] = p[i], p[cur]
           cur = i
           work(p, vis, cur, true)
       }
   }
   dfs(p, &v, p[cur], cur)
   return v
}

// F maps a possibly negative shift to 1-based index
func F(n, x int) int {
   if x > 0 {
       return x
   }
   return n + 1 + x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, n+1)
   b := make([]int, m+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for j := 1; j <= m; j++ {
       fmt.Fscan(reader, &b[j])
   }
   f1 := make([]int, n+1)
   f2 := make([]int, m+1)
   for i := 0; i <= n; i++ {
       f1[i] = calc(n, a, i)
   }
   for j := 0; j <= m; j++ {
       f2[j] = calc(m, b, j)
   }
   ans := int(1e9)
   bi, bj := -1, -1
   for i := 0; i <= n; i++ {
       for j := 0; j <= m; j++ {
           if (f1[i]^f2[j])%2 == 0 {
               mx := f1[i]
               if f2[j] > mx {
                   mx = f2[j]
               }
               if mx < ans {
                   ans = mx
                   bi, bj = i, j
               }
           }
       }
   }
   if bi < 0 {
       fmt.Println(-1)
       return
   }
   v1 := solve(n, a, bi)
   v2 := solve(m, b, bj)
   // pad operations to equal length
   for len(v1) < len(v2) {
       v1 = append(v1, 1, n)
   }
   for len(v2) < len(v1) {
       v2 = append(v2, 1, m)
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, len(v1))
   for k := range v1 {
       fmt.Fprintln(w, F(n, v1[k]), F(m, v2[k]))
   }
}
