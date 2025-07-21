package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read target number as string
   var tStr string
   if _, err := fmt.Fscan(reader, &tStr); err != nil {
       return
   }
   // Read available pieces
   var pieces string
   if _, err := fmt.Fscan(reader, &pieces); err != nil {
       return
   }
   // Count available digits
   avail := make([]int, 10)
   for _, ch := range pieces {
       avail[int(ch-'0')]++
   }
   // Combined groups: 6/9 and 2/5
   avail69 := avail[6] + avail[9]
   avail25 := avail[2] + avail[5]
   // Count requirements per instance
   req := make([]int, 10)
   req69 := 0
   req25 := 0
   for _, ch := range tStr {
       switch ch {
       case '6', '9':
           req69++
       case '2', '5':
           req25++
       default:
           req[int(ch-'0')]++
       }
   }
   // Maximum possible instances cannot exceed number of pieces
   maxK := len(pieces)
   // Check fixed digits
   for d := 0; d <= 9; d++ {
       if req[d] > 0 {
           k := avail[d] / req[d]
           if k < maxK {
               maxK = k
           }
       }
   }
   // Check 6/9 group
   if req69 > 0 {
       k := avail69 / req69
       if k < maxK {
           maxK = k
       }
   }
   // Check 2/5 group
   if req25 > 0 {
       k := avail25 / req25
       if k < maxK {
           maxK = k
       }
   }
   // Output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, maxK)
}
