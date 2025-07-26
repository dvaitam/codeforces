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

   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   ratings := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &ratings[i])
   }
   // copy and sort ratings for counting smaller
   sorted := make([]int, n)
   copy(sorted, ratings)
   sort.Ints(sorted)

   // initial answer: number of ratings strictly less than rating[i]
   ans := make([]int, n)
   for i := 0; i < n; i++ {
       // lower_bound of ratings[i]
       cnt := sort.Search(n, func(j int) bool { return sorted[j] >= ratings[i] })
       ans[i] = cnt
   }

   // process friend pairs
   for i := 0; i < k; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--
       v--
       if ratings[u] > ratings[v] {
           ans[u]--
       } else if ratings[v] > ratings[u] {
           ans[v]--
       }
   }

   // output
   for i := 0; i < n; i++ {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, ans[i])
   }
   out.WriteByte('\n')
}
