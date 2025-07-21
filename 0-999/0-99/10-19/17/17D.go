package main

import (
   "bufio"
   "fmt"
   "os"
)

// decStringMinusOne subtracts one from a decimal string s (s >= "1").
// Returns the result without leading zeros (except "0").
func decStringMinusOne(s string) string {
   b := []byte(s)
   i := len(b) - 1
   for i >= 0 {
       if b[i] > '0' {
           b[i]--
           break
       }
       b[i] = '9'
       i--
   }
   // trim leading zeros
   j := 0
   for j < len(b)-1 && b[j] == '0' {
       j++
   }
   return string(b[j:])
}

// modPowInt computes x^e mod m, where e is small (e >= 0).
func modPowInt(x int64, e int, m int64) int64 {
   res := int64(1)
   base := x % m
   for e > 0 {
       if e&1 != 0 {
           res = (res * base) % m
       }
       base = (base * base) % m
       e >>= 1
   }
   return res
}

// modPowDecimalExp computes base^exp (mod m), where exp is a decimal string.
func modPowDecimalExp(base int64, exp string, m int64) int64 {
   result := int64(1)
   for i := 0; i < len(exp); i++ {
       d := int(exp[i] - '0')
       // result = result^10 mod m
       result = modPowInt(result, 10, m)
       // multiply by base^d mod m
       if d != 0 {
           result = (result * modPowInt(base, d, m)) % m
       }
   }
   return result
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var bStr, nStr string
   var c int64
   if _, err := fmt.Fscan(reader, &bStr, &nStr, &c); err != nil {
       return
   }
   // compute b mod c
   var bMod int64
   for i := 0; i < len(bStr); i++ {
       bMod = (bMod*10 + int64(bStr[i]-'0')) % c
   }
   // compute exponent = n - 1
   expStr := decStringMinusOne(nStr)
   // compute p = b^exp mod c
   p := modPowDecimalExp(bMod, expStr, c)
   // total numbers mod c = (b-1)*b^{n-1} mod c
   bMinus1Mod := (bMod - 1 + c) % c
   totalMod := (bMinus1Mod * p) % c
   var ans int64
   if totalMod == 0 {
       ans = c
   } else {
       ans = totalMod
   }
   fmt.Println(ans)
}
