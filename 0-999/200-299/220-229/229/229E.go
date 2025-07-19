package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   k := make([]int, m)
   type pair struct{ val, grp int }
   var v []pair
   for i := 0; i < m; i++ {
       var ki int
       fmt.Fscan(reader, &ki)
       k[i] = ki
       for j := 0; j < ki; j++ {
           var x int
           fmt.Fscan(reader, &x)
           v = append(v, pair{x, i})
       }
   }
   sort.Slice(v, func(i, j int) bool {
       return v[i].val > v[j].val
   })
   // threshold = nth largest value
   threshold := v[n-1].val
   a := make([]int, m)
   // make a mutable copy of k
   kRem := make([]int, m)
   copy(kRem, k)
   p := 1.0
   done := 0
   // process values greater than threshold
   for i := 0; i < len(v) && v[i].val != threshold; i++ {
       cur := v[i].grp
       denom := kRem[cur]
       if denom < 1 {
           denom = 1
       }
       p *= float64(a[cur]+1) / float64(denom)
       a[cur]++
       kRem[cur]--
       done++
   }
   // collect groups with value == threshold among top items
   var ids []int
   for i := 0; i < len(v) && v[i].val >= threshold; i++ {
       if v[i].val == threshold {
           ids = append(ids, v[i].grp)
       }
   }
   s := len(ids)
   tmp := make([]float64, s)
   for j := 0; j < s; j++ {
       cur := ids[j]
       denom := kRem[cur]
       if denom < 1 {
           denom = 1
       }
       tmp[j] = float64(a[cur]+1) / float64(denom)
   }
   // DP f[i][j]: probability after j processed threshold items, i selected
   f := make([][]float64, s+1)
   for i := range f {
       f[i] = make([]float64, s+1)
   }
   f[0][0] = p
   for j := 0; j < s; j++ {
       for i := 0; i <= j; i++ {
           cur := f[i][j]
           if cur == 0 {
               continue
           }
           // pick this threshold item
           f[i+1][j+1] += cur * tmp[j] * float64(i+1) / float64(j+1)
           // skip this one
           f[i][j+1] += cur * float64(j-i+1) / float64(j+1)
       }
   }
   res := f[n-done][s]
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintf(writer, "%.10f\n", res)
}
