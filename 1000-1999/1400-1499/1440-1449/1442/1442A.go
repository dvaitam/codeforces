package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       // prefix minima and suffix minima
       L := make([]int, n)
       R := make([]int, n)
       for i := 0; i < n; i++ {
           if i == 0 || a[i] < L[i-1] {
               L[i] = a[i]
           } else {
               L[i] = L[i-1]
           }
       }
       for i := n - 1; i >= 0; i-- {
           if i == n-1 || a[i] < R[i+1] {
               R[i] = a[i]
           } else {
               R[i] = R[i+1]
           }
       }
       ok := true
       for i := 0; i < n; i++ {
           // maximum reductions at a[i] via prefix k=i+1 and suffix k=n-i
           // prefix k=i+1 max L[i], suffix k=n-i max R[i]
           if L[i] + R[i] < a[i] {
               ok = false
               break
           }
       }
       if ok {
           fmt.Fprintln(out, "YES")
       } else {
           fmt.Fprintln(out, "NO")
       }
   }
}
