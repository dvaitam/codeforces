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
   r := n % 4
   var res int
   if r == 1 || r == 2 {
       res = 1
   } else {
       res = 0
   }
   fmt.Fprint(writer, res)
}
