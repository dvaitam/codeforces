package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       grid[i] = []byte(s)
   }
   // mark zeros reachable from border (outside zeros)
   outside := make([][]bool, n)
   for i := 0; i < n; i++ {
       outside[i] = make([]bool, m)
   }
   dq := make([][2]int, 0, n*m)
   // enqueue border zeros
   for i := 0; i < n; i++ {
       for _, j := range []int{0, m - 1} {
           if grid[i][j] == '0' && !outside[i][j] {
               outside[i][j] = true
               dq = append(dq, [2]int{i, j})
           }
       }
   }
   for j := 0; j < m; j++ {
       for _, i := range []int{0, n - 1} {
           if grid[i][j] == '0' && !outside[i][j] {
               outside[i][j] = true
               dq = append(dq, [2]int{i, j})
           }
       }
   }
   // BFS for outside zeros
   for bi := 0; bi < len(dq); bi++ {
       r, c := dq[bi][0], dq[bi][1]
       for _, d := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
           nr, nc := r+d[0], c+d[1]
           if nr >= 0 && nr < n && nc >= 0 && nc < m && grid[nr][nc] == '0' && !outside[nr][nc] {
               outside[nr][nc] = true
               dq = append(dq, [2]int{nr, nc})
           }
       }
   }
   // Traverse faces in 1s graph to find finite faces (cool cycles)
   // adjacency directions: up, right, down, left
   dr := [4]int{-1, 0, 1, 0}
   dc := [4]int{0, 1, 0, -1}
   // visited half-edges
   visitedEdge := make([][][]bool, n)
   for i := 0; i < n; i++ {
       visitedEdge[i] = make([][]bool, m)
       for j := 0; j < m; j++ {
           visitedEdge[i][j] = make([]bool, 4)
       }
   }
   best := 0
   for r := 0; r < n; r++ {
       for c := 0; c < m; c++ {
           if grid[r][c] != '1' {
               continue
           }
           for d := 0; d < 4; d++ {
               nr, nc := r+dr[d], c+dc[d]
               // half-edge exists
               if nr < 0 || nr >= n || nc < 0 || nc >= m || grid[nr][nc] != '1' {
                   continue
               }
               if visitedEdge[r][c][d] {
                   continue
               }
               // traverse face with face on right
               cr, cc, dir := r, c, d
               startR, startC, startDir := r, c, d
               length := 0
               infinite := false
               for {
                   // mark half-edge visited
                   visitedEdge[cr][cc][dir] = true
                   // move to next vertex
                   nr, nc := cr+dr[dir], cc+dc[dir]
                   if nr < 0 || nr >= n || nc < 0 || nc >= m || grid[nr][nc] != '1' {
                       infinite = true
                       break
                   }
                   // at (nr,nc), pick next direction
                   found := false
                   // right, straight, left, back
                   for _, turn := range []int{1, 0, 3, 2} {
                       nd := (dir + turn) & 3
                       wr, wc := nr+dr[nd], nc+dc[nd]
                       if wr >= 0 && wr < n && wc >= 0 && wc < m && grid[wr][wc] == '1' {
                           cr, cc, dir = nr, nc, nd
                           found = true
                           break
                       }
                   }
                   if !found {
                       infinite = true
                       break
                   }
                   length++
                   if cr == startR && cc == startC && dir == startDir {
                       break
                   }
               }
               if !infinite && length > 0 && length > best {
                   best = length
               }
           }
       }
   }
   fmt.Fprint(writer, best)
}
