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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   L := 2*n - 1
   a := make([]int, L)
   for i := 0; i < L; i++ {
       fmt.Fscan(reader, &a[i])
   }
   if n == 1 {
       // single node, sequence length 1
       if a[0] == 0 {
           a[0] = 1
       }
       if a[0] != 1 {
           fmt.Fprintln(writer, "no")
       } else {
           fmt.Fprintln(writer, "yes")
           fmt.Fprintln(writer, "1")
       }
       return
   }
   // determine center s at odd positions (1-based)
   s := 0
   for i := 0; i < L; i += 2 {
       if a[i] != 0 {
           if s == 0 {
               s = a[i]
           } else if a[i] != s {
               fmt.Fprintln(writer, "no")
               return
           }
       }
   }
   if s == 0 {
       s = 1
   }
   // check boundaries at first and last
   if a[0] != 0 && a[0] != s {
       fmt.Fprintln(writer, "no")
       return
   }
   if a[L-1] != 0 && a[L-1] != s {
       fmt.Fprintln(writer, "no")
       return
   }
   // leaves appear at even positions (2-based), i.e., i%2==1
   used := make([]bool, n+1)
   for i := 1; i < L; i += 2 {
       if a[i] != 0 {
           if a[i] == s || a[i] < 1 || a[i] > n {
               fmt.Fprintln(writer, "no")
               return
           }
           if used[a[i]] {
               fmt.Fprintln(writer, "no")
               return
           }
           used[a[i]] = true
       }
   }
   // collect unused leaves
   leaves := make([]int, 0, n-1)
   for v := 1; v <= n; v++ {
       if v == s {
           continue
       }
       if !used[v] {
           leaves = append(leaves, v)
       }
   }
   // fill sequence
   b := make([]int, L)
   idx := 0
   for i := 0; i < L; i++ {
       if i%2 == 0 {
           b[i] = s
       } else {
           if a[i] != 0 {
               b[i] = a[i]
           } else {
               b[i] = leaves[idx]
               idx++
           }
       }
   }
   // output
   fmt.Fprintln(writer, "yes")
   for i, v := range b {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
