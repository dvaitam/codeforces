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

   var l, r uint64
   fmt.Fscan(reader, &l, &r)
   x := l ^ r
   var ans uint64
   for x > 0 {
       ans = (ans << 1) | 1
       x >>= 1
   }
   fmt.Fprint(writer, ans)
}
