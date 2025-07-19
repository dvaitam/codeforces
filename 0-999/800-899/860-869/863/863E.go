package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   l := make([]int, n)
   r := make([]int, n)
   p := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &l[i], &r[i])
       p[i] = i
   }
   sort.Slice(p, func(i, j int) bool {
       if l[p[i]] == l[p[j]] {
           return r[p[i]] < r[p[j]]
       }
       return l[p[i]] < l[p[j]]
   })
   // Check same left endpoints or containment
   for i := 1; i < n; i++ {
       if l[p[i-1]] == l[p[i]] {
           fmt.Println(p[i-1] + 1)
           return
       }
       if r[p[i-1]] >= r[p[i]] {
           fmt.Println(p[i] + 1)
           return
       }
   }
   // Check overlap with neighbor gap
   for i := 1; i+1 < n; i++ {
       if r[p[i-1]]+1 >= l[p[i+1]] {
           fmt.Println(p[i] + 1)
           return
       }
   }
   fmt.Println(-1)
}
