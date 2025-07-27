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
   var ans int64
   if n%2 == 0 {
       k := n / 2
       ans = (k + 1) * (k + 1)
   } else {
       ans = (n + 1) * (n + 3) / 2
   }
   fmt.Fprintln(writer, ans)
}
