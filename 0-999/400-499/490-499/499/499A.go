package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, x int
   if _, err := fmt.Fscan(reader, &n, &x); err != nil {
       return
   }
   cur := 1
   watched := 0
   for i := 0; i < n; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       // skip gaps using skip button and watch leftover
       gap := l - cur
       if gap > 0 {
           skips := gap / x
           cur += skips * x
           leftover := l - cur
           watched += leftover
           cur += leftover
       }
       // watch the best moment interval
       length := r - l + 1
       watched += length
       cur = r + 1
   }
   fmt.Println(watched)
}
