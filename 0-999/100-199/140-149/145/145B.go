package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var a1, a2, a3, a4 int
   if _, err := fmt.Fscan(reader, &a1, &a2, &a3, &a4); err != nil {
       return
   }
   var start, end byte
   switch {
   case a3 == a4:
       start, end = '4', '4'
   case a3 == a4+1:
       start, end = '4', '7'
   case a4 == a3+1:
       start, end = '7', '4'
   default:
       fmt.Fprintln(writer, "-1")
       return
   }
   var seg4, seg7, k int
   switch {
   case start == '4' && end == '4':
       k = a3
       seg4, seg7 = k+1, k
   case start == '4' && end == '7':
       k = a4
       seg4, seg7 = k+1, k+1
   case start == '7' && end == '4':
       k = a3
       seg4, seg7 = k+1, k+1
   case start == '7' && end == '7':
       k = a3
       seg4, seg7 = k, k+1
   }
   if a1 < seg4 || a2 < seg7 {
       fmt.Fprintln(writer, "-1")
       return
   }
   // determine distribution of digits over segments
   size4First := a1 - (seg4 - 1)
   size7Last := a2 - (seg7 - 1)
   totalLen := a1 + a2
   var b strings.Builder
   b.Grow(totalLen)
   idx4, idx7 := 0, 0
   current := start
   // build segments alternating
   for i := 0; i < seg4+seg7; i++ {
       if current == '4' {
           idx4++
           cnt := 1
           if idx4 == 1 {
               cnt = size4First
           }
           b.WriteString(strings.Repeat("4", cnt))
           current = '7'
       } else {
           idx7++
           cnt := 1
           if idx7 == seg7 {
               cnt = size7Last
           }
           b.WriteString(strings.Repeat("7", cnt))
           current = '4'
       }
   }
   fmt.Fprint(writer, b.String())
}
