package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   sBytes, err := reader.ReadBytes('\n')
   if err != nil && len(sBytes) == 0 {
       return
   }
   // Remove trailing newline if present
   if sBytes[len(sBytes)-1] == '\n' {
       sBytes = sBytes[:len(sBytes)-1]
   }
   s := sBytes
   n := len(s)
   // Precompute arrays
   L := make([]int, n)      // consecutive local-part chars ending at i
   PL := make([]int, n)     // prefix sum of letters
   isLetter := func(c byte) bool { return c >= 'a' && c <= 'z' }
   isDigit := func(c byte) bool { return c >= '0' && c <= '9' }
   isLocal := func(c byte) bool { return isLetter(c) || isDigit(c) || c == '_' }
   isDomain := func(c byte) bool { return isLetter(c) || isDigit(c) }
   // Build L and PL
   for i := 0; i < n; i++ {
       if isLocal(s[i]) {
           if i > 0 {
               L[i] = L[i-1] + 1
           } else {
               L[i] = 1
           }
       }
       if isLetter(s[i]) {
           if i > 0 {
               PL[i] = PL[i-1] + 1
           } else {
               PL[i] = 1
           }
       } else if i > 0 {
           PL[i] = PL[i-1]
       }
   }
   // R1: domain chars
   R1 := make([]int, n)
   for i := n - 1; i >= 0; i-- {
       if isDomain(s[i]) {
           if i+1 < n {
               R1[i] = R1[i+1] + 1
           } else {
               R1[i] = 1
           }
       }
   }
   // R2: tld letters
   R2 := make([]int, n)
   for i := n - 1; i >= 0; i-- {
       if isLetter(s[i]) {
           if i+1 < n {
               R2[i] = R2[i+1] + 1
           } else {
               R2[i] = 1
           }
       }
   }
   // DotValue and prefix sum
   DotValue := make([]int64, n)
   for i := 0; i < n-1; i++ {
       if s[i] == '.' {
           DotValue[i] = int64(R2[i+1])
       }
   }
   PDOV := make([]int64, n)
   for i := 0; i < n; i++ {
       if i > 0 {
           PDOV[i] = PDOV[i-1] + DotValue[i]
       } else {
           PDOV[i] = DotValue[i]
       }
   }
   // Iterate over '@'
   var ans int64
   for i := 0; i < n; i++ {
       if s[i] != '@' {
           continue
       }
       if i == 0 || i+1 >= n {
           continue
       }
       leftLen := L[i-1]
       if leftLen < 1 {
           continue
       }
       // count letters in [i-leftLen, i-1]
       start := i - leftLen
       end := i - 1
       var letters int
       if start > 0 {
           letters = PL[end] - PL[start-1]
       } else {
           letters = PL[end]
       }
       if letters == 0 {
           continue
       }
       domMax := 0
       if i+1 < n {
           domMax = R1[i+1]
       }
       if domMax < 1 {
           continue
       }
       dLow := i + 2
       if dLow >= n {
           continue
       }
       dHigh := i + 1 + domMax
       if dHigh >= n {
           dHigh = n - 1
       }
       var rightSum int64
       if dLow > 0 {
           rightSum = PDOV[dHigh] - PDOV[dLow-1]
       } else {
           rightSum = PDOV[dHigh]
       }
       if rightSum <= 0 {
           continue
       }
       ans += int64(letters) * rightSum
   }
   fmt.Println(ans)
}
