package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func modPow(a, e int64) int64 {
   res := int64(1)
   a %= mod
   for e > 0 {
       if e&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       e >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   // Count occurrences of A, C, G, T
   cnt := make(map[byte]int)
   for i := 0; i < n; i++ {
       cnt[s[i]]++
   }
   // Find maximum count and number of chars achieving it
   maxCnt := 0
   for _, c := range []byte{'A', 'C', 'G', 'T'} {
       if cnt[c] > maxCnt {
           maxCnt = cnt[c]
       }
   }
   chars := 0
   for _, c := range []byte{'A', 'C', 'G', 'T'} {
       if cnt[c] == maxCnt {
           chars++
       }
   }
   // Number of strings is chars^n mod mod
   ans := modPow(int64(chars), int64(n))
   fmt.Println(ans)
}
