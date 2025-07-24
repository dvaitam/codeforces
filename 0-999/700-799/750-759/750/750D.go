package main

import (
   "bufio"
   "fmt"
   "os"
)

type state struct {
   x, y  int
   dir   int
   level int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   t := make([]int, n+1)
   sum := 0
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &t[i])
       sum += t[i]
   }
   // grid size
   max := sum + 2
   size := max*2 + 1
   offset := max
   // visited cells
   visited := make([]byte, size*size)
   // visited states
   seen := make(map[uint32]struct{}, 1<<20)
   // directions: 0=right,1=up-right,2=up,3=up-left,4=left,5=down-left,6=down,7=down-right
   dx := [8]int{1, 1, 0, -1, -1, -1, 0, 1}
   dy := [8]int{0, 1, 1, 1, 0, -1, -1, -1}
   // stack for BFS
   stack := make([]state, 0, 1024)
   // start at origin
   sx, sy := offset, offset
   stack = append(stack, state{sx, sy, 2, 1})
   // mark starting cell
   visited[sy*size+sx] = 1
   // process
   for i := 0; i < len(stack); i++ {
       s := stack[i]
       key := packKey(s.x, s.y, s.dir, s.level)
       if _, ok := seen[key]; ok {
           continue
       }
       seen[key] = struct{}{}
       x, y := s.x, s.y
       // step along segment: move t[level]-1 steps (starting cell counted)
       for step := 0; step < t[s.level]-1; step++ {
           x += dx[s.dir]
           y += dy[s.dir]
           visited[y*size+x] = 1
       }
       if s.level < n {
           // branch
           d1 := (s.dir + 1) & 7
           d2 := (s.dir + 7) & 7
           stack = append(stack, state{x, y, d1, s.level + 1})
           stack = append(stack, state{x, y, d2, s.level + 1})
       }
   }
   // count visited cells
   cnt := 0
   for _, v := range visited {
       if v != 0 {
           cnt++
       }
   }
   fmt.Println(cnt)
}

// pack state into 32-bit key: x:9b, y:9b, dir:3b, level:5b -> total 26 bits
func packKey(x, y, dir, level int) uint32 {
   return uint32(x) | uint32(y)<<9 | uint32(dir)<<18 | uint32(level)<<21
}
