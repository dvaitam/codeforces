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
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   n, m := len(s), len(t)
   // prev[j] stores prefix sums of f[i-1][1..j]
   prev := make([]int, m+1)
   var ans int
   for i := 1; i <= n; i++ {
       curr := make([]int, m+1)
       prefix := make([]int, m+1)
       for j := 1; j <= m; j++ {
           if s[i-1] == t[j-1] {
               curr[j] = prev[j-1] + 1
               if curr[j] >= mod {
                   curr[j] -= mod
               }
           }
           prefix[j] = prefix[j-1] + curr[j]
           if prefix[j] >= mod {
               prefix[j] -= mod
           }
       }
       ans += prefix[m]
       if ans >= mod {
           ans -= mod
       }
       // move to next row
       prev = prefix
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, ans)
}
