package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(in, &n)
       var res []int
       switch n {
       case 2:
           res = []int{1, 2}
       case 3:
           res = []int{1, 2, 3}
       default:
           if n%2 == 0 {
               res = make([]int, n)
               for i := 0; i < n; i++ {
                   if i == 0 {
                       res[i] = n/2 + 1
                   } else if i%2 == 0 {
                       res[i] = (n+i)/2 + 1
                   } else {
                       res[i] = i/2 + 1
                   }
               }
           } else {
               res = make([]int, 0, n)
               res = append(res, n)
               nn := n - 1
               for i := 0; i < nn; i++ {
                   if i%2 == 0 {
                       res = append(res, (nn+i)/2 + 1)
                   } else {
                       res = append(res, i/2 + 1)
                   }
               }
           }
       }
       for i, v := range res {
           if i > 0 {
               fmt.Fprint(out, " ")
           }
           fmt.Fprint(out, v)
       }
       fmt.Fprint(out, '\n')
   }
}
