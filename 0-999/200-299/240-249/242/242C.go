package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var x0, y0, x1, y1 int
   fmt.Fscan(reader, &x0, &y0, &x1, &y1)
   var n int
   fmt.Fscan(reader, &n)
   allowed := make(map[uint64]bool, n)
   for i := 0; i < n; i++ {
       var r, a, b int
       fmt.Fscan(reader, &r, &a, &b)
       for c := a; c <= b; c++ {
           key := (uint64(r) << 32) | uint64(c)
           allowed[key] = true
       }
   }
   start := (uint64(x0) << 32) | uint64(y0)
   goal := (uint64(x1) << 32) | uint64(y1)
   dirs := [8][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
   visited := make(map[uint64]int, len(allowed))
   queue := make([]uint64, 0, len(allowed))
   visited[start] = 0
   queue = append(queue, start)
   for head := 0; head < len(queue); head++ {
       cur := queue[head]
       dist := visited[cur]
       if cur == goal {
           fmt.Println(dist)
           return
       }
       x := int(cur >> 32)
       y := int(cur & 0xFFFFFFFF)
       for _, d := range dirs {
           nx := x + d[0]
           ny := y + d[1]
           nk := (uint64(nx) << 32) | uint64(ny)
           if !allowed[nk] {
               continue
           }
           if _, seen := visited[nk]; !seen {
               visited[nk] = dist + 1
               queue = append(queue, nk)
           }
       }
   }
   fmt.Println(-1)
}
