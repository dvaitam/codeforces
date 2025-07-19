package main

import (
   "bufio"
   "fmt"
   "os"
)

// gcd computes the greatest common divisor of a and b.
func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   arr := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i])
   }
   // find maximum element b
   b := 0
   for _, v := range arr {
       if v > b {
           b = v
       }
   }
   // find largest element coprime with b
   second := 0
   for _, v := range arr {
       if v != b && gcd(v, b) == 1 && v > second {
           second = v
       }
   }
   fmt.Fprintln(writer, second, b)
}
