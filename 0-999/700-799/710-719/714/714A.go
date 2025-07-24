package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var l1, r1, l2, r2, k int64
   if _, err := fmt.Fscan(reader, &l1, &r1, &l2, &r2, &k); err != nil {
       return
   }
   // Intersection of [l1, r1] and [l2, r2]
   start := l1
   if l2 > start {
       start = l2
   }
   end := r1
   if r2 < end {
       end = r2
   }
   if start > end {
       fmt.Println(0)
       return
   }
   // Total minutes in intersection
   ans := end - start + 1
   // Subtract minute k if it's within the intersection
   if k >= start && k <= end {
       ans--
   }
   if ans < 0 {
       ans = 0
   }
   fmt.Println(ans)
}
