package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   p := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &p[i])
   }
   var t int
   fmt.Fscan(reader, &t)
   isForcedClose := make([]bool, n)
   for i := 0; i < t; i++ {
       var q int
       fmt.Fscan(reader, &q)
       if q >= 1 && q <= n {
           isForcedClose[q-1] = true
       }
   }
   res := make([]int, n)
   stack := make([]int, 0, n)
   for i := n - 1; i >= 0; i-- {
       if isForcedClose[i] {
           stack = append(stack, p[i])
           res[i] = -p[i]
       } else {
           if len(stack) > 0 && stack[len(stack)-1] == p[i] {
               // match as opening
               stack = stack[:len(stack)-1]
               res[i] = p[i]
           } else {
               // make closing
               stack = append(stack, p[i])
               res[i] = -p[i]
           }
       }
   }
   if len(stack) != 0 {
       writer.WriteString("NO\n")
       return
   }
   writer.WriteString("YES\n")
   tmp := make([]byte, 0, 20)
   for i, v := range res {
       if i > 0 {
           writer.WriteByte(' ')
       }
       tmp = tmp[:0]
       tmp = strconv.AppendInt(tmp, int64(v), 10)
       writer.Write(tmp)
   }
   writer.WriteByte('\n')
}
