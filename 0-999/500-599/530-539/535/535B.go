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

   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   l := len(s)
   idx := (1 << l) - 1
   for i := 0; i < l; i++ {
       if s[i] == '7' {
           idx += 1 << (l - i - 1)
       }
   }
   fmt.Fprintln(writer, idx)
}
