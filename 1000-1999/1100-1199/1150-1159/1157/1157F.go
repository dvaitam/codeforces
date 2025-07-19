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

   var n int
   fmt.Fscan(reader, &n)
   const MAX = 200005
   cnt := make([]int, MAX)
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if x >= 0 && x < MAX {
           cnt[x]++
       }
   }
   sumcnt := make([]int, MAX)
   for i := 1; i < MAX; i++ {
       sumcnt[i] = sumcnt[i-1] + cnt[i]
   }
   ans := -1000000000 // sufficiently small initial value
   ansL, ansR := 0, 0
   r := 0
   for l := 1; l < MAX; l++ {
       if cnt[l] == 0 {
           continue
       }
       if r < l {
           r = l
       }
       // expand r while next is valid
       for r+1 < MAX && cnt[r+1] >= 1 && (cnt[r] >= 2 || l == r) {
           r++
       }
       tans := sumcnt[r] - sumcnt[l-1]
       if tans > ans {
           ans = tans
           ansL = l
           ansR = r
       }
   }
   // output
   fmt.Fprintln(writer, ans)
   // first pass: one each from ansL to ansR
   for i := ansL; i <= ansR; i++ {
       if cnt[i] > 0 {
           fmt.Fprint(writer, i, " ")
           cnt[i]--
       }
   }
   // second pass: remaining in reverse
   for i := ansR; i >= ansL; i-- {
       for cnt[i] > 0 {
           fmt.Fprint(writer, i, " ")
           cnt[i]--
       }
   }
   fmt.Fprintln(writer)
}
