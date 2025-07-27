package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

const mod = 1000000007
const inv9 = 111111112 // inverse of 9 mod

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, _ := reader.ReadString('\n')
   s = strings.TrimSpace(s)
   n := len(s)
   // digits 1-indexed
   digits := make([]int, n+1)
   for i := 1; i <= n; i++ {
       digits[i] = int(s[i-1] - '0')
   }
   // pow10[i] = 10^i mod
   pow10 := make([]int, n+2)
   pow10[0] = 1
   for i := 1; i <= n+1; i++ {
       pow10[i] = int(int64(pow10[i-1]) * 10 % mod)
   }
   // prefix P[i] = value of s[1..i]
   P := make([]int, n+1)
   for i := 1; i <= n; i++ {
       P[i] = int((int64(P[i-1])*10 + int64(digits[i])) % mod)
   }
   // suffix Suf[i] = value of s[i..n]
   Suf := make([]int, n+2)
   Suf[n+1] = 0
   for i := n; i >= 1; i-- {
       // digits[i]*10^(n-i) + Suf[i+1]
       Suf[i] = int((int64(digits[i])*int64(pow10[n-i]) + int64(Suf[i+1])) % mod)
   }
   var sum1, sum2 int64
   // sum1: contributions from prefixes
   for l := 1; l <= n; l++ {
       // P[l-1] * (sum_{r=l..n} 10^{n-r})
       // sum_{k=0..n-l}10^k = (10^{n-l+1}-1)/9
       x := (int64(pow10[n-l+1]) - 1 + mod) % mod
       x = x * inv9 % mod
       sum1 = (sum1 + int64(P[l-1]) * x) % mod
   }
   // sum2: contributions from suffixes
   for r := 1; r <= n; r++ {
       // suffix starting at r+1, counted r times
       sum2 = (sum2 + int64(r) * int64(Suf[r+1])) % mod
   }
   res := (sum1 + sum2) % mod
   fmt.Println(res)
}
