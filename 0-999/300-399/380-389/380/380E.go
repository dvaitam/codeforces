package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const K = 45

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // build powers of two
   p2 := make([]float64, K)
   p2[0] = 1.0
   for i := 1; i < K; i++ {
       p2[i] = p2[i-1] * 2.0
   }
   // prepare sorted values
   type pair struct{ val, idx int }
   v := make([]pair, n)
   for i := 1; i <= n; i++ {
       v[i-1] = pair{a[i], i}
   }
   sort.Slice(v, func(i, j int) bool { return v[i].val < v[j].val })

   // doubly linked list via arrays
   tlft := make([]int, n+2)
   trgt := make([]int, n+2)
   for i := 0; i <= n+1; i++ {
       tlft[i] = i - 1
       trgt[i] = i + 1
   }
   // adjust boundaries
   tlft[0] = 0
   trgt[n] = n + 1

   ans := 0.0
   for _, pr := range v {
       i := pr.idx
       tl := 0.0
       tr := 0.0
       pl := i
       prt := i
       for j := 0; j < K; j++ {
           if pl != 0 {
               dist := float64(pl - tlft[pl])
               tl += dist / p2[j]
               pl = tlft[pl]
           }
           if prt <= n {
               dist := float64(trgt[prt] - prt)
               tr += dist / p2[j]
               prt = trgt[prt]
           }
       }
       // remove i
       left := tlft[i]
       right := trgt[i]
       trgt[left] = right
       tlft[right] = left
       ans += tl * tr * float64(a[i]) * 0.5
   }
   result := ans / (float64(n) * float64(n))
   fmt.Printf("%.12f\n", result)
}
