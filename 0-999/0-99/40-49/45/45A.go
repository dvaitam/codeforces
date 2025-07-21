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

   var s string
   var k int
   fmt.Fscan(reader, &s)
   fmt.Fscan(reader, &k)
   months := []string{
       "January", "February", "March", "April", "May", "June",
       "July", "August", "September", "October", "November", "December",
   }
   idx := 0
   for i, m := range months {
       if m == s {
           idx = i
           break
       }
   }
   res := months[(idx+k)%12]
   fmt.Fprint(writer, res)
}
