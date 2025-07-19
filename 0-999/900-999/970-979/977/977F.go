package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n+1)
   hs := make([]int, 0, n*3)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
       hs = append(hs, a[i], a[i]-1, a[i]+1)
   }
   sort.Ints(hs)
   // unique values
   uniq := make([]int, 0, len(hs))
   for i, v := range hs {
       if i == 0 || v != hs[i-1] {
           uniq = append(uniq, v)
       }
   }
   // coordinate compression
   for i := 1; i <= n; i++ {
       x := a[i]
       idx := sort.SearchInts(uniq, x)
       a[i] = idx + 1
   }
   m := len(uniq)
   f := make([]int, m+2)
   pre := make([]int, n+1)
   pos := make([]int, m+2)
   ans, anspos := 0, 0
   for i := 1; i <= n; i++ {
       x := a[i]
       if f[x-1]+1 > f[x] {
           f[x] = f[x-1] + 1
           pre[i] = pos[x-1]
           pos[x] = i
       }
       if f[x] > ans {
           ans = f[x]
           anspos = i
       }
   }
   fmt.Fprintln(writer, ans)
   // reconstruct sequence of positions
   seq := make([]int, 0, ans)
   for p := anspos; p != 0; p = pre[p] {
       seq = append(seq, p)
   }
   // reverse
   for i, j := 0, len(seq)-1; i < j; i, j = i+1, j-1 {
       seq[i], seq[j] = seq[j], seq[i]
   }
   for i, v := range seq {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
