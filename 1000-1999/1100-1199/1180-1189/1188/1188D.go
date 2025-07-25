package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int64, n)
   var maxA int64
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
       if a[i] > maxA {
           maxA = a[i]
       }
   }
   // carry (borrow) flags for each a[i]
   carry := make([]bool, n)
   var ans int64
   // process bits 0..60
   for j := 0; j <= 60; j++ {
       var cost0, cost1 int64
       next0 := make([]bool, n)
       next1 := make([]bool, n)
       mask := int64(1) << j
       for i := 0; i < n; i++ {
           b := (a[i] >> j) & 1
           c := int64(0)
           if carry[i] {
               c = 1
           }
           // t=0
           // d0 = (0 - b - c) mod 2
           d0 := (2 - b - c) & 1
           cost0 += d0
           // carry out if b+c+d0 >= 2
           if b+c+d0 >= 2 {
               next0[i] = true
           }
           // t=1
           // d1 = (1 - b - c) mod 2
           d1 := (1 - b - c + 2) & 1
           cost1 += d1
           if b+c+d1 >= 2 {
               next1[i] = true
           }
       }
       if cost0 <= cost1 {
           ans += cost0
           carry = next0
       } else {
           ans += cost1
           carry = next1
       }
   }
   // output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, ans)
}
