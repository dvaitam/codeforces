package main

import (
   "bufio"
   "fmt"
   "os"
)

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var n, k int
   fmt.Fscan(reader, &n, &k)
   idx := make([][]int, n)
   xs := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &xs[i])
       xs[i]--
       idx[xs[i]] = append(idx[xs[i]], i)
   }
   ans := make([]int, n)
   left := 0
   for _, a := range idx {
       if len(a) >= k {
           for c := 1; c <= k; c++ {
               ans[a[c-1]] = c
           }
       } else {
           left += len(a)
       }
   }
   times := left / k
   if times > 0 {
       c := 1
       for _, a := range idx {
           if len(a) >= k {
               continue
           }
           for _, i := range a {
               if times == 0 {
                   break
               }
               ans[i] = c
               c++
               if c > k {
                   c = 1
                   times--
               }
           }
           if times == 0 {
               break
           }
       }
   }
   for i, v := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprint(v))
   }
   writer.WriteByte('\n')
}

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
