package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   // Read permutations
   p := make([][]int, k)
   for i := 0; i < k; i++ {
       p[i] = make([]int, n)
       for j := 0; j < n; j++ {
           fmt.Fscan(reader, &p[i][j])
           p[i][j]-- // zero-based
       }
   }
   // pos[j][v] = position of value v in p[j]
   pos := make([][]int, k)
   for i := 0; i < k; i++ {
       pos[i] = make([]int, n)
       for idx, v := range p[i] {
           pos[i][v] = idx
       }
   }
   // dp over base sequence p[0]
   base := p[0]
   dp := make([]int, n)
   ans := 0
   for i := 0; i < n; i++ {
       dp[i] = 1
       vi := base[i]
       for j := 0; j < i; j++ {
           vj := base[j]
           ok := true
           // check increasing in all other sequences
           for t := 1; t < k; t++ {
               if pos[t][vj] >= pos[t][vi] {
                   ok = false
                   break
               }
           }
           if ok && dp[j]+1 > dp[i] {
               dp[i] = dp[j] + 1
           }
       }
       if dp[i] > ans {
           ans = dp[i]
       }
   }
   // Output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}
