package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil && len(s) == 0 {
       return
   }
   // Remove trailing newline if present
   if s[len(s)-1] == '\n' {
       s = s[:len(s)-1]
   }
   n := len(s)
   // f[i]: palindrome degree of prefix s[0..i]
   f := make([]int, n)
   var ans int64
   const P uint64 = 131
   var h1, h2, pPow uint64 = 0, 0, 1
   for i := 0; i < n; i++ {
       c := uint64(s[i])
       // forward-reverse hash and forward hash
       h1 = h1*P + c
       h2 = h2 + c*pPow
       pPow *= P
       if h1 == h2 {
           // prefix [0..i] is palindrome
           halfIdx := ((i + 1) >> 1) - 1
           if halfIdx >= 0 {
               f[i] = f[halfIdx] + 1
           } else {
               f[i] = 1
           }
       } else {
           f[i] = 0
       }
       ans += int64(f[i])
   }
   // output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}
