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

   var s string
   var k int
   if _, err := fmt.Fscan(reader, &s, &k); err != nil {
       return
   }
   var n int
   fmt.Fscan(reader, &n)
   // bonuses c[x][y]
   var c [26][26]int
   for i := 0; i < n; i++ {
       var xs, ys string
       var val int
       fmt.Fscan(reader, &xs, &ys, &val)
       x := int(xs[0] - 'a')
       y := int(ys[0] - 'a')
       c[x][y] = val
   }
   m := len(s)
   // dpPrev[j][l]: max euphony up to prev pos with j changes and letter l
   const NEG_INF = -1000000000
   dpPrev := make([][]int, k+1)
   dpCurr := make([][]int, k+1)
   for j := 0; j <= k; j++ {
       dpPrev[j] = make([]int, 26)
       dpCurr[j] = make([]int, 26)
       for l := 0; l < 26; l++ {
           dpPrev[j][l] = NEG_INF
           dpCurr[j][l] = NEG_INF
       }
   }
   // Initialize first position
   orig0 := int(s[0] - 'a')
   for l := 0; l < 26; l++ {
       chg := 0
       if l != orig0 {
           chg = 1
       }
       if chg <= k {
           dpPrev[chg][l] = 0
       }
   }
   // DP over positions
   for pos := 1; pos < m; pos++ {
       // reset dpCurr
       for j := 0; j <= k; j++ {
           for l := 0; l < 26; l++ {
               dpCurr[j][l] = NEG_INF
           }
       }
       curOrig := int(s[pos] - 'a')
       for used := 0; used <= k; used++ {
           for prevL := 0; prevL < 26; prevL++ {
               prevVal := dpPrev[used][prevL]
               if prevVal <= NEG_INF {
                   continue
               }
               for l := 0; l < 26; l++ {
                   nj := used
                   if l != curOrig {
                       nj++
                   }
                   if nj > k {
                       continue
                   }
                   val := prevVal + c[prevL][l]
                   if val > dpCurr[nj][l] {
                       dpCurr[nj][l] = val
                   }
               }
           }
       }
       // swap dpPrev and dpCurr
       dpPrev, dpCurr = dpCurr, dpPrev
   }
   // find answer
   ans := NEG_INF
   for used := 0; used <= k; used++ {
       for l := 0; l < 26; l++ {
           if dpPrev[used][l] > ans {
               ans = dpPrev[used][l]
           }
       }
   }
   // single character string yields 0
   if m <= 1 {
       ans = 0
   }
   fmt.Fprintln(writer, ans)
}
