package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   type domino struct {
       x, h int
       idx  int
   }
   dom := make([]domino, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &dom[i].x, &dom[i].h)
       dom[i].idx = i
   }
   sort.Slice(dom, func(i, j int) bool { return dom[i].x < dom[j].x })

   far := make([]int, n)
   R := make([]int, n)
   stack := make([]int, 0, n)
   // process from right to left
   for i := n - 1; i >= 0; i-- {
       far[i] = dom[i].x + dom[i].h - 1
       R[i] = i
       // absorb all dominoes that fall within reach
       for len(stack) > 0 {
           j := stack[len(stack)-1]
           if dom[j].x > far[i] {
               break
           }
           // pop j
           stack = stack[:len(stack)-1]
           if far[j] > far[i] {
               far[i] = far[j]
           }
           if R[j] > R[i] {
               R[i] = R[j]
           }
       }
       stack = append(stack, i)
   }

   // prepare answers
   ans := make([]int, n)
   for i := 0; i < n; i++ {
       count := R[i] - i + 1
       ans[dom[i].idx] = count
   }
   // output
   for i, v := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
