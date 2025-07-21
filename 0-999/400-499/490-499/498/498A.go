package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var x1, y1, x2, y2 int64
   if _, err := fmt.Fscan(reader, &x1, &y1); err != nil {
       return
   }
   fmt.Fscan(reader, &x2, &y2)
   var n int
   fmt.Fscan(reader, &n)
   count := 0
   for i := 0; i < n; i++ {
       var a, b, c int64
       fmt.Fscan(reader, &a, &b, &c)
       f1 := a*x1 + b*y1 + c
       f2 := a*x2 + b*y2 + c
       if (f1 < 0 && f2 > 0) || (f1 > 0 && f2 < 0) {
           count++
       }
   }
   fmt.Println(count)
}
