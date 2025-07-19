package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int64
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   var max1, max2 int64
   for i := int64(0); i < n; i++ {
       var x int64
       fmt.Fscan(reader, &x)
       if x > max1 {
           max2 = max1
           max1 = x
       } else if x > max2 {
           max2 = x
       }
   }
   cnt := m / (k + 1)
   result := max1*(m-cnt) + max2*cnt
   fmt.Println(result)
}
