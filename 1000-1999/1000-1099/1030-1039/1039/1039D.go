package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n       int
   edges   [][]int
   vs, ps  []int
   tptr    int
   m1, m2  []int
   answers []int
)

func solve(k int) int {
   // reset m1, m2
   for i := 0; i < tptr; i++ {
       m1[i] = 0
       m2[i] = 0
   }
   ans := 0
   for i := 0; i < tptr; i++ {
       v := vs[i]
       p := ps[i]
       if m1[v]+m2[v]+1 >= k {
           ans++
       } else if v != 0 {
           if m1[v]+1 > m2[p] {
               m2[p] = m1[v] + 1
           }
           if m2[p] > m1[p] {
               m1[p], m2[p] = m2[p], m1[p]
           }
       }
   }
   return ans
}

func compute(x, y, minv, maxv int) {
   if x > y {
       return
   }
   m := (x + y) >> 1
   if minv == maxv {
       answers[m] = minv
   } else {
       answers[m] = solve(m)
   }
   compute(x, m-1, answers[m], maxv)
   compute(m+1, y, minv, answers[m])
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n)
   edges = make([][]int, n)
   for i := 0; i < n-1; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       a--
       b--
       edges[a] = append(edges[a], b)
       edges[b] = append(edges[b], a)
   }
   // iterative DFS for post-order
   vs = make([]int, n)
   ps = make([]int, n)
   tptr = 0
   // stack elements: v, parent, state (0 enter, 1 exit)
   stack := make([][3]int, 0, n*2)
   stack = append(stack, [3]int{0, -1, 0})
   for len(stack) > 0 {
       elem := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       v, p, state := elem[0], elem[1], elem[2]
       if state == 0 {
           // enter
           stack = append(stack, [3]int{v, p, 1})
           for _, c := range edges[v] {
               if c == p {
                   continue
               }
               stack = append(stack, [3]int{c, v, 0})
           }
       } else {
           // exit
           vs[tptr] = v
           ps[tptr] = p
           tptr++
       }
   }
   // prepare auxiliary arrays
   m1 = make([]int, n)
   m2 = make([]int, n)
   answers = make([]int, n+1)
   const C = 1000
   // initial solves
   for i := 1; i <= C && i <= n; i++ {
       answers[i] = solve(i)
   }
   if C < n {
       compute(C+1, n, 0, answers[C])
   }
   // output
   for i := 1; i <= n; i++ {
       fmt.Fprintln(out, answers[i])
   }
}
