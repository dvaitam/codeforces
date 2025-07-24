package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int64
   fmt.Fscan(reader, &n, &k)
   for {
       mid := int64(1) << (n - 1)
       if k == mid {
           fmt.Println(n)
           return
       } else if k < mid {
           n--
       } else {
           k -= mid
           n--
       }
   }
}
