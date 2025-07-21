package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int64
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // Convert to binary without leading zeros
   s := strconv.FormatInt(n, 2)
   writer.WriteString(s)
}
