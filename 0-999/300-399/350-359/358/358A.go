package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   xs := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &xs[i])
   }
   // Build segments between consecutive points
   segs := make([][2]int, 0, n-1)
   for i := 0; i+1 < n; i++ {
       a, b := xs[i], xs[i+1]
       if a > b {
           a, b = b, a
       }
       segs = append(segs, [2]int{a, b})
   }
   // Check for intersecting semicircles: overlapping but not nested
   for i := 0; i < len(segs); i++ {
       a1, b1 := segs[i][0], segs[i][1]
       for j := i + 1; j < len(segs); j++ {
           a2, b2 := segs[j][0], segs[j][1]
           if (a1 < a2 && a2 < b1 && b1 < b2) || (a2 < a1 && a1 < b2 && b2 < b1) {
               fmt.Println("yes")
               return
           }
       }
   }
   fmt.Println("no")
}
