package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &p[i])
   }
   // adjacency via to/nxt/first as in C++ code
   to := make([]int, m+2)
   nxt := make([]int, m+2)
   first := make([]int, n+2)
   bj := make([]bool, n+2)
   flag := make([]bool, n+2)
   tot := 1
   // add edge u->v
   add := func(u, v int) {
       to[tot] = v
       nxt[tot] = first[u]
       first[u] = tot
       tot++
   }
   var u, v int
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &u, &v)
       add(u, v)
       if v == p[n] {
           bj[u] = true
       }
   }
   ans := 0
   a := make([]int, n+2)
   num := 0
   // process from n-1 down to 1
   for i := n - 1; i >= 1; i-- {
       pi := p[i]
       if bj[pi] {
           ans++
           // mark flags for neighbors of pi
           for j := first[pi]; j != 0; j = nxt[j] {
               flag[to[j]] = true
           }
           // check previous a entries
           for j := 1; j <= num; j++ {
               idx := a[j]
               if !flag[p[idx]] {
                   ans--
                   num++
                   a[num] = i
                   break
               }
           }
           // reset flags
           for j := first[pi]; j != 0; j = nxt[j] {
               flag[to[j]] = false
           }
       } else {
           num++
           a[num] = i
       }
   }
   // output answer
   fmt.Println(ans)
}
