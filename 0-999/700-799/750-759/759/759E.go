package main

import "fmt"

const mod = 1000000007

func main() {
   var x, y int
   if _, err := fmt.Scan(&x, &y); err != nil {
       return
   }
   // compute x^y mod mod using fast exponentiation
   base := x % mod
   if base < 0 {
       base += mod
   }
   res := 1
   exp := y
   for exp > 0 {
       if exp&1 == 1 {
           res = int((int64(res) * int64(base)) % mod)
       }
       // square base
       base = int((int64(base) * int64(base)) % mod)
       exp >>= 1
   }
   fmt.Println(res)
}
