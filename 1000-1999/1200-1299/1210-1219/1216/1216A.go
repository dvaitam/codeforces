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
   var str string
   fmt.Fscan(reader, &str)
   s := []byte(str)
   cnt := 0
   for i := 1; i < n; i += 2 {
       if s[i] == s[i-1] {
           if s[i-1] == 'a' {
               s[i] = 'b'
           } else {
               s[i] = 'a'
           }
           cnt++
       }
   }
   fmt.Fprintln(writer, cnt)
   writer.Write(s)
}
