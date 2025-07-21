package main

import (
   "fmt"
)

func main() {
   var p, q, l, r int
   // Read counts and time bounds
   if _, err := fmt.Scan(&p, &q, &l, &r); err != nil {
       return
   }
   // Read Little Z's schedule intervals
   a := make([]int, p)
   b := make([]int, p)
   for i := 0; i < p; i++ {
       fmt.Scan(&a[i], &b[i])
   }
   // Read Little X's base schedule intervals
   c := make([]int, q)
   d := make([]int, q)
   for j := 0; j < q; j++ {
       fmt.Scan(&c[j], &d[j])
   }
   // Count valid wake-up times
   ans := 0
   for t := l; t <= r; t++ {
       ok := false
       // Check if any interval overlaps when shifted by t
       for i := 0; i < p && !ok; i++ {
           for j := 0; j < q; j++ {
               // Shift X's interval by t: [c[j]+t, d[j]+t]
               if c[j]+t <= b[i] && d[j]+t >= a[i] {
                   ok = true
                   break
               }
           }
       }
       if ok {
           ans++
       }
   }
   fmt.Println(ans)
}
