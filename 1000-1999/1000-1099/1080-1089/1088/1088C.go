package main

import "fmt"

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   start := n * 10
   // print number of edges
   fmt.Println(n + 1)
   // first two edges
   fmt.Printf("2 %d 1\n", n)
   fmt.Printf("1 %d %d\n", n, start)
   start--
   // remaining edges
   for i := 0; i < n-1; i++ {
       fmt.Printf("2 %d %d\n", i+1, start)
       start--
   }
}
