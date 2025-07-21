package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, err := reader.ReadString('\n')
   if err != nil && len(line) == 0 {
       return
   }
   line = strings.TrimSpace(line)
   parts := strings.Fields(line)
   if len(parts) < 2 {
       return
   }
   a := parts[0]
   b := parts[1]
   la, lb := len(a), len(b)
   // maximum length of input numbers
   L := la
   if lb > L {
       L = lb
   }
   // convert to digit arrays aligned from least significant (index 0)
   ai := make([]int, L)
   bi := make([]int, L)
   for i := 0; i < la; i++ {
       ai[i] = int(a[la-1-i] - '0')
   }
   for i := 0; i < lb; i++ {
       bi[i] = int(b[lb-1-i] - '0')
   }
   // determine minimum base and sum-of-digits maximum
   maxDigit := 0
   sumMax := 0
   for i := 0; i < L; i++ {
       if ai[i] > maxDigit {
           maxDigit = ai[i]
       }
       if bi[i] > maxDigit {
           maxDigit = bi[i]
       }
       s := ai[i] + bi[i]
       if s > sumMax {
           sumMax = s
       }
   }
   minBase := maxDigit + 1
   if minBase < 2 {
       minBase = 2
   }
   // initial answer for large bases (no carries)
   ans := L
   // check all bases where carries may occur
   for p := minBase; p <= sumMax; p++ {
       carry := 0
       length := 0
       for i := 0; i < L; i++ {
           s := ai[i] + bi[i] + carry
           carry = s / p
           length++
       }
       // count remaining carry digits
       for carry > 0 {
           carry /= p
           length++
       }
       if length > ans {
           ans = length
       }
   }
   fmt.Println(ans)
}
