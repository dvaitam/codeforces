package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   const mod = 1000003
   exp := n - 1
   if exp < 0 {
       exp = 0
   }
   result := 1
   base := 3
   for exp > 0 {
       if exp&1 == 1 {
           result = (result * base) % mod
       }
       base = (base * base) % mod
       exp >>= 1
   }
   fmt.Println(result)
}
