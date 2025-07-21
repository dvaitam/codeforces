package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var x, t, a, b, da, db int
   if _, err := fmt.Fscan(reader, &x, &t, &a, &b, &da, &db); err != nil {
       return
   }

   // Precompute possible scores for each problem
   set1 := make(map[int]bool)
   set2 := make(map[int]bool)
   for i := 0; i < t; i++ {
       s1 := a - i*da
       s2 := b - i*db
       set1[s1] = true
       set2[s2] = true
   }

   possible := false
   // Zero problems
   if x == 0 {
       possible = true
   }
   // One problem
   if !possible {
       if set1[x] || set2[x] {
           possible = true
       }
   }
   // Two problems
   if !possible {
       for s1 := range set1 {
           if set2[x-s1] {
               possible = true
               break
           }
       }
   }

   if possible {
       fmt.Fprint(writer, "YES")
   } else {
       fmt.Fprint(writer, "NO")
   }
}
