package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   var ans []int
   for a[n-1] != a[0] {
       d1 := a[n-1]/2 - a[0]/2
       d2 := (a[n-1]+1)/2 - (a[0]+1)/2
       if d1 < d2 {
           ans = append(ans, 0)
           a[n-1] /= 2
           a[0] /= 2
       } else {
           ans = append(ans, 1)
           a[n-1] = (a[n-1] + 1) / 2
           a[0] = (a[0] + 1) / 2
       }
   }
   fmt.Fprintln(writer, len(ans))
   if len(ans) <= n {
       if len(ans) > 0 {
           for _, e := range ans {
               fmt.Fprint(writer, e, " ")
           }
           fmt.Fprintln(writer)
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   for i := 0; i < t; i++ {
       solve(reader, writer)
   }
}
