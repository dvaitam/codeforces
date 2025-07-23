package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   total := 0
   for i := 0; i < n; i++ {
       var x1, y1, x2, y2 int
       fmt.Fscan(in, &x1, &y1, &x2, &y2)
       width := x2 - x1 + 1
       height := y2 - y1 + 1
       total += width * height
   }
   fmt.Println(total)
}
