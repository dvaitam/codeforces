package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
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
   // g[i]: length of decreasing run ending at i from left
   g := make([]int, n)
   g[0] = 1
   for i := 1; i < n; i++ {
       if a[i-1] > a[i] {
           g[i] = g[i-1] + 1
       } else {
           g[i] = 1
       }
   }
   // f[i]: length of increasing run starting at i to right
   f := make([]int, n)
   f[n-1] = 1
   for i := n - 2; i >= 0; i-- {
       if a[i+1] > a[i] {
           f[i] = f[i+1] + 1
       } else {
           f[i] = 1
       }
   }

   l, r := 0, n-1
   now, ans := 0, 0
   var sb strings.Builder
   for l <= r && (a[l] > now || a[r] > now) {
       j := -1
       if a[l] > now && a[r] > now {
           if a[l] == a[r] {
               if f[l] > g[r] {
                   j = l
               } else {
                   j = r
               }
           } else if a[l] < a[r] {
               j = l
           } else {
               j = r
           }
       } else if a[l] > now {
           j = l
       } else if a[r] > now {
           j = r
       }
       now = a[j]
       ans++
       if j == l {
           sb.WriteByte('L')
           l++
       } else {
           sb.WriteByte('R')
           r--
       }
   }
   fmt.Fprintln(writer, ans)
   fmt.Fprint(writer, sb.String())
}
