package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var k int64
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   a := make([]int64, n)
   a[0] = 0
   for i := 1; i < n; i++ {
       var p int
       var w int64
       fmt.Fscan(in, &p, &w)
       p--
       a[i] = a[p] ^ w
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   cl := make([]int, n)
   cr := make([]int, n)
   for i := 0; i < n; i++ {
       cl[i] = 0
       cr[i] = n
   }
   nxt := make([]int, n+1)
   var cpref int64
   for z := 61; z >= 0; z-- {
       nxt[n] = n
       for i := n - 1; i >= 0; i-- {
           if (a[i]>>uint(z))&1 == 1 {
               nxt[i] = i
           } else {
               nxt[i] = nxt[i+1]
           }
       }
       var cnt0 int64
       for i := 0; i < n; i++ {
           mid := min(cr[i], nxt[cl[i]])
           if (a[i]>>uint(z))&1 == 1 {
               cnt0 += int64(cr[i] - mid)
           } else {
               cnt0 += int64(mid - cl[i])
           }
       }
       if cnt0 > k {
           for i := 0; i < n; i++ {
               mid := min(cr[i], nxt[cl[i]])
               if (a[i]>>uint(z))&1 == 1 {
                   cl[i] = mid
               } else {
                   cr[i] = mid
               }
           }
       } else {
           k -= cnt0
           for i := 0; i < n; i++ {
               mid := min(cr[i], nxt[cl[i]])
               if (a[i]>>uint(z))&1 == 1 {
                   cr[i] = mid
               } else {
                   cl[i] = mid
               }
           }
           cpref |= 1 << uint(z)
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, cpref)
}
