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
   const mod = 1000000007
   result := int64(1)
   for i := 2; i <= n; i++ {
       result = (result * int64(i)) % mod
   }
   fmt.Fprintln(writer, result)
}
