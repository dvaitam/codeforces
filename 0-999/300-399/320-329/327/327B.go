package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   start := 10000000 - n + 1
   for i := 0; i < n; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, start+i)
   }
   writer.WriteByte('\n')
}
