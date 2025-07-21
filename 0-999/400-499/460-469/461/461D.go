package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func modPow(a, e int64) int64 {
   res := int64(1)
   a %= MOD
   for e > 0 {
       if e&1 == 1 {
           res = res * a % MOD
       }
       a = a * a % MOD
       e >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   // map diff -> value (0 for 'x', 1 for 'o')
   assign := make(map[int]int)
   ok := true
   for i := 0; i < k; i++ {
       var a, b int
       var c byte
       fmt.Fscan(reader, &a, &b, &c)
       diff := a - b
       if diff < 0 {
           diff = -diff
       }
       val := 0
       if c == 'o' {
           val = 1
       }
       if prev, found := assign[diff]; found {
           if prev != val {
               ok = false
           }
       } else {
           assign[diff] = val
       }
   }
   if !ok {
       fmt.Println(0)
       return
   }
   m := len(assign)
   // number of degrees of freedom is n - m
   exp := int64(n - m)
   if exp < 0 {
       exp = 0
   }
   ans := modPow(2, exp)
   fmt.Println(ans)
}
