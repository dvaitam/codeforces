package main

import "fmt"

func main() {
   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   // Precompute bit masks
   bin := make([]int, n)
   for i := 0; i < n; i++ {
       bin[i] = 1 << i
   }
   // Adjacency in bitmask form
   w := make([]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Scan(&u, &v)
       u--
       v--
       w[u] |= bin[v]
       w[v] |= bin[u]
   }
   size := 1 << n
   // f[mask][j]: number of paths covering mask ending at j
   f := make([][]int64, size)
   for i := range f {
       f[i] = make([]int64, n)
   }
   // sz[mask]: number of bits in mask
   sz := make([]int, size)
   var ans int64
   // Iterate over all non-empty subsets
   for mask := 1; mask < size; mask++ {
       // find lowest set bit
       pos := 0
       for pos < n && (mask&bin[pos]) == 0 {
           pos++
       }
       // single-bit subset
       if mask == bin[pos] {
           f[mask][pos] = 1
           sz[mask] = 1
           continue
       }
       // size of this subset
       sz[mask] = sz[mask^bin[pos]] + 1
       // extend paths
       for j := pos + 1; j < n; j++ {
           if (mask & bin[j]) == 0 {
               continue
           }
           prev := mask ^ bin[j]
           var cnt int64
           for k := 0; k < n; k++ {
               if (w[j]&bin[k]) != 0 {
                   cnt += f[prev][k]
               }
           }
           f[mask][j] = cnt
           // if closing a cycle and length > 2
           if (w[j]&bin[pos]) != 0 && sz[mask] > 2 {
               ans += cnt
           }
       }
   }
   // each cycle counted twice
   fmt.Println(ans / 2)
}
