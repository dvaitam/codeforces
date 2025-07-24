package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var c, n int
   if _, err := fmt.Fscan(reader, &c); err != nil {
       return
   }
   fmt.Fscan(reader, &n)
   counts := make(map[int]int)
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       counts[x]++
   }
   // get unique values sorted descending
   vals := make([]int, 0, len(counts))
   for v := range counts {
       vals = append(vals, v)
   }
   sort.Sort(sort.Reverse(sort.IntSlice(vals)))
   rem := int64(c)
   var best int = 0
   for _, v := range vals {
       cnt := counts[v]
       need := int64(cnt+1) * int64(v)
       if rem >= need {
           if best == 0 || v < best {
               best = v
           }
       }
       // greedy pick
       take := int(rem) / v
       if take > cnt {
           take = cnt
       }
       rem -= int64(take * v)
   }
   if best > 0 {
       fmt.Println(best)
   } else {
       fmt.Println("Greed is good")
   }
}
