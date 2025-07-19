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
   fmt.Fscan(reader, &n)
   var s string
   fmt.Fscan(reader, &s)

   half := n / 2
   for i := 0; i < half; i++ {
       if s[i] == s[i+half] || (s[i] != s[i+half] && s[i] != 'R') {
           fmt.Fprintln(writer, i+1, i+1+half)
       } else {
           fmt.Fprintln(writer, i+1+half, i+1)
       }
   }
}
