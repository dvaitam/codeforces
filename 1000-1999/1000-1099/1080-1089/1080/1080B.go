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

   var q int
   if _, err := fmt.Fscan(reader, &q); err != nil {
       return
   }
   for i := 0; i < q; i++ {
       var l, r int64
       fmt.Fscan(reader, &l, &r)
       // compute a = ceil((l-1)/2.0) == l/2
       a := l / 2
       a1 := a
       if (l-1)%2 != 0 {
           a1 = -a
       }
       // compute b = ceil(r/2.0) == (r+1)/2
       b := (r + 1) / 2
       b1 := b
       if r%2 != 0 {
           b1 = -b
       }
       res := b1 - a1
       fmt.Fprintln(writer, res)
   }
}
