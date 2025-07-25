package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var k int64
   if _, err := fmt.Fscan(reader, &k); err != nil {
       return
   }
   var digitLen int64 = 1
   var count int64 = 9
   var start int64 = 1
   // Find the range where the k-th digit is
   for {
       digits := count * digitLen
       if k <= digits {
           break
       }
       k -= digits
       digitLen++
       count *= 10
       start *= 10
   }
   // Determine the exact number and digit
   index := (k - 1) / digitLen
   num := start + index
   offset := (k - 1) % digitLen
   s := fmt.Sprintf("%d", num)
   // Output the specific digit
   fmt.Println(string(s[offset]))
}
