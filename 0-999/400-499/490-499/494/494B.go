package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s, t string
   fmt.Fscan(reader, &s)
   fmt.Fscan(reader, &t)
   n, m := len(s), len(t)
   // Build prefix function for pattern t
   pi := make([]int, m)
   for i := 1; i < m; i++ {
       j := pi[i-1]
       for j > 0 && t[i] != t[j] {
           j = pi[j-1]
       }
       if t[i] == t[j] {
           j++
       }
       pi[i] = j
   }
   // Find occurrences of t in s
   occ := make([]bool, n+1)
   j := 0
   for i := 0; i < n; i++ {
       for j > 0 && s[i] != t[j] {
           j = pi[j-1]
       }
       if s[i] == t[j] {
           j++
       }
       if j == m {
           // match ends at i (0-indexed), record at position i+1 (1-indexed)
           occ[i+1] = true
           j = pi[j-1]
       }
   }
   // dp arrays
   d := make([]int64, n+1)
   sdp := make([]int64, n+1)
   sdp[0] = 1
   var sumLast int64
   // Compute DP
   for i := 1; i <= n; i++ {
       if occ[i] {
           pos := i - m + 1
           if pos-1 >= 0 {
               sumLast = (sumLast + sdp[pos-1]) % mod
           }
       }
       d[i] = sumLast
       sdp[i] = (sdp[i-1] + d[i]) % mod
   }
   // subtract empty selection
   ans := (sdp[n] - 1 + mod) % mod
   fmt.Println(ans)
}
