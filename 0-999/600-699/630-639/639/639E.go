package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Item struct {
   p, t, l, r int64
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   A := make([]Item, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &A[i].p)
   }
   var T int64
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &A[i].t)
       T += A[i].t
   }
   // sort by ratio t/p ascending (i.e., t[i]*p[j] < t[j]*p[i])
   sort.Slice(A, func(i, j int) bool {
       return A[i].t*A[j].p < A[j].t*A[i].p
   })
   // assign l and r for groups with same ratio
   var suml, sumr int64
   for i := 0; i < n; i++ {
       sumr += A[i].t
       j := i
       for j+1 < n && A[j].t*A[j+1].p == A[j+1].t*A[j].p {
           j++
           sumr += A[j].t
       }
       for k := i; k <= j; k++ {
           A[k].l = suml + A[k].t
           A[k].r = sumr
       }
       suml = sumr
       i = j
   }
   // sort by p ascending
   sort.Slice(A, func(i, j int) bool {
       return A[i].p < A[j].p
   })
   // binary search on x
   lo, hi := 0.0, 1.0
   const iter = 50
   for it := 0; it < iter; it++ {
       mid := (lo + hi) * 0.5
       if check(A, mid, T) {
           lo = mid
       } else {
           hi = mid
       }
   }
   // output
   fmt.Printf("%.12f\n", lo)
}

func check(A []Item, x float64, T int64) bool {
   var max1, max2 float64
   n := len(A)
   i := 0
   for i < n {
       // for group of same p
       // process A[i]
       t1 := float64(A[i].p) * (1.0 - x*float64(A[i].r)/float64(T))
       if max1 > t1+1e-12 {
           return false
       }
       t2 := float64(A[i].p) * (1.0 - x*float64(A[i].l)/float64(T))
       if t2 > max2 {
           max2 = t2
       }
       j := i
       for j+1 < n && A[j+1].p == A[i].p {
           j++
           t1 = float64(A[j].p) * (1.0 - x*float64(A[j].r)/float64(T))
           if max1 > t1+1e-12 {
               return false
           }
           t2 = float64(A[j].p) * (1.0 - x*float64(A[j].l)/float64(T))
           if t2 > max2 {
               max2 = t2
           }
       }
       max1 = max2
       i = j + 1
   }
   return true
}
