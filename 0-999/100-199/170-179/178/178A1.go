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
   // process and print prefix sums for k = 1 to n-1
   for i := 1; i <= n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       sum += a
       if i < n {
           // print sum of first i elements
           fmt.Println(sum)
       }
   }
}
