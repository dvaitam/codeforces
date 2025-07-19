package main

import "fmt"

func main() {
   var n int
   var d int64
   if _, err := fmt.Scan(&n, &d); err != nil {
       return
   }
   var last int64 = -1
   var moves int64
   for i := 0; i < n; i++ {
       var cur int64
       fmt.Scan(&cur)
       if last >= cur {
           diff := last - cur
           k := diff/d + 1
           moves += k
           last = cur + k*d
       } else {
           last = cur
       }
   }
   fmt.Println(moves)
}
