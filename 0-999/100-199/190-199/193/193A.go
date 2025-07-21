package main

import (
   "bufio"
   "fmt"
   "os"
)

type point struct { x, y int }

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(in, &line)
       grid[i] = []byte(line)
   }
   // count painted cells
   total := 0
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == '#' {
               total++
           }
       }
   }
   // impossible if 1 or 2 cells
   if total <= 2 {
       fmt.Println(-1)
       return
   }
   // directions for adjacency
   dirs := []point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
   // check for articulation point
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] != '#' {
               continue
           }
           // remove this cell
           grid[i][j] = '.'
           // find a start cell for BFS
           var sx, sy int
           found := false
           for x := 0; x < n && !found; x++ {
               for y := 0; y < m; y++ {
                   if grid[x][y] == '#' {
                       sx, sy = x, y
                       found = true
                       break
                   }
               }
           }
           // perform BFS if there is any remaining
           reach := 0
           if found {
               vis := make([][]bool, n)
               for ii := range vis {
                   vis[ii] = make([]bool, m)
               }
               queue := make([]point, 0, total)
               queue = append(queue, point{sx, sy})
               vis[sx][sy] = true
               reach = 1
               for qi := 0; qi < len(queue); qi++ {
                   u := queue[qi]
                   for _, d := range dirs {
                       nx, ny := u.x+d.x, u.y+d.y
                       if nx >= 0 && nx < n && ny >= 0 && ny < m && !vis[nx][ny] && grid[nx][ny] == '#' {
                           vis[nx][ny] = true
                           reach++
                           queue = append(queue, point{nx, ny})
                       }
                   }
               }
           }
           // restore cell
           grid[i][j] = '#'
           // if disconnected
           if reach < total-1 {
               fmt.Println(1)
               return
           }
       }
   }
   // no single removal disconnects, so need at least 2
   fmt.Println(2)
}
