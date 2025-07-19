package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Pair holds a value and its original index
type Pair struct {
   v   int
   idx int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]Pair, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i].v)
       a[i].idx = i
   }
   sort.Slice(a, func(i, j int) bool {
       if a[i].v != a[j].v {
           return a[i].v < a[j].v
       }
       return a[i].idx < a[j].idx
   })
   maxCount := 1
   retL, retR := 0, 0
   l := 0
   r := 0
   for r < n {
       if a[l].v == a[r].v {
           r++
       } else {
           count := r - l
           if count > maxCount || (count == maxCount && a[r-1].idx-a[l].idx < a[retR].idx-a[retL].idx) {
               maxCount = count
               retL = l
               retR = r - 1
           }
           l = r
       }
   }
   // final segment
   if r-l > maxCount || (r-l == maxCount && a[r-1].idx-a[l].idx < a[retR].idx-a[retL].idx) {
       retL = l
       retR = r - 1
   }
   // output 1-based original indices of first and last
   fmt.Printf("%d %d\n", a[retL].idx+1, a[retR].idx+1)
}
