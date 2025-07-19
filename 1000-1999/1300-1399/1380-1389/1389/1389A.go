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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       if r%l == 0 {
           fmt.Fprintln(writer, l, r)
       } else if r%2 == 0 && r/2 >= l {
           fmt.Fprintln(writer, r/2, r)
       } else if (r-1)%2 == 0 && (r-1)/2 >= l {
           fmt.Fprintln(writer, (r-1)/2, r-1)
       } else {
           fmt.Fprintln(writer, -1, -1)
       }
   }
}
