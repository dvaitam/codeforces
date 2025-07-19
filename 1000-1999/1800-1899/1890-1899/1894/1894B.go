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
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       v := make([]int, n)
       cnt := make([]int, 101)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &v[i])
           if v[i] >= 1 && v[i] <= 100 {
               cnt[v[i]]++
           }
       }
       // find first two values with count > 1
       x := make([]int, 0, 2)
       for i := 1; i <= 100 && len(x) < 2; i++ {
           if cnt[i] > 1 {
               x = append(x, i)
           }
       }
       if len(x) < 2 {
           fmt.Fprintln(writer, -1)
           continue
       }
       // labels default to 1
       b := make([]int, n)
       for i := range b {
           b[i] = 1
       }
       // mark all occurrences beyond first for x[0] as 2
       cur := make([]int, 101)
       for i, val := range v {
           if val >= 1 && val <= 100 {
               cur[val]++
               if cur[val] > 1 && val == x[0] {
                   b[i] = 2
               }
           }
       }
       // mark all occurrences beyond first for x[1] as 3
       for i := range cur {
           cur[i] = 0
       }
       for i, val := range v {
           if val >= 1 && val <= 100 {
               cur[val]++
               if cur[val] > 1 && val == x[1] {
                   b[i] = 3
               }
           }
       }
       // output labels
       for i, bi := range b {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, bi)
       }
       writer.WriteByte('\n')
   }
}
