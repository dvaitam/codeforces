package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
   "strings"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   var q int
   fmt.Fscan(in, &q)
   groupSize := n / k
   // copy of all elements for leftovers
   rest := make([]int, n)
   copy(rest, a)
   // track unique groups
   seen := make(map[string]bool)
   sums := make([]int, 0, k)
   for i := 0; i < q; i++ {
       vals := make([]int, groupSize)
       for j := 0; j < groupSize; j++ {
           var idx int
           fmt.Fscan(in, &idx)
           vals[j] = a[idx-1]
       }
       sort.Ints(vals)
       // build key from sorted values
       ss := make([]string, len(vals))
       for j, v := range vals {
           ss[j] = strconv.Itoa(v)
       }
       key := strings.Join(ss, ",")
       if !seen[key] {
           seen[key] = true
           // remove from rest one by one
           for _, v := range vals {
               for ri, rv := range rest {
                   if rv == v {
                       // remove element at ri
                       rest = append(rest[:ri], rest[ri+1:]...)
                       break
                   }
               }
           }
           sum := 0
           for _, v := range vals {
               sum += v
           }
           sums = append(sums, sum)
       }
   }
   // compute min and max sum
   INF := int(1e9)
   mi := INF
   ma := -INF
   for _, s := range sums {
       mi = min(mi, s)
       ma = max(ma, s)
   }
   // consider possible extra groups from leftovers
   if len(rest) >= groupSize && len(sums) < k {
       sort.Ints(rest)
       // smallest group
       ssmall := 0
       for i := 0; i < groupSize; i++ {
           ssmall += rest[i]
       }
       mi = min(mi, ssmall)
       ma = max(ma, ssmall)
       // largest group
       slarge := 0
       for i := 0; i < groupSize; i++ {
           slarge += rest[len(rest)-1-i]
       }
       mi = min(mi, slarge)
       ma = max(ma, slarge)
   }
   // output averages
   avgMin := float64(mi) / float64(groupSize)
   avgMax := float64(ma) / float64(groupSize)
   // print with sufficient precision
   fmt.Printf("%.15f %.15f\n", avgMin, avgMax)
}
