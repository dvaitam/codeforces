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

   var n int64
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var na, nb, ans int64
   for i := int64(0); i < n; i++ {
       var a, b int64
       fmt.Fscan(reader, &a, &b)
       // Determine overlapping range between previous (na,nb) and current [a,b]
       x := na
       if nb > x {
           x = nb
       }
       y := a
       if b < y {
           y = b
       }
       if y-x+1 > 0 {
           ans += y - x + 1
       }
       na, nb = a, b
       // Move the smaller one forward
       if na < nb {
           na++
       } else {
           nb++
       }
   }
   fmt.Fprint(writer, ans)
}
