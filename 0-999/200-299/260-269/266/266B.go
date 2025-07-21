package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, t int
   if _, err := fmt.Fscan(reader, &n, &t); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)
   // simulate t seconds
   b := []byte(s)
   for step := 0; step < t; step++ {
       for i := 0; i+1 < n; {
           if b[i] == 'B' && b[i+1] == 'G' {
               b[i], b[i+1] = b[i+1], b[i]
               i += 2
           } else {
               i++
           }
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, string(b))
}
