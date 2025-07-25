package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353

// modPow computes a^e mod mod
func modPow(a, e int) int {
   res := 1
   a %= mod
   for e > 0 {
       if e&1 == 1 {
           res = int(int64(res) * int64(a) % mod)
       }
       a = int(int64(a) * int64(a) % mod)
       e >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var w, h int
   if _, err := fmt.Fscan(reader, &w, &h); err != nil {
       return
   }
   // Number of tilings is 2^(w+h) mod mod
   ans := modPow(2, w+h)
   fmt.Println(ans)
}
