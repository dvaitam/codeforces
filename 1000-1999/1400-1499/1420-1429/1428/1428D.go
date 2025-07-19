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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   q := make([]int, 0, n)
   head := 0
   res := make([][2]int, 0, 2*n)
   r := 0
   pre := 0
   for i := 1; i <= n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       if a == 0 {
           continue
       }
       if a == 1 {
           if head < len(q) {
               res = append(res, [2]int{q[head], i})
               head++
           } else {
               if pre == 3 {
                   res = append(res, [2]int{r, i})
               }
               r++
               res = append(res, [2]int{r, i})
           }
       } else if a == 2 {
           if pre == 3 {
               res = append(res, [2]int{r, i})
           }
           r++
           res = append(res, [2]int{r, i})
           q = append(q, r)
       } else if a == 3 {
           if pre == 3 {
               res = append(res, [2]int{r, i})
           }
           r++
           res = append(res, [2]int{r, i})
       }
       pre = a
   }
   if head < len(q) || pre == 3 {
       fmt.Fprintln(writer, -1)
       return
   }
   fmt.Fprintln(writer, len(res))
   for _, p := range res {
       fmt.Fprintln(writer, p[0], p[1])
   }
}
