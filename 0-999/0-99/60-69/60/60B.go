package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var k, n, m int
   fmt.Fscan(reader, &k, &n, &m)
   // read grid layers
   grid := make([][][]byte, k)
   for z := 0; z < k; z++ {
       grid[z] = make([][]byte, n)
       for i := 0; i < n; i++ {
           var line string
           fmt.Fscan(reader, &line)
           grid[z][i] = []byte(line)
       }
   }
   // read tap coordinates (1-based)
   var x, y int
   fmt.Fscan(reader, &x, &y)
   sx, sy := x-1, y-1
   // BFS to count reachable empty cells
   visited := make([][][]bool, k)
   for z := 0; z < k; z++ {
       visited[z] = make([][]bool, n)
       for i := 0; i < n; i++ {
           visited[z][i] = make([]bool, m)
       }
   }
   type node struct{ z, x, y int }
   queue := make([]node, 0, k*n*m)
   // start at top layer z=0
   visited[0][sx][sy] = true
   queue = append(queue, node{0, sx, sy})
   // 6-directional moves: dz, dx, dy
   dirs := []node{{1, 0, 0}, {-1, 0, 0}, {0, 1, 0}, {0, -1, 0}, {0, 0, 1}, {0, 0, -1}}
   count := 0
   for head := 0; head < len(queue); head++ {
       u := queue[head]
       count++
       for _, d := range dirs {
           nz, nx, ny := u.z+d.z, u.x+d.x, u.y+d.y
           if nz >= 0 && nz < k && nx >= 0 && nx < n && ny >= 0 && ny < m {
               if !visited[nz][nx][ny] && grid[nz][nx][ny] == '.' {
                   visited[nz][nx][ny] = true
                   queue = append(queue, node{nz, nx, ny})
               }
           }
       }
   }
   fmt.Println(count)
}
