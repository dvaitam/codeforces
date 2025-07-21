package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // build set for quick lookup
   denomSet := make(map[int]struct{}, n)
   for _, v := range a {
       denomSet[v] = struct{}{}
   }
   var q int
   fmt.Fscan(reader, &q)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for qi := 0; qi < q; qi++ {
       var x int
       fmt.Fscan(reader, &x)
       ans := k + 1
       // single denomination
       for _, di := range a {
           if di > 0 && x%di == 0 {
               b := x / di
               if b <= k && b < ans {
                   ans = b
               }
           }
       }
       // two denominations
       for _, d1 := range a {
           // use c1 bills of d1
           // c1 from 1 up to min(k-1, x/d1)
           maxC1 := x / d1
           if maxC1 > k-1 {
               maxC1 = k - 1
           }
           for c1 := 1; c1 <= maxC1; c1++ {
               if c1 >= ans {
                   break
               }
               rem := x - c1*d1
               // need c2 bills of some d2
               maxC2 := k - c1
               // c2 >=1
               for c2 := 1; c2 <= maxC2; c2++ {
                   if c1+c2 >= ans {
                       break
                   }
                   if rem%c2 != 0 {
                       continue
                   }
                   d2 := rem / c2
                   if _, ok := denomSet[d2]; ok {
                       ans = c1 + c2
                       break
                   }
               }
           }
       }
       if ans > k {
           fmt.Fprintln(out, -1)
       } else {
           fmt.Fprintln(out, ans)
       }
   }
}
