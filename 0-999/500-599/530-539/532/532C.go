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

   var xp, yp, xv, yv int
   if _, err := fmt.Fscan(reader, &xp, &yp, &xv, &yv); err != nil {
       return
   }
   dP := xp + yp
   dV := xv
   if yv > dV {
       dV = yv
   }
   if dP <= dV {
       fmt.Fprint(writer, "Polycarp")
   } else {
       fmt.Fprint(writer, "Vasiliy")
   }
}
