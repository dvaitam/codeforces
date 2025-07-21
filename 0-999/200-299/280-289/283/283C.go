package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, q, t int
   fmt.Fscan(reader, &n, &q, &t)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   indegree := make([]int, n)
   out := make([]int, n)
   nxt := make([]int, n)
   for i := 0; i < n; i++ {
       nxt[i] = -1
   }
   prev := make([]int, n)
   for i := 0; i < n; i++ {
       prev[i] = -1
   }
   const MOD = 1000000007
   for i := 0; i < q; i++ {
       var b, c int
       fmt.Fscan(reader, &b, &c)
       b--
       c--
       out[b]++
       indegree[c]++
       if out[b] > 1 || indegree[c] > 1 {
           fmt.Println(0)
           return
       }
       nxt[b] = c
       prev[c] = b
   }
   visited := make([]bool, n)
   weights := make([]int, 0, n)
   C := 0
   for i := 0; i < n; i++ {
       if indegree[i] == 0 {
           path := []int{}
           cur := i
           for cur != -1 {
               path = append(path, cur)
               visited[cur] = true
               cur = nxt[cur]
           }
           k := len(path)
           prefixSum := make([]int, k)
           sum := 0
           for j, v := range path {
               sum += a[v]
               prefixSum[j] = sum
           }
           for _, w := range prefixSum {
               weights = append(weights, w)
           }
           for j, v := range path {
               C += (k-1-j) * a[v]
           }
       }
   }
   for i := 0; i < n; i++ {
       if !visited[i] {
           fmt.Println(0)
           return
       }
   }
   T := t - C
   if T < 0 {
       fmt.Println(0)
       return
   }
   dp := make([]int, T+1)
   dp[0] = 1
   for _, w := range weights {
       for s := w; s <= T; s++ {
           dp[s] += dp[s-w]
           if dp[s] >= MOD {
               dp[s] -= MOD
           }
       }
   }
   fmt.Println(dp[T])
}
