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
   for tc := 0; tc < t; tc++ {
       var n, s, r int
       fmt.Fscan(reader, &n, &s, &r)
       // first element
       first := s - r
       // distribute r over n-1 elements as equally as possible
       base := r / (n - 1)
       rem := r % (n - 1)
       // output
       fmt.Fprint(writer, first)
       for i := 0; i < n-1; i++ {
           val := base
           if i < rem {
               val++
           }
           fmt.Fprint(writer, " ", val)
       }
       fmt.Fprint(writer, '\n')
   }
}
