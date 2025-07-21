package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   var sx, sy, sx2, sy2 int64
   for i := 0; i < n; i++ {
       var x, y int64
       if _, err := fmt.Fscan(in, &x, &y); err != nil {
           return
       }
       sx += x
       sy += y
       sx2 += x * x
       sy2 += y * y
   }
   N := int64(n)
   ansX := N*sx2 - sx*sx
   ansY := N*sy2 - sy*sy
   ans := ansX + ansY
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprint(out, ans)
}
