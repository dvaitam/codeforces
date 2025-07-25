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
   fmt.Fscan(reader, &n, &m)
   pockets := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &pockets[i])
   }
   sort.Slice(pockets, func(i, j int) bool { return pockets[i] < pockets[j] })
   var L int64
   if n > 0 {
       L = pockets[n-1] - int64(n)
       if L < 0 {
           L = 0
       }
   }
   // count pockets <= v
   countPocket := func(v int64) int64 {
       i := sort.Search(n, func(i int) bool { return pockets[i] > v })
       return int64(i)
   }
   // h(x): one filter application
   var h func(x int64) int64
   h = func(x int64) int64 {
       c := countPocket(x)
       for {
           nc := countPocket(x + c)
           if nc == c {
               break
           }
           c = nc
       }
       return x + c
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for qi := 0; qi < m; qi++ {
       var x, k int64
       fmt.Fscan(reader, &x, &k)
       res := x
       for t := int64(0); t < k; t++ {
           if res >= L {
               res += (k - t) * int64(n)
               break
           }
           nr := h(res)
           if nr == res {
               break
           }
           res = nr
       }
       fmt.Fprintln(writer, res)
   }
}
