package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   c := make([]int, n+1)
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &c[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   f := make([]bool, n+1)
   used := make([]bool, n+1)
   s := make([]int, 0, n)
   ans := 0

   var work func(int)
   work = func(i int) {
       cur := i
       s = s[:0]
       for {
           if f[cur] {
               if !used[cur] {
                   for _, v := range s {
                       used[v] = false
                   }
                   return
               }
               // found a cycle, compute minimal cost in it
               res := c[cur]
               for j := len(s) - 1; j >= 0; j-- {
                   v := s[j]
                   res = min(res, c[v])
                   if v == cur {
                       break
                   }
               }
               ans += res
               for _, v := range s {
                   used[v] = false
               }
               return
           }
           f[cur] = true
           used[cur] = true
           s = append(s, cur)
           cur = a[cur]
       }
   }

   for i := 1; i <= n; i++ {
       if !f[i] {
           work(i)
       }
   }

   fmt.Fprintln(writer, ans)
}
