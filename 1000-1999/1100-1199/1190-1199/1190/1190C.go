package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}
func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)

   // positions of ones and zeros
   L0, R0 := n, -1
   L1, R1 := n, -1
   for i, ch := range s {
       if ch == '0' {
           L0 = min(L0, i)
           R0 = max(R0, i)
       } else {
           L1 = min(L1, i)
           R1 = max(R1, i)
       }
   }
   // immediate win
   if R0-L0+1 <= k || R1-L1+1 <= k {
       fmt.Fprintln(writer, "tokitsukaze")
       return
   }
   // prev and next for zeros and ones
   prev0 := make([]int, n)
   prev1 := make([]int, n)
   last0, last1 := -1, -1
   for i := 0; i < n; i++ {
       if s[i] == '0' {
           last0 = i
       } else {
           last1 = i
       }
       prev0[i] = last0
       prev1[i] = last1
   }
   next0 := make([]int, n)
   next1 := make([]int, n)
   nxt0, nxt1 := n, n
   for i := n - 1; i >= 0; i-- {
       if s[i] == '0' {
           nxt0 = i
       } else {
           nxt1 = i
       }
       next0[i] = nxt0
       next1[i] = nxt1
   }

   // check for draw possibility
   onceAgain := false
   for i := 0; i + k <= n; i++ {
       l, r := i, i+k-1
       // try flip to 0 or to 1
       // val = 0
       {
           newL0 := min(L0, l)
           newR0 := max(R0, r)
           // ones outside
           nl1, nr1 := n, -1
           if l > 0 {
               p := prev1[l-1]
               if p >= 0 {
                   nl1 = min(nl1, p)
                   nr1 = max(nr1, p)
               }
           }
           if r+1 < n {
               nx := next1[r+1]
               if nx < n {
                   nl1 = min(nl1, nx)
                   nr1 = max(nr1, nx)
               }
           }
           // if after flip span of both > k, Q cannot win immediately
           if newR0 - newL0 + 1 > k && nr1 - nl1 + 1 > k {
               onceAgain = true
               break
           }
       }
       // val = 1
       {
           newL1 := min(L1, l)
           newR1 := max(R1, r)
           // zeros outside
           nl0, nr0 := n, -1
           if l > 0 {
               p := prev0[l-1]
               if p >= 0 {
                   nl0 = min(nl0, p)
                   nr0 = max(nr0, p)
               }
           }
           if r+1 < n {
               nx := next0[r+1]
               if nx < n {
                   nl0 = min(nl0, nx)
                   nr0 = max(nr0, nx)
               }
           }
           if newR1 - newL1 + 1 > k && nr0 - nl0 + 1 > k {
               onceAgain = true
               break
           }
       }
   }
   if onceAgain {
       fmt.Fprintln(writer, "once again")
   } else {
       fmt.Fprintln(writer, "quailty")
   }
}
