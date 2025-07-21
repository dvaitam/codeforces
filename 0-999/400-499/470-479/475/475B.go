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
   hor := make([]byte, n)
   ver := make([]byte, m)
   // read strings
   var hs, vs string
   fmt.Fscan(reader, &hs)
   fmt.Fscan(reader, &vs)
   hor = []byte(hs)
   ver = []byte(vs)
   N := n * m
   adj := make([][]int, N)
   radj := make([][]int, N)
   // horizontal streets
   for i := 0; i < n; i++ {
       if hor[i] == '>' {
           for j := 0; j < m-1; j++ {
               u := i*m + j
               v := i*m + j + 1
               adj[u] = append(adj[u], v)
               radj[v] = append(radj[v], u)
           }
       } else {
           for j := 0; j < m-1; j++ {
               u := i*m + j
               v := i*m + j + 1
               // from v to u
               adj[v] = append(adj[v], u)
               radj[u] = append(radj[u], v)
           }
       }
   }
   // vertical streets
   for j := 0; j < m; j++ {
       if ver[j] == 'v' {
           for i := 0; i < n-1; i++ {
               u := i*m + j
               v2 := (i+1)*m + j
               adj[u] = append(adj[u], v2)
               radj[v2] = append(radj[v2], u)
           }
       } else {
           for i := 0; i < n-1; i++ {
               u := i*m + j
               v2 := (i+1)*m + j
               adj[v2] = append(adj[v2], u)
               radj[u] = append(radj[u], v2)
           }
       }
   }
   // DFS on graph
   visited := make([]bool, N)
   var stack []int
   stack = append(stack, 0)
   visited[0] = true
   cnt := 1
   for len(stack) > 0 {
       u := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       for _, v := range adj[u] {
           if !visited[v] {
               visited[v] = true
               cnt++
               stack = append(stack, v)
           }
       }
   }
   if cnt != N {
       fmt.Println("NO")
       return
   }
   // DFS on reverse graph
   visited = make([]bool, N)
   stack = stack[:0]
   stack = append(stack, 0)
   visited[0] = true
   cnt = 1
   for len(stack) > 0 {
       u := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       for _, v := range radj[u] {
           if !visited[v] {
               visited[v] = true
               cnt++
               stack = append(stack, v)
           }
       }
   }
   if cnt != N {
       fmt.Println("NO")
   } else {
       fmt.Println("YES")
   }
}
