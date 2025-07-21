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
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   var m int
   fmt.Fscan(reader, &m)
   b := make([]int, m)
   for j := 0; j < m; j++ {
       fmt.Fscan(reader, &b[j])
   }

   // dp[j]: length of LCIS ending with b[j]
   dp := make([]int, m)
   prev := make([]int, m)
   for j := range prev {
       prev[j] = -1
   }

   for i := 0; i < n; i++ {
       current := 0
       last := -1
       for j := 0; j < m; j++ {
           if a[i] == b[j] {
               if current+1 > dp[j] {
                   dp[j] = current + 1
                   prev[j] = last
               }
           } else if a[i] > b[j] {
               if dp[j] > current {
                   current = dp[j]
                   last = j
               }
           }
       }
   }

   // find maximum
   length := 0
   endIdx := -1
   for j := 0; j < m; j++ {
       if dp[j] > length {
           length = dp[j]
           endIdx = j
       }
   }

   // reconstruct sequence
   seq := make([]int, 0, length)
   for idx := endIdx; idx != -1; idx = prev[idx] {
       seq = append(seq, b[idx])
   }
   // reverse
   for i, j := 0, len(seq)-1; i < j; i, j = i+1, j-1 {
       seq[i], seq[j] = seq[j], seq[i]
   }

   // output
   fmt.Fprintln(writer, length)
   // print sequence (blank line if empty)
   for i, v := range seq {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
