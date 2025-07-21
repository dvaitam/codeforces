package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var s string
   var k int
   fmt.Fscan(reader, &s)
   fmt.Fscan(reader, &k)
   n := len(s)
   // cost[i][j]: min changes to make s[i..j] a palindrome
   cost := make([][]int, n)
   for i := range cost {
       cost[i] = make([]int, n)
   }
   // compute costs
   for i := n - 1; i >= 0; i-- {
       for j := i + 1; j < n; j++ {
           c := 0
           if s[i] != s[j] {
               c = 1
           }
           if i+1 <= j-1 {
               cost[i][j] = cost[i+1][j-1] + c
           } else {
               cost[i][j] = c
           }
       }
   }
   // dp[p][i]: min changes for prefix 0..i into p palindromes
   const inf = 1000000000
   dp := make([][]int, k+1)
   prev := make([][]int, k+1)
   for p := 0; p <= k; p++ {
       dp[p] = make([]int, n)
       prev[p] = make([]int, n)
       for i := 0; i < n; i++ {
           dp[p][i] = inf
           prev[p][i] = -1
       }
   }
   // base p=1
   for i := 0; i < n; i++ {
       dp[1][i] = cost[0][i]
       prev[1][i] = -1
   }
   // fill dp for p=2..k
   for p := 2; p <= k; p++ {
       for i := p - 1; i < n; i++ {
           // split at t, prefix 0..t with p-1, segment t+1..i
           for t := p - 2; t < i; t++ {
               cur := dp[p-1][t] + cost[t+1][i]
               if cur < dp[p][i] {
                   dp[p][i] = cur
                   prev[p][i] = t
               }
           }
       }
   }
   // find best p
   best := inf
   bp := 1
   for p := 1; p <= k; p++ {
       if dp[p][n-1] < best {
           best = dp[p][n-1]
           bp = p
       }
   }
   // reconstruct segments
   segments := make([][2]int, 0, bp)
   p := bp
   i := n - 1
   for p > 0 {
       t := prev[p][i]
       l := 0
       if t >= 0 {
           l = t + 1
       }
       segments = append(segments, [2]int{l, i})
       i = t
       p--
   }
   // reverse segments
   for l, r := 0, len(segments)-1; l < r; l, r = l+1, r-1 {
       segments[l], segments[r] = segments[r], segments[l]
   }
   // rebuild string
   runes := []rune(s)
   for _, seg := range segments {
       l, r := seg[0], seg[1]
       for a, b := l, r; a < b; a, b = a+1, b-1 {
           if runes[a] != runes[b] {
               // change right to match left
               runes[b] = runes[a]
           }
       }
   }
   // output
   fmt.Fprintln(writer, best)
   // print segments with plus
   out := make([]rune, 0, n+len(segments)-1)
   for idx, seg := range segments {
       if idx > 0 {
           out = append(out, '+')
       }
       for j := seg[0]; j <= seg[1]; j++ {
           out = append(out, runes[j])
       }
   }
   fmt.Fprintln(writer, string(out))
}
