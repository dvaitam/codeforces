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
       var l, r int64
       fmt.Fscan(reader, &l, &r)
       // output l and 2*l as per problem statement
       fmt.Fprintf(writer, "%d %d\n", l, l*2)
   }
}
