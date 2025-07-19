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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       A := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &A[i])
       }
       // count frequencies and find max
       B := make([]int, n+2)
       k := 0
       for i := 0; i < n; i++ {
           v := A[i]
           if v > k {
               k = v
           }
           B[v]++
       }
       // initial validity check
       seen1 := false
       invalid := false
       for i := 1; i <= k; i++ {
           if B[i] == 1 {
               seen1 = true
           } else if B[i] == 2 {
               if seen1 {
                   invalid = true
                   break
               }
           } else {
               invalid = true
               break
           }
       }
       if invalid {
           fmt.Fprintln(writer, 0)
           continue
       }
       if n == k {
           // two splits: (0,n) and (n,0)
           fmt.Fprintln(writer, 2)
           fmt.Fprintf(writer, "0 %d\n", n)
           fmt.Fprintf(writer, "%d 0\n", n)
           continue
       }
       // check prefix [0:k)
       okPref := true
       for i := 1; i <= k; i++ {
           B[i] = 0
       }
       for i := 0; i < k; i++ {
           B[A[i]]++
       }
       for i := 1; i <= k; i++ {
           if B[i] == 0 {
               okPref = false
               break
           }
       }
       // check suffix [n-k:n)
       okSuf := true
       for i := 1; i <= k; i++ {
           B[i] = 0
       }
       for i := n - k; i < n; i++ {
           B[A[i]]++
       }
       for i := 1; i <= k; i++ {
           if B[i] == 0 {
               okSuf = false
               break
           }
       }
       // count valid splits
       cnt := 0
       if okPref {
           cnt++
       }
       if okSuf && (!okPref || k != n-k) {
           cnt++
       }
       fmt.Fprintln(writer, cnt)
       if okPref {
           fmt.Fprintf(writer, "%d %d\n", k, n-k)
       }
       if okSuf && (!okPref || k != n-k) {
           fmt.Fprintf(writer, "%d %d\n", n-k, k)
       }
   }
}
