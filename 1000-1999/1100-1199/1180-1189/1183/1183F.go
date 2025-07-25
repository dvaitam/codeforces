package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var q int
   if _, err := fmt.Fscan(in, &q); err != nil {
       return
   }
   for q > 0 {
       q--
       var n int
       fmt.Fscan(in, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &a[i])
       }
       // sort descending and dedupe
       sort.Ints(a)
       m := 0
       for i := 0; i < n; i++ {
           if m == 0 || a[i] != a[m-1] {
               a[m] = a[i]
               m++
           }
       }
       a = a[:m]
       // descending order
       for i := 0; i < m/2; i++ {
           a[i], a[m-1-i] = a[m-1-i], a[i]
       }
       if m == 0 {
           fmt.Fprintln(out, 0)
           continue
       }
       best1 := a[0]
       best2 := 0
       // best pair
       for i := 0; i < m-1; i++ {
           if a[i]+a[i+1] <= best2 {
               break
           }
           ai := a[i]
           for j := i + 1; j < m; j++ {
               sum2 := ai + a[j]
               if sum2 <= best2 {
                   break
               }
               if ai%a[j] != 0 {
                   best2 = sum2
                   break
               }
           }
       }
       best3 := best2
       // best triplet
       for i := 0; i < m-2; i++ {
           // pruning: max possible sum with next two
           if a[i]+a[i+1]+a[i+2] <= best3 {
               break
           }
           ai := a[i]
           for j := i + 1; j < m-1; j++ {
               // pruning for this pair
               if ai+a[j]+a[j+1] <= best3 {
                   break
               }
               if ai%a[j] == 0 {
                   continue
               }
               for k := j + 1; k < m; k++ {
                   sum3 := ai + a[j] + a[k]
                   if sum3 <= best3 {
                       break
                   }
                   if ai%a[k] != 0 && a[j]%a[k] != 0 {
                       best3 = sum3
                       break
                   }
               }
           }
       }
       res := best1
       if best2 > res {
           res = best2
       }
       if best3 > res {
           res = best3
       }
       fmt.Fprintln(out, res)
   }
}
