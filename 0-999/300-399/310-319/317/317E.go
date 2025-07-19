package main

import (
   "bufio"
   "fmt"
   "os"
)

// Node represents a point on the grid
type Node struct { x, y int }

var (
   // Directions: right, up, down, left
   D = [4][2]int{{1, 0}, {0, 1}, {0, -1}, {-1, 0}}
   DIRCH = "RUDL"
   s, t Node
   n    int
   // obstacle grid
   c = [411][411]bool
   // BFS visitation
   visited = [411][411]bool
   // extremal obstacle points
   P = [4]Node
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   // read input
   fmt.Fscan(in, &s.x, &s.y, &t.x, &t.y, &n)
   if n == 0 {
       out.WriteString("-1\n")
       return
   }
   // shift coordinates
   s.x += 200; s.y += 200
   t.x += 200; t.y += 200
   // initialize extremal points
   P[0].x = -1000000000
   P[1].y = -1000000000
   P[2].y = 1000000000
   P[3].x = 1000000000
   // read obstacles
   for i := 0; i < n; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       x += 200; y += 200
       c[x][y] = true
       // update extremal
       if x < P[3].x || (x == P[3].x && y < P[3].y) {
           P[3].x, P[3].y = x, y
       }
       if x > P[0].x || (x == P[0].x && y > P[0].y) {
           P[0].x, P[0].y = x, y
       }
       if y < P[2].y || (y == P[2].y && x < P[2].x) {
           P[2].y, P[2].x = y, x
       }
       if y > P[1].y || (y == P[1].y && x > P[1].x) {
           P[1].y, P[1].x = y, x
       }
   }
   // find initial path
   path := bfs()
   if path == nil {
       out.WriteString("-1\n")
       return
   }
   // follow path, appending successful moves
   p := make([]int, len(path))
   copy(p, path)
   m := len(path)
   for i := 0; i < m; i++ {
       if moveDir(p[i], out) {
           p = append(p, p[i])
       }
       if s == t {
           out.WriteByte('\n')
           return
       }
       if !inGrid(s) {
           break
       }
   }
   // adjust vertical
   if s.y != t.y {
       var d int
       if s.y < t.y {
           d = 2
       } else {
           d = 1
       }
       ud(d, out)
   }
   // adjust horizontal
   if s.x != t.x {
       var d int
       if s.x < t.x {
           d = 3
       } else {
           d = 0
       }
       lr(d, out)
   }
   out.WriteByte('\n')
}

// bfs finds a path of directions from s to t ignoring obstacles
func bfs() []int {
   // clear visited
   for i := 0; i <= 410; i++ {
       for j := 0; j <= 410; j++ {
           visited[i][j] = false
       }
   }
   type QItem struct { x, y, pre, fr int }
   q := make([]QItem, 0, 411*411)
   q = append(q, QItem{s.x, s.y, -1, -1})
   visited[s.x][s.y] = true
   for head := 0; head < len(q); head++ {
       u := q[head]
       for i := 0; i < 4; i++ {
           nx := u.x + D[i][0]
           ny := u.y + D[i][1]
           if nx >= 0 && nx <= 410 && ny >= 0 && ny <= 410 && !visited[nx][ny] {
               visited[nx][ny] = true
               q = append(q, QItem{nx, ny, head, i})
               if nx == t.x && ny == t.y {
                   // backtrack
                   path := []int{}
                   idx := len(q) - 1
                   for q[idx].pre != -1 {
                       path = append(path, q[idx].fr)
                       idx = q[idx].pre
                   }
                   // reverse
                   for l, r := 0, len(path)-1; l < r; l, r = l+1, r-1 {
                       path[l], path[r] = path[r], path[l]
                   }
                   return path
               }
           }
       }
   }
   return nil
}

// moveDir moves s and t by direction i, writes the move, and returns true if no obstacle hit at t
func moveDir(i int, out *bufio.Writer) bool {
   s.x += D[i][0]
   s.y += D[i][1]
   t.x += D[i][0]
   t.y += D[i][1]
   out.WriteByte(DIRCH[i])
   if t.x >= 0 && t.x <= 410 && t.y >= 0 && t.y <= 410 && c[t.x][t.y] {
       t.x -= D[i][0]
       t.y -= D[i][1]
       return false
   }
   return true
}

// inGrid checks if a node is within [0,410]
func inGrid(u Node) bool {
   return u.x >= 0 && u.x <= 410 && u.y >= 0 && u.y <= 410
}

// getout moves both points outside the grid bounds
func getout() {
   for (s.x >= 0 && s.x <= 410) || (t.x >= 0 && t.x <= 410) {
       moveDir(0, bufio.NewWriter(os.Stdout))
   }
   for (s.y >= 0 && s.y <= 410) || (t.y >= 0 && t.y <= 410) {
       moveDir(1, bufio.NewWriter(os.Stdout))
   }
}

// ud adjusts vertical alignment using direction d (1=up,2=down)
func ud(d int, out *bufio.Writer) {
   getout()
   for (d == 2 && t.y >= 0) || (d == 1 && t.y <= 410) {
       moveDir(d, out)
   }
   for t.x != P[d].x {
       if t.x < P[d].x {
           moveDir(0, out)
       } else {
           moveDir(3, out)
       }
   }
   for t.y != s.y {
       moveDir(3-d, out)
   }
}

// lr adjusts horizontal alignment using direction d (0=right,3=left)
func lr(d int, out *bufio.Writer) {
   getout()
   for (d == 3 && t.x >= 0) || (d == 0 && t.x <= 410) {
       moveDir(d, out)
   }
   for t.y != P[d].y {
       if t.y < P[d].y {
           moveDir(1, out)
       } else {
           moveDir(2, out)
       }
   }
   for t.x != s.x {
       moveDir(3-d, out)
   }
}
