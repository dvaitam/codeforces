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
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       line, err := reader.ReadString('\n')
       if err != nil && len(line) == 0 {
           return
       }
       if len(line) < m {
           // read until full line
           more, _ := reader.ReadString('\n')
           line += more
       }
       grid[i] = []byte(line[:m])
   }
   size := n * m
   dist := make([]int, size)
   queue := make([]int, 0, size)
   // initial nodes: arrows pointing to blocked
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           u := grid[i][j]
           if u == '<' || u == '>' || u == '^' || u == 'v' {
               var ni, nj int
               switch u {
               case '<': ni, nj = i, j-1
               case '>': ni, nj = i, j+1
               case '^': ni, nj = i-1, j
               case 'v': ni, nj = i+1, j
               }
               if grid[ni][nj] == '#' {
                   idx := i*m + j
                   dist[idx] = 1
                   queue = append(queue, idx)
               }
           }
       }
   }
   // BFS in reverse graph
   head := 0
   for head < len(queue) {
       v := queue[head]
       head++
       vi := v / m
       vj := v % m
       // check neighbors u that point to v
       // up neighbor (i-1,j) with 'v'
       if vi > 0 {
           ui, uj := vi-1, vj
           if grid[ui][uj] == 'v' {
               u := ui*m + uj
               if dist[u] == 0 {
                   dist[u] = dist[v] + 1
                   queue = append(queue, u)
               }
           }
       }
       // down neighbor (i+1,j) with '^'
       if vi+1 < n {
           ui, uj := vi+1, vj
           if grid[ui][uj] == '^' {
               u := ui*m + uj
               if dist[u] == 0 {
                   dist[u] = dist[v] + 1
                   queue = append(queue, u)
               }
           }
       }
       // left neighbor (i,j-1) with '>'
       if vj > 0 {
           ui, uj := vi, vj-1
           if grid[ui][uj] == '>' {
               u := ui*m + uj
               if dist[u] == 0 {
                   dist[u] = dist[v] + 1
                   queue = append(queue, u)
               }
           }
       }
       // right neighbor (i,j+1) with '<'
       if vj+1 < m {
           ui, uj := vi, vj+1
           if grid[ui][uj] == '<' {
               u := ui*m + uj
               if dist[u] == 0 {
                   dist[u] = dist[v] + 1
                   queue = append(queue, u)
               }
           }
       }
   }
   // detect cycles and find top two
   max1, max2 := 0, 0
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           c := grid[i][j]
           if c == '<' || c == '>' || c == '^' || c == 'v' {
               d := dist[i*m+j]
               if d == 0 {
                   fmt.Println(-1)
                   return
               }
               if d > max1 {
                   max2 = max1
                   max1 = d
               } else if d > max2 {
                   max2 = d
               }
           }
       }
   }
   fmt.Println(max1 + max2)
}
