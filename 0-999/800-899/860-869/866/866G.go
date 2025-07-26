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
   if n <= 0 {
       fmt.Println("YES")
       return
   }
   var prev int
   if _, err := fmt.Fscan(reader, &prev); err != nil {
       return
   }
   for i := 1; i < n; i++ {
       var x int
       if _, err := fmt.Fscan(reader, &x); err != nil {
           return
       }
       if x < prev {
           fmt.Println("NO")
           return
       }
       prev = x
   }
   fmt.Println("YES")
}
