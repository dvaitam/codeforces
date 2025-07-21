package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var sx, sy, sz int
   for i := 0; i < n; i++ {
       var x, y, z int
       fmt.Fscan(reader, &x, &y, &z)
       sx += x
       sy += y
       sz += z
   }
   if sx == 0 && sy == 0 && sz == 0 {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
