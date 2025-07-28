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
       fmt.Println(0)
       return
   }
   var x, max int
   // Read first value to initialize max
   if _, err := fmt.Fscan(reader, &x); err != nil {
       return
   }
   max = x
   for i := 1; i < n; i++ {
       if _, err := fmt.Fscan(reader, &x); err != nil {
           return
       }
       if x > max {
           max = x
       }
   }
   fmt.Println(max)
}
