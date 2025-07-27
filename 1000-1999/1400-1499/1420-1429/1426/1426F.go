package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }

   var dp0, dp1, dp2, dp3 int64
   dp0 = 1 // empty subsequence count (total ways for prefix)

   for i := 0; i < n; i++ {
       c := s[i]
       switch c {
       case 'a':
           dp1 = (dp1 + dp0) % mod
       case 'b':
           dp2 = (dp2 + dp1) % mod
       case 'c':
           dp3 = (dp3 + dp2) % mod
       case '?':
           // for '?', consider as 'c', 'b', 'a'
           dp3 = (3*dp3 + dp2) % mod
           dp2 = (3*dp2 + dp1) % mod
           dp1 = (3*dp1 + dp0) % mod
           dp0 = (dp0 * 3) % mod
       }
   }
   fmt.Println(dp3)
}
