package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
)

type pair struct {
   val int
   idx int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   A := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &A[i])
   }

   P := make([]pair, n)
   for i := 0; i < n; i++ {
       P[i] = pair{val: A[i], idx: i}
   }
   sort.Slice(P, func(i, j int) bool {
       return P[i].val < P[j].val
   })

   res := make([]int64, n)
   var curr int64
   for _, p := range P {
       ai := int64(p.val)
       if ai > curr {
           curr = ai
       } else {
           curr++
       }
       res[p.idx] = curr
   }

   for i := 0; i < n; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(strconv.FormatInt(res[i], 10))
   }
   writer.WriteByte('\n')
}
