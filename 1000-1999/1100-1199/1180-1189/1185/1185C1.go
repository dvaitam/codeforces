package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, M int
   fmt.Fscan(reader, &n, &M)
   t := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &t[i])
   }
   res := make([]int, n)
   prev := make([]int, 0, n)
   sum := 0
   for i := 0; i < n; i++ {
       need := sum + t[i] - M
       if need <= 0 {
           res[i] = 0
       } else {
           tmp := make([]int, len(prev))
           copy(tmp, prev)
           sort.Sort(sort.Reverse(sort.IntSlice(tmp)))
           removed, acc := 0, 0
           for _, v := range tmp {
               acc += v
               removed++
               if acc >= need {
                   break
               }
           }
           res[i] = removed
       }
       prev = append(prev, t[i])
       sum += t[i]
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i, v := range res {
       if i > 0 {
           writer.WriteString(" ")
       }
       writer.WriteString(fmt.Sprint(v))
   }
   writer.WriteByte('\n')
}
