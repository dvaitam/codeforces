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

   var n, p int
   fmt.Fscan(in, &n, &p)
   deg := make([]int, n+1)
   // map for counts of exact pairs
   cnt := make(map[int64]int)
   for i := 0; i < n; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       deg[x]++
       deg[y]++
       if x > y {
           x, y = y, x
       }
       key := int64(x)<<32 | int64(y)
       cnt[key]++
   }
   // sorted degrees
   ds := make([]int, n)
   for i := 1; i <= n; i++ {
       ds[i-1] = deg[i]
   }
   sort.Ints(ds)
   // count initial pairs with sum >= p
   total := 0
   for i := 0; i < n; i++ {
       // find minimal j > i such that ds[i] + ds[j] >= p
       need := p - ds[i]
       j := sort.Search(n, func(j int) bool {
           return ds[j] >= need
       })
       if j < 0 {
           j = 0
       }
       if j <= i {
           j = i + 1
       }
       if j < n {
           total += n - j
       }
   }
   // subtract overcounts where deg[u]+deg[v]-cnt[u,v] < p
   for key, c := range cnt {
       u := int(key >> 32)
       v := int(key & 0xFFFFFFFF)
       if deg[u]+deg[v] >= p && deg[u]+deg[v]-c < p {
           total--
       }
   }
   fmt.Fprintln(out, total)
}
