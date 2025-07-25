package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var c, s, t string
   fmt.Fscan(reader, &c)
   fmt.Fscan(reader, &s)
   fmt.Fscan(reader, &t)
   n := len(c)
   m1 := len(s)
   m2 := len(t)
   // build prefix functions
   pi1 := make([]int, m1)
   for i := 1; i < m1; i++ {
       j := pi1[i-1]
       for j > 0 && s[i] != s[j] {
           j = pi1[j-1]
       }
       if s[i] == s[j] {
           j++
       }
       pi1[i] = j
   }
   pi2 := make([]int, m2)
   for i := 1; i < m2; i++ {
       j := pi2[i-1]
       for j > 0 && t[i] != t[j] {
           j = pi2[j-1]
       }
       if t[i] == t[j] {
           j++
       }
       pi2[i] = j
   }
   // build automata
   next1 := make([][26]int, m1)
   next2 := make([][26]int, m2)
   for j := 0; j < m1; j++ {
       for ci := 0; ci < 26; ci++ {
           ch := byte('a' + ci)
           k := j
           for k > 0 && s[k] != ch {
               k = pi1[k-1]
           }
           if s[k] == ch {
               k++
           }
           next1[j][ci] = k
       }
   }
   for j := 0; j < m2; j++ {
       for ci := 0; ci < 26; ci++ {
           ch := byte('a' + ci)
           k := j
           for k > 0 && t[k] != ch {
               k = pi2[k-1]
           }
           if t[k] == ch {
               k++
           }
           next2[j][ci] = k
       }
   }
   // dp[j][k] = best score at current position
   const INF = 1000000000
   dp := make([][]int, m1)
   for i := 0; i < m1; i++ {
       dp[i] = make([]int, m2)
       for j := 0; j < m2; j++ {
           dp[i][j] = -INF
       }
   }
   dp[0][0] = 0
   // iterate positions
   for i := 0; i < n; i++ {
       // prepare next dp
       ndp := make([][]int, m1)
       for x := 0; x < m1; x++ {
           ndp[x] = make([]int, m2)
           for y := 0; y < m2; y++ {
               ndp[x][y] = -INF
           }
       }
       // possible chars
       var chars []int
       if c[i] == '*' {
           chars = make([]int, 26)
           for ci := 0; ci < 26; ci++ {
               chars[ci] = ci
           }
       } else {
           chars = []int{int(c[i] - 'a')}
       }
       for j := 0; j < m1; j++ {
           for k := 0; k < m2; k++ {
               if dp[j][k] <= -INF/2 {
                   continue
               }
               base := dp[j][k]
               for _, ci := range chars {
                   nj := next1[j][ci]
                   add := 0
                   if nj == m1 {
                       add++
                       nj = pi1[m1-1]
                   }
                   nk := next2[k][ci]
                   if nk == m2 {
                       add--
                       nk = pi2[m2-1]
                   }
                   val := base + add
                   if val > ndp[nj][nk] {
                       ndp[nj][nk] = val
                   }
               }
           }
       }
       dp = ndp
   }
   // find answer
   ans := -INF
   for j := 0; j < m1; j++ {
       for k := 0; k < m2; k++ {
           if dp[j][k] > ans {
               ans = dp[j][k]
           }
       }
   }
   fmt.Println(ans)
}
