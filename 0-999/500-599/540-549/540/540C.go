package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       grid[i] = []byte(line)
   }
   var r1, c1, r2, c2 int
   fmt.Fscan(reader, &r1, &c1)
   fmt.Fscan(reader, &r2, &c2)
   r1--
   c1--
   r2--
   c2--
   // BFS from start
   visited := make([][]bool, n)
   for i := range visited {
       visited[i] = make([]bool, m)
   }
   type pt struct{ x, y int }
   queue := make([]pt, 0, n*m)
   visited[r1][c1] = true
   queue = append(queue, pt{r1, c1})
   dirs := []pt{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
   for i := 0; i < len(queue); i++ {
       p := queue[i]
       for _, d := range dirs {
           nx, ny := p.x+d.x, p.y+d.y
           if nx < 0 || nx >= n || ny < 0 || ny >= m {
               continue
           }
           if visited[nx][ny] {
               continue
           }
           // can move if intact or it's the target
           if grid[nx][ny] == '.' || (nx == r2 && ny == c2) {
               visited[nx][ny] = true
               queue = append(queue, pt{nx, ny})
           }
       }
   }
   if !visited[r2][c2] {
       fmt.Println("NO")
       return
   }
   // count accessible neighbors of target
   cnt := 0
   for _, d := range dirs {
       nx, ny := r2+d.x, c2+d.y
       if nx < 0 || nx >= n || ny < 0 || ny >= m {
           continue
       }
       if grid[nx][ny] == '.' || (nx == r1 && ny == c1) {
           cnt++
       }
   }
   // decide
   if r1 == r2 && c1 == c2 {
       if cnt >= 1 {
           fmt.Println("YES")
       } else {
           fmt.Println("NO")
       }
   } else if grid[r2][c2] == 'X' {
       fmt.Println("YES")
   } else {
       if cnt >= 2 {
           fmt.Println("YES")
       } else {
           fmt.Println("NO")
       }
   }
}
