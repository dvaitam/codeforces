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
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &grid[i])
   }
   // Prepare walls, apple position, and snake segments
   total := n * m
   walls := make([]bool, total)
   var applePos int
   var length int
   // determine snake length
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           ch := grid[i][j]
           if ch >= '1' && ch <= '9' {
               idx := int(ch - '1') + 1
               if idx > length {
                   length = idx
               }
           }
       }
   }
   // record positions
   segs := make([]int, length)
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           ch := grid[i][j]
           pos := i*m + j
           switch {
           case ch == '#':
               walls[pos] = true
           case ch == '@':
               applePos = pos
           case ch >= '1' && ch <= '9':
               idx := int(ch - '1')
               if idx < length {
                   segs[idx] = pos
               }
           }
       }
   }
   // BFS state: sequence of segment positions (head first)
   start := make([]byte, length)
   for i := 0; i < length; i++ {
       start[i] = byte(segs[i])
   }
   target := byte(applePos)
   visited := make(map[string]bool)
   type node struct {
       s    []byte
       dist int
   }
   queue := []node{{start, 0}}
   visited[string(start)] = true
   // Directions: up, right, down, left
   dirs := [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
   for head := 0; head < len(queue); head++ {
       cur := queue[head]
       s := cur.s
       d := cur.dist
       headIdx := int(s[0])
       if byte(headIdx) == target {
           fmt.Println(d)
           return
       }
       x0 := headIdx / m
       y0 := headIdx % m
       for _, dir := range dirs {
           nx := x0 + dir[0]
           ny := y0 + dir[1]
           if nx < 0 || nx >= n || ny < 0 || ny >= m {
               continue
           }
           np := nx*m + ny
           if walls[np] {
               continue
           }
           // check collision with body except tail
           collide := false
           for i := 0; i < length-1; i++ {
               if int(s[i]) == np {
                   collide = true
                   break
               }
           }
           if collide {
               continue
           }
           // build next state
           ns := make([]byte, length)
           ns[0] = byte(np)
           copy(ns[1:], s[:length-1])
           key := string(ns)
           if !visited[key] {
               visited[key] = true
               queue = append(queue, node{ns, d + 1})
           }
       }
   }
   // unreachable
   fmt.Println(-1)
}
