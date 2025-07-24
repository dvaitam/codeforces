package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, x, y int64
   if _, err := fmt.Fscan(reader, &n, &x, &y); err != nil {
       return
   }
   var ans int64
   for n > 0 {
       if n == 1 {
           ans += x
           break
       }
       if n%2 == 0 {
           half := n / 2
           // cost to delete down to half by single deletions
           costDel := (n - half) * x
           // cost to reverse duplicate
           costDup := y
           if costDup < costDel {
               ans += costDup
           } else {
               ans += costDel
           }
           n = half
       } else {
           // delete one to make even
           ans += x
           n--
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, ans)
}
