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
   adj := make([][]int, n+1)
   for i := 1; i < n; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       adj[x] = append(adj[x], y)
       adj[y] = append(adj[y], x)
   }
   a := make([]int, n+1)
   b := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &b[i])
   }

   var ans []int
   type stackElem struct {
       x, fa, fc, fs, idx, cur int
   }
   // iterative DFS using explicit stack
   stack := []stackElem{{x: 1, fa: 0, fc: 0, fs: 0, idx: 0, cur: -1}}
   for len(stack) > 0 {
       // pop
       e := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       x, fa, fc, fs, idx, cur := e.x, e.fa, e.fc, e.fs, e.idx, e.cur
       if idx == 0 {
           cur = a[x] ^ fc ^ b[x]
           if cur == 1 {
               ans = append(ans, x)
           }
       }
       if idx < len(adj[x]) {
           // push state back with next index
           e.idx = idx + 1
           e.cur = cur
           stack = append(stack, e)
           y := adj[x][idx]
           if y != fa {
               // push child
               stack = append(stack, stackElem{x: y, fa: x, fc: fs, fs: fc ^ cur, idx: 0, cur: -1})
           }
       }
   }
   sort.Ints(ans)
   fmt.Fprintln(writer, len(ans))
   for _, v := range ans {
       fmt.Fprintln(writer, v)
   }
}
