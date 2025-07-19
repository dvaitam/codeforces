package main

import "fmt"

func main() {
   var n, S int64
   if _, err := fmt.Scan(&n, &S); err != nil {
       return
   }
   var count int64
   for i := n; i >= 1 && S > 0; i-- {
       cnt := S / i
       if cnt > 0 {
           count += cnt
           S -= cnt * i
       }
   }
   fmt.Println(count)
}
