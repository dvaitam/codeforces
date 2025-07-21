package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

var (
   n      int
   a      []int64
   tree   [][]int
   childs [][]int
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read n
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a = make([]int64, n+1)
   for i := 1; i <= n; i++ {
       var x int64
       fmt.Fscan(reader, &x)
       a[i] = x
   }
   tree = make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       tree[u] = append(tree[u], v)
       tree[v] = append(tree[v], u)
   }
   // build rooted tree at 1
   childs = make([][]int, n+1)
   parent := make([]int, n+1)
   parent[1] = 0
   queue := make([]int, 0, n)
   queue = append(queue, 1)
   for i := 0; i < len(queue); i++ {
       v := queue[i]
       for _, u := range tree[v] {
           if u == parent[v] {
               continue
           }
           parent[u] = v
           childs[v] = append(childs[v], u)
           queue = append(queue, u)
       }
   }
   // dfs to balance
   removals, _ := dfs(1)
   // output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   writer.WriteString(strconv.FormatInt(removals, 10))
   writer.WriteByte('\n')
}

// dfs returns (removals needed, subtree weight after balancing)
func dfs(v int) (int64, int64) {
   if len(childs[v]) == 0 {
       return 0, a[v]
   }
   var totalRem int64
   var sumW int64
   minW := int64(1<<62 - 1)
   for _, u := range childs[v] {
       r, w := dfs(u)
       totalRem += r
       sumW += w
       if w < minW {
           minW = w
       }
   }
   cnt := int64(len(childs[v]))
   // reduce all to minW
   totalRem += sumW - minW*cnt
   return totalRem, minW * cnt
}
