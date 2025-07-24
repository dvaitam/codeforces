package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var m int
   if _, err := fmt.Fscan(reader, &m); err != nil {
       return
   }
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   res := make([]byte, 0, (n+m-1)/m)
   for i := 0; i < n; i += m {
       end := i + m
       if end > n {
           end = n
       }
       minc := byte('z' + 1)
       for j := i; j < end; j++ {
           if s[j] < minc {
               minc = s[j]
           }
       }
       res = append(res, minc)
   }
   sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   writer.Write(res)
}
