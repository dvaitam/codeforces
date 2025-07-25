package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   oddA, evenA := 0, 0
   for i := 0; i < n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       if a&1 == 1 {
           oddA++
       } else {
           evenA++
       }
   }
   oddB, evenB := 0, 0
   for j := 0; j < m; j++ {
       var b int
       fmt.Fscan(reader, &b)
       if b&1 == 1 {
           oddB++
       } else {
           evenB++
       }
   }
   // A chest and key sum to odd if one is odd and the other is even
   // Maximum matches: match odd chests with even keys, and even chests with odd keys
   res := min(oddA, evenB) + min(evenA, oddB)
   fmt.Fprintln(writer, res)
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}
