package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n, m int
       fmt.Fscan(reader, &n, &m)
       var s string
       fmt.Fscan(reader, &s)
       // cnt[i]: number of wrong tries where p_j == i
       cnt := make([]int, n+2)
       for i := 0; i < m; i++ {
           var p int
           fmt.Fscan(reader, &p)
           if p <= n {
               cnt[p]++
           }
       }
       // build suffix sums: cnt[i] = number of wrong tries with p_j >= i
       for i := n; i >= 1; i-- {
           cnt[i] += cnt[i+1]
       }
       // compute answers for each letter
       ans := make([]int64, 26)
       for i := 1; i <= n; i++ {
           times := int64(cnt[i] + 1)
           ans[s[i-1]-'a'] += times
       }
       // output
       for i, v := range ans {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       writer.WriteByte('\n')
   }
}
