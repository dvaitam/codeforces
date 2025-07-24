package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, c int
   if _, err := fmt.Fscan(reader, &n, &c); err != nil {
       return
   }
   // Read first word
   var l int
   fmt.Fscan(reader, &l)
   prev := make([]int, l)
   for i := 0; i < l; i++ {
       fmt.Fscan(reader, &prev[i])
   }
   // Initialize allowed shift interval [L, R]
   L, R := 0, c-1
   ok := true

   // Process remaining words
   for i := 1; i < n; i++ {
       fmt.Fscan(reader, &l)
       curr := make([]int, l)
       for j := 0; j < l; j++ {
           fmt.Fscan(reader, &curr[j])
       }
       // Compare prev and curr words
       minl := len(prev)
       if l < minl {
           minl = l
       }
       k := 0
       for k < minl && prev[k] == curr[k] {
           k++
       }
       if k == minl {
           // One is prefix
           if len(prev) > len(curr) {
               ok = false
               break
           }
       } else {
           a, b := prev[k], curr[k]
           if a > b {
               // Determine shifts s where (a+s)%c < (b+s)%c
               Li := c - a + 1
               Ri := c - b
               if Li > R || Ri < L {
                   ok = false
                   break
               }
               if Li > L {
                   L = Li
               }
               if Ri < R {
                   R = Ri
               }
           }
       }
       prev = curr
   }
   if !ok || L > R {
       fmt.Fprintln(writer, -1)
       return
   }
   fmt.Fprintln(writer, L)
}
