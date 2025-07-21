package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs64(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

// gcd returns the greatest common divisor of a and b
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
   // use map to store distinct x = -b/k in reduced form
   m := make(map[[2]int64]struct{})
   for i := 0; i < n; i++ {
       var k, b int64
       fmt.Fscan(reader, &k, &b)
       if k == 0 {
           continue
       }
       num := -b
       den := k
       if den < 0 {
           num = -num
           den = -den
       }
       g := gcd(abs64(num), den)
       num /= g
       den /= g
       m[[2]int64{num, den}] = struct{}{}
   }
   // number of kinks is number of distinct breakpoints
   fmt.Println(len(m))
}
