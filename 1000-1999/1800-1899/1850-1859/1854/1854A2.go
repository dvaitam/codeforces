package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct{ x, y int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       solve(reader, writer)
   }
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n+1)
   idp := make([]int, 0, n)
   idn := make([]int, 0, n)
   mx, mn := -1000000001, 1000000001
   mxid, mnid := 0, 0
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] > 0 {
           idp = append(idp, i)
       } else if a[i] < 0 {
           idn = append(idn, i)
       }
       if a[i] > mx {
           mx = a[i]
           mxid = i
       }
       if a[i] < mn {
           mn = a[i]
           mnid = i
       }
   }
   // all non-positive
   if mx <= 0 {
       // chain backwards
       fmt.Fprintln(writer, n-1)
       for i := n - 1; i >= 1; i-- {
           fmt.Fprintln(writer, i, i+1)
       }
       return
   }
   // all non-negative
   if mn >= 0 {
       fmt.Fprintln(writer, n-1)
       for i := 2; i <= n; i++ {
           fmt.Fprintln(writer, i, i-1)
       }
       return
   }
   // mixed
   var ans []pair
   // try positive amplification
   k := mx
   for k+mn < 0 {
       ans = append(ans, pair{mxid, mxid})
       k *= 2
   }
   for _, idx := range idn {
       ans = append(ans, pair{idx, mxid})
   }
   for i := 2; i <= n; i++ {
       ans = append(ans, pair{i, i - 1})
   }
   if len(ans) <= 31 {
       fmt.Fprintln(writer, len(ans))
       for _, p := range ans {
           fmt.Fprintln(writer, p.x, p.y)
       }
       return
   }
   // try negative amplification
   ans = ans[:0]
   k = mn
   for k+mx > 0 {
       ans = append(ans, pair{mnid, mnid})
       k *= 2
   }
   for _, idx := range idp {
       ans = append(ans, pair{idx, mnid})
   }
   for i := n - 1; i >= 1; i-- {
       ans = append(ans, pair{i, i + 1})
   }
   if len(ans) <= 31 {
       fmt.Fprintln(writer, len(ans))
       for _, p := range ans {
           fmt.Fprintln(writer, p.x, p.y)
       }
       return
   }
}
