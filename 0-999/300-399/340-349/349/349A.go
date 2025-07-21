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
   var bill int
   count25, count50 := 0, 0
   ok := true
   for i := 0; i < n; i++ {
       if _, err := fmt.Fscan(reader, &bill); err != nil {
           ok = false
           break
       }
       switch bill {
       case 25:
           count25++
       case 50:
           if count25 <= 0 {
               ok = false
           } else {
               count25--
               count50++
           }
       case 100:
           // Prefer giving 50+25 change
           if count50 > 0 && count25 > 0 {
               count50--
               count25--
           } else if count25 >= 3 {
               count25 -= 3
           } else {
               ok = false
           }
       default:
           ok = false
       }
       if !ok {
           // No need to read further
           break
       }
   }
   if ok {
       fmt.Fprintln(writer, "YES")
   } else {
       fmt.Fprintln(writer, "NO")
   }
}
