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
   var Q int
   fmt.Fscan(reader, &Q)
   for ; Q > 0; Q-- {
       var N, K int
       fmt.Fscan(reader, &N, &K)
       if K == 0 || K == 1 {
           for i := 1; i <= N; i++ {
               if i > 1 {
                   fmt.Fprint(writer, " ")
               }
               fmt.Fprint(writer, i)
           }
           fmt.Fprint(writer, "\n")
           fmt.Fprint(writer, N, "\n")
           for i := 1; i <= N; i++ {
               if i > 1 {
                   fmt.Fprint(writer, " ")
               }
               fmt.Fprint(writer, i)
           }
           fmt.Fprint(writer, "\n")
       } else {
           ans := make([]int, 0, N)
           for i := 0; i < N; i += K {
               s := i
               t := i + K - 1
               if t >= N {
                   t = N - 1
               }
               length := t - s + 1
               if length == 1 {
                   ans = append(ans, s)
               } else if length == 2 {
                   ans = append(ans, s)
                   ans = append(ans, t)
               } else if length%2 == 1 {
                   x := length / 2
                   for j := s + x - 1; j >= s; j-- {
                       ans = append(ans, j)
                   }
                   for j := t; j > s + x - 1; j-- {
                       ans = append(ans, j)
                   }
               } else {
                   x := length/2 - 1
                   for j := s + x - 1; j >= s; j-- {
                       ans = append(ans, j)
                   }
                   for j := t; j > s + x - 1; j-- {
                       ans = append(ans, j)
                   }
               }
           }
           for i, v := range ans {
               if i > 0 {
                   fmt.Fprint(writer, " ")
               }
               fmt.Fprint(writer, v+1)
           }
           fmt.Fprint(writer, "\n")
           m := (N + K - 1) / K
           fmt.Fprint(writer, m, "\n")
           for i := 0; i < N; i++ {
               if i > 0 {
                   fmt.Fprint(writer, " ")
               }
               fmt.Fprint(writer, i/K+1)
           }
           fmt.Fprint(writer, "\n")
       }
   }
}
