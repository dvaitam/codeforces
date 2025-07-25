package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, b, c int64
   if _, err := fmt.Fscan(reader, &a, &b, &c); err != nil {
       return
   }
   x := a + c
   y := b + c
   // Option1: start with 'a'
   opt1 := min(2*x, 2*y+1)
   // Option2: start with 'b'
   opt2 := min(2*y, 2*x+1)
   ans := max(opt1, opt2)
   fmt.Println(ans)
}
