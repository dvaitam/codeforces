package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var s string
   fmt.Fscan(reader, &n)
   fmt.Fscan(reader, &s)
   m := len(s)
   totalLen := 2 * n
   if m > totalLen {
       fmt.Fprintln(writer, 0)
       return
   }
   // build prefix function
   pi := make([]int, m)
   for i := 1; i < m; i++ {
       j := pi[i-1]
       for j > 0 && s[j] != s[i] {
           j = pi[j-1]
       }
       if s[j] == s[i] {
           j++
       }
       pi[i] = j
   }
   // build automaton: states 0..m, m is absorbing (matched)
   // nxt[state][0] for '(', nxt[state][1] for ')'
   nxt := make([][2]int, m+1)
   for st := 0; st <= m; st++ {
       for c := 0; c < 2; c++ {
           if st == m {
               nxt[st][c] = m
               continue
           }
           var ch byte
           if c == 0 {
               ch = '('
           } else {
               ch = ')'
           }
           j := st
           for j > 0 && s[j] != ch {
               j = pi[j-1]
           }
           if s[j] == ch {
               j++
           }
           nxt[st][c] = j
       }
   }
   // dp[b][st] = ways
   dp := make([][]int, n+1)
   for i := range dp {
       dp[i] = make([]int, m+1)
   }
   dp[0][0] = 1
   // iterate positions
   for step := 0; step < totalLen; step++ {
       // next dp
       nxtdp := make([][]int, n+1)
       for i := range nxtdp {
           nxtdp[i] = make([]int, m+1)
       }
       for b := 0; b <= n; b++ {
           for st := 0; st <= m; st++ {
               v := dp[b][st]
               if v == 0 {
                   continue
               }
               // add '('
               if b+1 <= n {
                   ns := nxt[st][0]
                   nxtdp[b+1][ns] = (nxtdp[b+1][ns] + v) % MOD
               }
               // add ')'
               if b > 0 {
                   ns := nxt[st][1]
                   nxtdp[b-1][ns] = (nxtdp[b-1][ns] + v) % MOD
               }
           }
       }
       dp = nxtdp
   }
   // answer: balance 0 and matched state m
   ans := dp[0][m]
   fmt.Fprintln(writer, ans)
}
