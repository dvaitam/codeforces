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

   var n, w int
   fmt.Fscan(reader, &n, &w)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   b := make([]int64, w)
   for i := 0; i < w; i++ {
       fmt.Fscan(reader, &b[i])
   }
   if w == 1 {
       // Any single tower matches after shift
       fmt.Fprintln(writer, n)
       return
   }
   // Build difference sequences
   m := w - 1
   pattern := make([]int64, m)
   for i := 0; i < m; i++ {
       pattern[i] = b[i+1] - b[i]
   }
   textLen := n - 1
   text := make([]int64, textLen)
   for i := 0; i < textLen; i++ {
       text[i] = a[i+1] - a[i]
   }
   // Build prefix function for KMP
   pi := make([]int, m)
   for i, j := 1, 0; i < m; i++ {
       for j > 0 && pattern[i] != pattern[j] {
           j = pi[j-1]
       }
       if pattern[i] == pattern[j] {
           j++
       }
       pi[i] = j
   }
   // KMP search
   count := 0
   for i, j := 0, 0; i < textLen; i++ {
       for j > 0 && text[i] != pattern[j] {
           j = pi[j-1]
       }
       if text[i] == pattern[j] {
           j++
       }
       if j == m {
           count++
           j = pi[j-1]
       }
   }
   fmt.Fprintln(writer, count)
}
