package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, L, a int
   fmt.Fscan(reader, &n, &L, &a)
   ans := 0
   r := 0
   for i := 0; i < n; i++ {
       var l, t int
       fmt.Fscan(reader, &l)
       ans += (l - r) / a
       fmt.Fscan(reader, &t)
       r = l + t
   }
   ans += (L - r) / a
   fmt.Println(ans)
}
