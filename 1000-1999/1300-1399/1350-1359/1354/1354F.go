package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   rdr := bufio.NewReader(os.Stdin)
   wr := bufio.NewWriter(os.Stdout)
   defer wr.Flush()

   var t int
   fmt.Fscan(rdr, &t)
   for t > 0 {
       t--
       solve(rdr, wr)
   }
}

func solve(rdr *bufio.Reader, wr *bufio.Writer) {
   var n, k int
   fmt.Fscan(rdr, &n, &k)
   a := make([]int64, n+1)
   b := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(rdr, &a[i], &b[i])
   }
   vec := make([]int, n)
   for i := 0; i < n; i++ {
       vec[i] = i + 1
   }
   sort.Slice(vec, func(i, j int) bool {
       return b[vec[i]] < b[vec[j]]
   })
   const INF = int64(1e18)
   // dp[i][j]: using first i, choose j
   dp := make([][]int64, n+1)
   lst := make([][]bool, n+1)
   for i := 0; i <= n; i++ {
       dp[i] = make([]int64, k+1)
       lst[i] = make([]bool, k+1)
       for j := 0; j <= k; j++ {
           dp[i][j] = -INF
       }
   }
   dp[0][0] = 0
   for i := 1; i <= n; i++ {
       p := vec[i-1]
       for j := 0; j <= k; j++ {
           // skip p
           dp[i][j] = dp[i-1][j] + int64(k-1)*b[p]
           lst[i][j] = false
           // take p
           if j > 0 {
               v := dp[i-1][j-1] + b[p]*int64(j-1) + a[p]
               if v > dp[i][j] {
                   dp[i][j] = v
                   lst[i][j] = true
               }
           }
       }
   }
   // backtrack
   e := make([]int, 0, k)
   f := make([]int, 0, 2*(n-k))
   p2 := k
   for i := n; i >= 1; i-- {
       if !lst[i][p2] {
           id := vec[i-1]
           f = append(f, -id)
           f = append(f, id)
       } else {
           e = append(e, vec[i-1])
           p2--
       }
   }
   // reverse e and f
   for i, j := 0, len(e)-1; i < j; i, j = i+1, j-1 {
       e[i], e[j] = e[j], e[i]
   }
   for i, j := 0, len(f)-1; i < j; i, j = i+1, j-1 {
       f[i], f[j] = f[j], f[i]
   }
   total := len(e) + len(f)
   fmt.Fprintln(wr, total)
   // print first k-1 from e
   for i := 0; i < len(e)-1; i++ {
       fmt.Fprint(wr, e[i], " ")
   }
   // then f
   for _, x := range f {
       fmt.Fprint(wr, x, " ")
   }
   // last of e
   if len(e) > 0 {
       fmt.Fprint(wr, e[len(e)-1])
   }
   fmt.Fprintln(wr)
}
