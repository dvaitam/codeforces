package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m  int
   head  []int
   toArr  []int
   nxt   []int
   C     []int
   ec, cnt int
   dfn, low, color []int
   instk, fa []bool
   stack []int
   ID    int
)

func addEdge(x, y int) {
   ec++
   toArr[ec] = y
   nxt[ec] = head[x]
   head[x] = ec
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func tarjan(x int) {
   instk[x], fa[x] = true, true
   stack = append(stack, x)
   ID++
   dfn[x], low[x] = ID, ID
   for i := head[x]; i != 0; i = nxt[i] {
       y := toArr[i]
       if dfn[y] == 0 {
           C[i] = color[x] ^ 1
           color[y] = color[x] ^ 1
           tarjan(y)
           low[x] = min(low[x], low[y])
       } else if instk[y] {
           low[x] = min(low[x], dfn[y])
           if fa[y] {
               C[i] = color[x] ^ 1
           } else {
               C[i] = color[y]
           }
       } else {
           C[i] = 1
       }
   }
   if dfn[x] == low[x] {
       cnt++
       for {
           v := stack[len(stack)-1]
           stack = stack[:len(stack)-1]
           instk[v] = false
           if v == x {
               break
           }
       }
   }
   fa[x] = false
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &m)
   head = make([]int, n+1)
   toArr = make([]int, m+1)
   nxt = make([]int, m+1)
   C = make([]int, m+1)
   dfn = make([]int, n+1)
   low = make([]int, n+1)
   color = make([]int, n+1)
   instk = make([]bool, n+1)
   fa = make([]bool, n+1)
   for i := 1; i <= m; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       addEdge(x, y)
   }
   for i := 1; i <= n; i++ {
       if dfn[i] == 0 {
           color[i] = 1
           tarjan(i)
       }
   }
   if cnt == n {
       fmt.Fprint(writer, "1\n")
       for i := 1; i <= m; i++ {
           fmt.Fprint(writer, "1 ")
       }
       fmt.Fprint(writer, "\n")
   } else {
       fmt.Fprint(writer, "2\n")
       for i := 1; i <= m; i++ {
           fmt.Fprintf(writer, "%d ", C[i]+1)
       }
       fmt.Fprint(writer, "\n")
   }
}
