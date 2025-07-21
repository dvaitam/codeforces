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
   sum := 0
   for i := 0; i < n; i++ {
       var x int
       if _, err := fmt.Fscan(reader, &x); err != nil {
           return
       }
       sum += x
   }
   fmt.Println(sum)
}
