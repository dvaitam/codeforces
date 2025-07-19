package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type pair struct {
   val, idx int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       b := make([]pair, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &b[i].val)
           b[i].idx = i
       }
       sort.Slice(b, func(i, j int) bool {
           return b[i].val < b[j].val
       })
       ans := make([]int, n)
       l, r := 0, n-1
       value := n
       cntP, cntM := 0, 0
       flag := true
       for l <= r {
           if b[l].val == cntP {
               cur := b[l].val
               for l <= r && b[l].val == cur {
                   ans[b[l].idx] = -value
                   cntM++
                   l++
               }
               value--
           } else if b[r].val == n-cntM {
               cur := b[r].val
               for l <= r && b[r].val == cur {
                   ans[b[r].idx] = value
                   cntP++
                   r--
               }
               value--
           } else {
               flag = false
               break
           }
       }
       if !flag {
           fmt.Fprintln(writer, "NO")
       } else {
           fmt.Fprintln(writer, "YES")
           for i := 0; i < n; i++ {
               fmt.Fprint(writer, ans[i])
               if i+1 < n {
                   fmt.Fprint(writer, " ")
               } else {
                   fmt.Fprint(writer, "\n")
               }
           }
       }
   }
}
