package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int64
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var result int64
   if n <= 0 {
       result = 1
   } else {
       result = 4 * n
   }
   writer := bufio.NewWriter(os.Stdout)
   fmt.Fprintln(writer, result)
   writer.Flush()
}
