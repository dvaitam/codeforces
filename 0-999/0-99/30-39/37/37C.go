package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   maxDepth := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] > maxDepth {
           maxDepth = a[i]
       }
   }
   // counts of code lengths
   b := make([]int, maxDepth+2)
   // indices of symbols by depth
   q := make([][]int, maxDepth+2)
   for i, depth := range a {
       b[depth]++
       q[depth] = append(q[depth], i)
   }
   // feasibility check using available tree nodes
   var cur int64 = 2
   const inf = 100000000
   for depth := 1; depth <= maxDepth; depth++ {
       if int64(b[depth]) > cur {
           fmt.Println("NO")
           return
       }
       cur = (cur - int64(b[depth])) * 2
       if cur > inf {
           break
       }
   }
   // prepare answer storage
   ans := make([][]byte, n)
   for i := range ans {
       ans[i] = make([]byte, a[i])
   }
   c := make([]byte, maxDepth+2)
   q1 := make([]int, maxDepth+2)
   all := n
   END := false
   var dfs func(depth int)
   dfs = func(depth int) {
       if END {
           return
       }
       if depth <= maxDepth && b[depth] > 0 {
           b[depth]--
           all--
           idx := q[depth][q1[depth]]
           // assign bits to answer
           for i := 0; i < depth; i++ {
               ans[idx][i] = c[i]
           }
           q1[depth]++
           if all == 0 {
               END = true
           }
           return
       }
       if depth > maxDepth {
           return
       }
       // explore left subtree
       c[depth] = '0'
       dfs(depth + 1)
       if END {
           return
       }
       // explore right subtree
       c[depth] = '1'
       dfs(depth + 1)
   }
   dfs(0)
   // output results
   fmt.Println("YES")
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < n; i++ {
       writer.Write(ans[i])
       writer.WriteByte('\n')
   }
}
