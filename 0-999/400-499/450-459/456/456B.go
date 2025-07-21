package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n string
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // If n == "0", 1^0+2^0+3^0+4^0 = 4
   if n == "0" {
       fmt.Println(4)
       return
   }
   // compute n mod 4 and n mod 2
   mod4, mod2 := 0, 0
   for i := 0; i < len(n); i++ {
       d := int(n[i] - '0')
       mod4 = (mod4*10 + d) % 4
       mod2 = (mod2*10 + d) % 2
   }
   // exponent for 2 and 3: if mod4==0 use 4
   e4 := mod4
   if e4 == 0 {
       e4 = 4
   }
   e2 := e4
   e3 := e4
   // exponent for 4: cycle length 2
   e_4 := mod2
   if e_4 == 0 {
       e_4 = 2
   }
   // compute powers mod 5
   val1 := 1
   val2 := powMod(2, e2, 5)
   val3 := powMod(3, e3, 5)
   val4 := powMod(4, e_4, 5)
   res := (val1 + val2 + val3 + val4) % 5
   fmt.Println(res)
}

// powMod computes a^e mod m, where e is small
func powMod(a, e, m int) int {
   res := 1
   a %= m
   for i := 0; i < e; i++ {
       res = (res * a) % m
   }
   return res
}
