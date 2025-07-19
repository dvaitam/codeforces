package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // prefix sum of vowels
   S := make([]int, n+1)
   isVowel := func(c byte) bool {
       switch c {
       case 'I', 'E', 'A', 'O', 'U', 'Y':
           return true
       }
       return false
   }
   for i := 1; i <= n; i++ {
       v := 0
       if isVowel(s[i-1]) {
           v = 1
       }
       S[i] = S[i-1] + v
   }
   // ans for length 1
   prev := float64(S[n])
   ret := prev
   // iterate lengths from 2 to n
   for l := 2; l <= n; l++ {
       // delta = S[n-l+1] - S[l-1]
       delta := float64(S[n-l+1] - S[l-1])
       ans := prev + delta
       ret += ans / float64(l)
       prev = ans
   }
   // output
   fmt.Printf("%.8f\n", ret)
}
