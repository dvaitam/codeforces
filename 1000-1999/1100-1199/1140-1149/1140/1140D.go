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
   if n < 2 {
       ans = 0
   } else {
       // sum_{i=2 to n-1} i*(i+1) = (n-1)*n*(n+1)/3 - 2
       ans = (n-1) * n * (n + 1) / 3
       ans -= 2
   }
   fmt.Fprintln(writer, ans)
}
