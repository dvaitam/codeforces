package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m int
   t    []int
   p    [][]int
   ans  [][2]int
)

func solve(x int) {
   nd := []int{x}
   y := t[x]
   for y != -1 && y != x {
       nd = append(nd, y)
       y = t[y]
   }
   if y == -1 {
       for i := len(nd) - 2; i >= 0; i-- {
           ans = append(ans, [2]int{nd[i], nd[i+1]})
       }
       t[nd[0]] = -1
       for i := 1; i < len(nd); i++ {
           t[nd[i]] = nd[i]
       }
   } else {
       e := -1
       for i := 1; i <= n; i++ {
           if t[i] == -1 {
               e = i
               break
           }
       }
       ans = append(ans, [2]int{nd[len(nd)-1], e})
       for i := len(nd) - 2; i >= 0; i-- {
           ans = append(ans, [2]int{nd[i], nd[i+1]})
       }
       ans = append(ans, [2]int{e, nd[0]})
       for i := 0; i < len(nd); i++ {
           t[nd[i]] = nd[i]
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n, &m)
   p = make([][]int, m+1)
   for i := 1; i <= m; i++ {
       var ni int
       fmt.Fscan(reader, &ni)
       p[i] = make([]int, ni)
       for j := 0; j < ni; j++ {
           fmt.Fscan(reader, &p[i][j])
       }
   }
   t = make([]int, n+1)
   for i := 1; i <= n; i++ {
       t[i] = -1
   }
   ans = make([][2]int, 0)
   S := 0
   for i := 1; i <= m; i++ {
       for _, x := range p[i] {
           S++
           t[x] = S
       }
   }
   for i := 1; i <= n; i++ {
       if t[i] != -1 && t[i] != i {
           solve(i)
       }
   }
   fmt.Fprintln(writer, len(ans))
   for _, pr := range ans {
       fmt.Fprintln(writer, pr[0], pr[1])
   }
}
