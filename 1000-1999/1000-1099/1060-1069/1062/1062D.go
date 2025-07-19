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
   var ans int64
   // maximum multiplier for 2 is n/2
   m := n / 2
   for i := int64(2); i <= m; i++ {
       d := n / i
       // Compute sum = (2 + d) * ((d - 1) / 2) + (d even ? (2 + d)/2 : 0)
       sum := int64(2+d) * ((d - 1) / 2)
       if d%2 == 0 {
           sum += (2 + d) / 2
       }
       ans += sum * 4
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, ans)
}
