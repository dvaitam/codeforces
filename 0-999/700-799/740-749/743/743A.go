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

   var n, a, b int
   if _, err := fmt.Fscan(reader, &n, &a, &b); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)

   // If start and end airports belong to same company, cost is 0
   if s[a-1] == s[b-1] {
       fmt.Fprintln(writer, 0)
       return
   }

   // Collect positions of each company
   var pos [2][]int
   for i, ch := range s {
       c := int(ch - '0')
       pos[c] = append(pos[c], i+1)
   }
   cs := int(s[a-1] - '0')
   ct := int(s[b-1] - '0')
   src := pos[cs]
   dst := pos[ct]

   // Two pointers to find minimal absolute difference
   // Maximum possible distance is n-1, so initialize ans = n
   i, j, ans := 0, 0, n
   for i < len(src) && j < len(dst) {
       di := src[i]
       dj := dst[j]
       diff := di - dj
       if diff < 0 {
           diff = -diff
       }
       if diff < ans {
           ans = diff
       }
       // Move pointer at smaller value
       if src[i] < dst[j] {
           i++
       } else {
           j++
       }
   }
   fmt.Fprintln(writer, ans)
}
