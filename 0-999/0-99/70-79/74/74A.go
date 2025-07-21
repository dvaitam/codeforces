package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   bestScore := -100000000
   var bestHandle string
   for i := 0; i < n; i++ {
       var handle string
       var plus, minus int
       var a, b, c, d, e int
       fmt.Fscan(in, &handle, &plus, &minus, &a, &b, &c, &d, &e)
       score := a + b + c + d + e + 100*plus - 50*minus
       if score > bestScore {
           bestScore = score
           bestHandle = handle
       }
   }
   fmt.Println(bestHandle)
}
