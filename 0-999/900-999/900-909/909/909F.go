package main

import (
   "bufio"
   "fmt"
   "math/bits"
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
   t := n

   // Subtask 1: p[i] != i and p[i] & i == 0
   if n%2 == 1 {
       fmt.Fprintln(writer, "NO")
   } else {
       fmt.Fprintln(writer, "YES")
       p := make([]int, n+1)
       nn := n
       for n > 0 {
           // highest power of two <= n
           m := 1 << (bits.Len(uint(n)) - 1)
           n -= m
           for i := 0; i <= n; i++ {
               p[m+i] = m - 1 - i
               p[m-1-i] = m + i
           }
           n = m - n - 2
       }
       // output p[1..nn]
       for i := 1; i <= nn; i++ {
           if i > 1 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, p[i])
       }
       writer.WriteByte('\n')
   }

   // Subtask 2: q[i] != i and q[i] & i != 0
   n = t
   if n < 6 || bits.OnesCount(uint(n)) == 1 {
       fmt.Fprintln(writer, "NO")
   } else {
       fmt.Fprintln(writer, "YES")
       // special case n == 7
       if n == 7 {
           // fixed solution for n=7
           fmt.Fprintln(writer, "5 3 2 6 1 7 4")
       } else {
           p := make([]int, n+1)
           // highest power of two <= n
           m := 1 << (bits.Len(uint(n)) - 1)
           p[m] = n
           p[n] = m
           rem := n - m - 1
           // build pairs
           if rem%2 == 1 {
               for i := 1; i <= rem; i++ {
                   p[i] = m + i
                   p[m+i] = i
               }
               for i := rem + 1; i < m-1; i += 2 {
                   p[i] = i + 1
                   p[i+1] = i
               }
           } else {
               for i := 1; i <= rem; i += 2 {
                   p[m+i] = m + i + 1
                   p[m+i+1] = m + i
               }
               // small adjustments
               p[1] = m - 1
               p[m-1] = 1
               for i := 2; i < m-4; i += 2 {
                   p[i] = i + 1
                   p[i+1] = i
               }
               p[m-2] = m-3
               p[m-3] = m-4
               p[m-4] = m-2
           }
           // output p[1..n]
           for i := 1; i <= n; i++ {
               if i > 1 {
                   writer.WriteByte(' ')
               }
               fmt.Fprint(writer, p[i])
           }
           writer.WriteByte('\n')
       }
   }
}
