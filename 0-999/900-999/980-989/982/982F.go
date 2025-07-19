package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m int
   edges [][]int
   check []bool
)
// variables for solve
var (
   c []int
   pre []int
   p []int
   tot int
   tag []bool
   pos []int
   flagArr []bool
   sumArr []int
   resCount int
)

func dfs1(x int) {
   if tot > 0 || check[x] {
       return
   }
   c[x] = -1
   for _, v := range edges[x] {
       if check[v] {
           continue
       }
       if c[v] == -1 {
           // found cycle, backtrack
           y := x
           for y != v {
               p = append(p, y)
               tot++
               y = pre[y]
           }
           p = append(p, v)
           tot++
           return
       }
       if c[v] == 0 {
           pre[v] = x
           dfs1(v)
           if tot > 0 {
               return
           }
       }
   }
   c[x] = 1
}

func dfs2(S, x int) {
   if check[x] {
       return
   }
   flagArr[x] = true
   for _, v := range edges[x] {
       if check[v] {
           continue
       }
       if tag[v] {
           if tag[x] && tag[v] {
               continue
           }
           if pos[S] >= pos[v] {
               continue
           }
           // cross edge from cycle to cycle
           resCount++
           sumArr[1]++
           sumArr[pos[S]+1]--
           sumArr[pos[v]]++
           sumArr[tot+1]--
       } else if !flagArr[v] {
           dfs2(S, v)
       }
   }
}

func dfs3(S, x int) {
   if check[x] {
       return
   }
   flagArr[x] = true
   for _, v := range edges[x] {
       if check[v] {
           continue
       }
       if tag[v] {
           if tag[x] && tag[v] {
               continue
           }
           if pos[S] < pos[v] {
               continue
           }
           resCount++
           sumArr[pos[v]]++
           sumArr[pos[S]+1]--
       } else if !flagArr[v] {
           dfs3(S, v)
       }
   }
}

// solve finds a vertex in all cycles or 0 if none
func solve() int {
   c = make([]int, n+1)
   pre = make([]int, n+1)
   p = p[:0]
   tot = 0
   resCount = 0
   // find a cycle
   for i := 1; i <= n && tot == 0; i++ {
       if c[i] == 0 {
           dfs1(i)
       }
   }
   if tot == 0 {
       return 0
   }
   // reverse p to correct order
   for i, j := 0, len(p)-1; i < j; i, j = i+1, j-1 {
       p[i], p[j] = p[j], p[i]
   }
   tag = make([]bool, n+1)
   pos = make([]int, n+1)
   sumArr = make([]int, tot+2)
   for i, v := range p {
       idx := i + 1
       tag[v] = true
       pos[v] = idx
   }
   // dfs2
   flagArr = make([]bool, n+1)
   for _, v := range p {
       dfs2(v, v)
       // reset flags
       for i := range flagArr {
           flagArr[i] = false
       }
   }
   // dfs3
   for i := len(p) - 1; i >= 0; i-- {
       v := p[i]
       dfs3(v, v)
       for j := range flagArr {
           flagArr[j] = false
       }
   }
   // prefix sums to find position
   for i := 1; i <= tot; i++ {
       sumArr[i] += sumArr[i-1]
       if sumArr[i] == resCount {
           po := p[i-1]
           check[po] = true
           return po
       }
   }
   return 0
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n, &m)
   edges = make([][]int, n+1)
   check = make([]bool, n+1)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       edges[a] = append(edges[a], b)
   }
   ans := solve()
   if ans == 0 {
       fmt.Println(-1)
       return
   }
   if solve() == 0 {
       fmt.Println(ans)
   } else {
       fmt.Println(-1)
   }
}
