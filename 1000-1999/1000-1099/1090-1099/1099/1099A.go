package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, h, w1, h1, w2, h2 int64
   fmt.Fscan(reader, &n, &h)
   fmt.Fscan(reader, &w1, &h1)
   fmt.Fscan(reader, &w2, &h2)

   curr, ans := h, n
   for curr > 0 {
       ans += curr
       if curr == h1 {
           ans -= w1
           if ans < 0 {
               ans = 0
           }
       }
       if curr == h2 {
           ans -= w2
           if ans < 0 {
               ans = 0
           }
       }
       curr--
   }
   if ans < 0 {
       ans = 0
   }
   fmt.Println(ans)
}
