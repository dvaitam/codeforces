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
   var t int
   fmt.Fscan(in, &t)
TestLoop:
   for t > 0 {
       t--
       var n, m int
       fmt.Fscan(in, &n, &m)
       a := make([]int, n)
       b := make([]int, n)
       c := make([]int, m)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &a[i])
       }
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &b[i])
       }
       for i := 0; i < m; i++ {
           fmt.Fscan(in, &c[i])
       }
       // prepare mapping of needed paints
       cnt := make(map[int][]int)
       kalV := -1
       lastColor := c[m-1]
       for i := 0; i < n; i++ {
           if a[i] != b[i] {
               cnt[b[i]] = append(cnt[b[i]], i)
           }
           if a[i] == b[i] && a[i] == lastColor {
               kalV = i
           }
       }
       resRev := make([]int, m)
       kal := -1
       // apply paints in reverse order
       for i := 0; i < m; i++ {
           color := c[m-1-i]
           stack := cnt[color]
           if len(stack) == 0 {
               if kal == -1 && kalV == -1 {
                   fmt.Fprintln(out, "No")
                   continue TestLoop
               }
               if kal == -1 {
                   kal = kalV
               }
               resRev[i] = kal
           } else {
               pos := stack[len(stack)-1]
               cnt[color] = stack[:len(stack)-1]
               if kal == -1 {
                   kal = pos
               }
               resRev[i] = pos
           }
       }
       // ensure all differences are painted
       for _, stack := range cnt {
           if len(stack) != 0 {
               fmt.Fprintln(out, "No")
               continue TestLoop
           }
       }
       // success: output result
       fmt.Fprintln(out, "Yes")
       // reverse resRev and print 1-based positions
       for i := m - 1; i >= 0; i-- {
           fmt.Fprint(out, resRev[i]+1)
           if i > 0 {
               fmt.Fprint(out, " ")
           }
       }
       fmt.Fprintln(out)
   }
}
