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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   ans := make([]int, n)
   sl := make([]int, n)
   ct := 0
   y := 0
   ok := true
   var x int
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x)
       if x > i+1 {
           ok = false
       }
       if ok {
           sl[ct] = i
           ct++
           for y < x {
               ct--
               pos := sl[ct]
               ans[pos] = y
               y++
           }
       }
   }
   // fill remaining positions
   if ok {
       fillVal := x + 1
       for ct > 0 {
           ct--
           pos := sl[ct]
           ans[pos] = fillVal
       }
       // output result
       for i, v := range ans {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       writer.WriteByte('\n')
   } else {
       fmt.Fprintln(writer, -1)
   }
}
