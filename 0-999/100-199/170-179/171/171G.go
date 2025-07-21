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

   var a1, a2, a3 int
   if _, err := fmt.Fscan(reader, &a1, &a2, &a3); err != nil {
       return
   }
   // Find the median (the middle value) of three integers
   var med int
   if (a1 <= a2 && a2 <= a3) || (a3 <= a2 && a2 <= a1) {
       med = a2
   } else if (a2 <= a1 && a1 <= a3) || (a3 <= a1 && a1 <= a2) {
       med = a1
   } else {
       med = a3
   }
   fmt.Fprintln(writer, med)
}
