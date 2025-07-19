package main

import (
   "bufio"
   "fmt"
   "os"
)

type pt struct{ x, y, d int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, dist int
   if _, err := fmt.Fscan(reader, &n, &m, &dist); err != nil {
       return
   }
   // read newline
   reader.ReadString('\n')
   b := make([][]bool, n+2)
   for i := range b {
       b[i] = make([]bool, m+2)
   }
   var rocks []pt
   for i := 1; i <= n; i++ {
       line, _ := reader.ReadString('\n')
       for j := 1; j <= m && j <= len(line); j++ {
           c := line[j-1]
           if c == 'R' {
               b[i][j] = true
               rocks = append(rocks, pt{i, j, 0})
           } else if c == '.' {
               b[i][j] = true
           }
       }
   }
   k := len(rocks)
   maxCells := 2 * (2*dist + 1) * (2*dist + 1)
   if k > maxCells {
       fmt.Println(-1)
       return
   }
   // coverage stamps
   c1 := make([][]int, n+2)
   c2 := make([][]int, n+2)
   for i := range c1 {
       c1[i] = make([]int, m+2)
       c2[i] = make([]int, m+2)
   }
   iter1, iter2 := 0, 0
   // directions (manhattan)
   dirs := []struct{dx, dy int}{{1,0},{-1,0},{0,1},{0,-1}}
   // explode2 function closure
   var explode2 func(int, int)
   explode2 = func(sx, sy int) {
       iter2++
       // BFS
       q := make([]pt, 0, (2*dist+1)*(2*dist+1))
       if sx < 1 || sy < 1 || sx > n || sy > m || !b[sx][sy] {
           return
       }
       q = append(q, pt{sx, sy, 0})
       c2[sx][sy] = iter2
       for head := 0; head < len(q); head++ {
           p := q[head]
           if p.d == dist {
               continue
           }
           nd := p.d + 1
           for _, dd := range dirs {
               x2, y2 := p.x+dd.dx, p.y+dd.dy
               if x2 >= 1 && y2 >= 1 && x2 <= n && y2 <= m && b[x2][y2] && c2[x2][y2] != iter2 {
                   c2[x2][y2] = iter2
                   q = append(q, pt{x2, y2, nd})
               }
           }
       }
   }
   // explode1 returns ok, second center
   var explode1 func(int, int) (bool, int, int)
   explode1 = func(sx, sy int) (bool, int, int) {
       if sx < 1 || sy < 1 || sx > n || sy > m || !b[sx][sy] {
           return false, -1, -1
       }
       iter1++
       // BFS for first
       q := make([]pt, 0, (2*dist+1)*(2*dist+1))
       q = append(q, pt{sx, sy, 0})
       c1[sx][sy] = iter1
       for head := 0; head < len(q); head++ {
           p := q[head]
           if p.d == dist {
               continue
           }
           nd := p.d + 1
           for _, dd := range dirs {
               x2, y2 := p.x+dd.dx, p.y+dd.dy
               if x2 >= 1 && y2 >= 1 && x2 <= n && y2 <= m && b[x2][y2] && c1[x2][y2] != iter1 {
                   c1[x2][y2] = iter1
                   q = append(q, pt{x2, y2, nd})
               }
           }
       }
       // determine bounds
       minx, maxx := 1, n
       miny, maxy := 1, m
       allDead := true
       for _, r := range rocks {
           if c1[r.x][r.y] != iter1 {
               allDead = false
               if r.x-dist > minx {
                   minx = r.x - dist
               }
               if r.x+dist < maxx {
                   maxx = r.x + dist
               }
               if r.y-dist > miny {
                   miny = r.y - dist
               }
               if r.y+dist < maxy {
                   maxy = r.y + dist
               }
           }
       }
       if allDead {
           return true, -1, -1
       }
       // try second
       for x2 := minx; x2 <= maxx; x2++ {
           for y2 := miny; y2 <= maxy; y2++ {
               if x2 < 1 || y2 < 1 || x2 > n || y2 > m || !b[x2][y2] {
                   continue
               }
               explode2(x2, y2)
               ok := true
               for _, r := range rocks {
                   if c1[r.x][r.y] != iter1 && c2[r.x][r.y] != iter2 {
                       ok = false
                       break
                   }
               }
               if ok {
                   return true, x2, y2
               }
           }
       }
       return false, -1, -1
   }
   // search first
   fx0, fy0 := rocks[0].x, rocks[0].y
   for i := fx0 - dist; i <= fx0 + dist; i++ {
       for j := fy0 - dist; j <= fy0 + dist; j++ {
           ok, sx2, sy2 := explode1(i, j)
           if !ok {
               continue
           }
           // adjust second
           if i == sx2 && j == sy2 {
               sx2, sy2 = -1, -1
           }
           if sx2 == -1 {
               // pick any other
               for x := 1; x <= n; x++ {
                   for y := 1; y <= m; y++ {
                       if b[x][y] && (x != i || y != j) {
                           sx2, sy2 = x, y
                           break
                       }
                   }
                   if sx2 != -1 {
                       break
                   }
               }
           }
           fmt.Printf("%d %d %d %d\n", i, j, sx2, sy2)
           return
       }
   }
   fmt.Println(-1)
}
