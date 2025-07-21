package main

import (
   "fmt"
   "sort"
)

func main() {
   var l, r int64
   if _, err := fmt.Scan(&l, &r); err != nil {
       return
   }
   // generate lucky numbers up to length 10
   var luckies []int64
   var gen func(cur int64, depth int)
   gen = func(cur int64, depth int) {
       if depth > 10 {
           return
       }
       if cur > 0 {
           luckies = append(luckies, cur)
       }
       // append digit 4
       gen(cur*10+4, depth+1)
       // append digit 7
       gen(cur*10+7, depth+1)
   }
   gen(0, 0)
   sort.Slice(luckies, func(i, j int) bool { return luckies[i] < luckies[j] })

   var sum int64
   start := l
   // iterate lucky numbers
   for _, z := range luckies {
       if z < start {
           continue
       }
       if z >= r {
           sum += (r - start + 1) * z
           start = r + 1
           break
       }
       // for range [start, z]
       sum += (z - start + 1) * z
       start = z + 1
       if start > r {
           break
       }
   }
   fmt.Println(sum)
}
