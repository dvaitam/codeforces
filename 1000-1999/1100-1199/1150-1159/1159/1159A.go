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
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   cur, minCur := 0, 0
   for _, ch := range s {
       if ch == '+' {
           cur++
       } else {
           cur--
       }
       if cur < minCur {
           minCur = cur
       }
   }
   // initial stones needed
   initStones := -minCur
   // final stones
   result := initStones + cur
   if result < 0 {
       result = 0
   }
   fmt.Println(result)
}
