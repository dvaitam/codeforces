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
   // Read first element to initialize maximum
   var x int
   if n <= 0 {
       fmt.Fprint(writer, 0)
       return
   }
   fmt.Fscan(reader, &x)
   mx := x
   for i := 1; i < n; i++ {
       fmt.Fscan(reader, &x)
       if x > mx {
           mx = x
       }
   }
   fmt.Fprint(writer, mx)
}
