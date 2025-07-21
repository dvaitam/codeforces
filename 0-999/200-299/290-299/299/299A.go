package main

import (
   "bufio"
   "fmt"
   "os"
)

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   arr := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i])
   }
   if n == 0 {
       fmt.Println(-1)
       return
   }
   G := arr[0]
   for i := 1; i < n; i++ {
       G = gcd(G, arr[i])
   }
   // Check if G appears in the array
   for _, v := range arr {
       if v == G {
           fmt.Println(G)
           return
       }
   }
   fmt.Println(-1)
}
