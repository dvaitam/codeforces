package main

import "fmt"

func main() {
   var a, b, c int64
   if _, err := fmt.Scan(&a, &b, &c); err != nil {
       return
   }
   // Compute a^b mod c using fast exponentiation
   result := int64(1) % c
   base := a % c
   for b > 0 {
       if b&1 == 1 {
           result = (result * base) % c
       }
       base = (base * base) % c
       b >>= 1
   }
   fmt.Println(result)
}
