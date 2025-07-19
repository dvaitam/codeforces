package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   parent := make([]int, n+2)
   children := make([][]int, n+2)
   parent[1] = 1
   deg := make([]int, n+2)
   for i := 2; i <= n; i++ {
       fmt.Fscan(in, &parent[i])
       p := parent[i]
       children[p] = append(children[p], i)
       deg[p]++
   }
   // Euler tour order (pre-order)
   et := make([]int, 0, n)
   stack := []int{1}
   for len(stack) > 0 {
       v := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       et = append(et, v)
       // push children in reverse to maintain order
       for i := len(children[v]) - 1; i >= 0; i-- {
           stack = append(stack, children[v][i])
       }
   }
   // compute subtree sizes and xor
   sz := make([]int, n+2)
   xr := make([]int, n+2)
   // initialize leaves
   q := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       sz[i] = 1
       xr[i] = a[i]
       if deg[i] == 0 {
           q = append(q, i)
       }
   }
   // prevent root processed early
   deg[1]++
   // process
   head := 0
   for head < len(q) {
       v := q[head]
       head++
       p := parent[v]
       sz[p] += sz[v]
       xr[p] ^= xr[v]
       deg[p]--
       if deg[p] == 0 {
           q = append(q, p)
       }
   }
   // DP
   const maxA = 32
   dp := make([][maxA]bool, n+2)
   // base case: dp[n+1][0] = true
   if n+1 < len(dp) {
       dp[n+1][0] = true
   }
   // fill from n down to 1
   for i := n; i >= 1; i-- {
       v := et[i-1]
       szi := sz[v]
       xri := xr[v]
       for x := 0; x < maxA; x++ {
           ok := dp[i+1][x]
           if !ok && szi%2 == 0 {
               nx := x ^ xri
               if nx < maxA && dp[i+szi][nx] {
                   ok = true
               }
           }
           dp[i][x] = ok
       }
   }
   startX := xr[1]
   if startX >= maxA || !dp[1][startX] {
       fmt.Fprintln(out, -1)
       return
   }
   // reconstruct answer
   ans := make([]int, 0, n)
   i := 1
   x := startX
   for i <= n {
       if dp[i+1][x] {
           i++
       } else {
           v := et[i-1]
           ans = append(ans, v)
           x ^= xr[v]
           i += sz[v]
       }
   }
   // always include root
   ans = append(ans, 1)
   // output
   fmt.Fprintln(out, len(ans))
   for idx, v := range ans {
       if idx > 0 {
           fmt.Fprint(out, " ")
       }
       fmt.Fprint(out, v)
   }
   fmt.Fprintln(out)
}
