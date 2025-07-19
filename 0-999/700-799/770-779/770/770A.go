package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   var n, k int
   if _, err := fmt.Scan(&n, &k); err != nil {
       return
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i := 0; i < k; i++ {
       w.WriteByte(byte('a' + i))
   }
   for i := k; i < n; i++ {
       if i%2 == k%2 {
           w.WriteByte(byte('a' + k - 2))
       } else {
           w.WriteByte(byte('a' + k - 1))
       }
   }
}
