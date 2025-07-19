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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   type pair struct { val, idx int }
   a := make([]pair, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i].val)
       a[i].idx = i
   }
   sort.Slice(a, func(i, j int) bool { return a[i].val > a[j].val })

   // Prepare answer matrix: (n+1) rows of length n, initialized to '0'
   ans := make([][]byte, n+1)
   for i := 0; i <= n; i++ {
       row := make([]byte, n)
       for j := range row {
           row[j] = '0'
       }
       ans[i] = row
   }
   used := make([]bool, n+1)

   for i := 0; i < n; i++ {
       for j := 0; j < a[i].val; j++ {
           idx := (i + j) % (n + 1)
           ans[idx][a[i].idx] = '1'
           used[idx] = true
       }
   }
   // Count and output
   cnt := 0
   for i := 0; i <= n; i++ {
       if used[i] {
           cnt++
       }
   }
   fmt.Fprintln(writer, cnt)
   for i := 0; i <= n; i++ {
       if used[i] {
           writer.Write(ans[i])
           writer.WriteByte('\n')
       }
   }
}
