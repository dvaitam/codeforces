package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int64
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   var sum int64
   for i := int64(0); i < n; i++ {
       var l, r int64
       fmt.Fscan(reader, &l, &r)
       sum += r - l + 1
   }
   rem := sum % k
   if rem == 0 {
       fmt.Println(0)
   } else {
       fmt.Println(k - rem)
   }
}
