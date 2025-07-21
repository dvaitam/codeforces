package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var h, w int
   if _, err := fmt.Fscan(reader, &h, &w); err != nil {
       return
   }
   grid := make([][]byte, h)
   for i := 0; i < h; i++ {
       grid[i] = make([]byte, w)
       for j := 0; j < w; j++ {
           var v int
           fmt.Fscan(reader, &v)
           if v != 0 {
               grid[i][j] = 1
           }
       }
   }
   visited := make([][]bool, h)
   for i := range visited {
       visited[i] = make([]bool, w)
   }
   var rayCounts []int
   // directions for 8-connected
   dirs := [8][2]int{{1,0},{-1,0},{0,1},{0,-1},{1,1},{1,-1},{-1,1},{-1,-1}}
   for i := 0; i < h; i++ {
       for j := 0; j < w; j++ {
           if grid[i][j] == 1 && !visited[i][j] {
               // BFS for component
               var q [][2]int
               q = append(q, [2]int{i, j})
               visited[i][j] = true
               comp := make([][2]int, 0, 256)
               for qi := 0; qi < len(q); qi++ {
                   x, y := q[qi][0], q[qi][1]
                   comp = append(comp, [2]int{x, y})
                   for _, d := range dirs {
                       nx, ny := x+d[0], y+d[1]
                       if nx >= 0 && nx < h && ny >= 0 && ny < w && grid[nx][ny] == 1 && !visited[nx][ny] {
                           visited[nx][ny] = true
                           q = append(q, [2]int{nx, ny})
                       }
                   }
               }
               // compute centroid
               var sx, sy float64
               for _, p := range comp {
                   sx += float64(p[0])
                   sy += float64(p[1])
               }
               cnt := float64(len(comp))
               cx := sx / cnt
               cy := sy / cnt
               // histogram of rounded distances
               hist := make(map[int]int)
               for _, p := range comp {
                   dx := float64(p[0]) - cx
                   dy := float64(p[1]) - cy
                   d := math.Hypot(dx, dy)
                   r := int(d + 0.5)
                   hist[r]++
               }
               // find mode radius
               r0, maxc := 0, 0
               for r, c := range hist {
                   if c > maxc {
                       maxc = c
                       r0 = r
                   }
               }
               // collect ray pixels: those with rounded dist > r0 + 2
               cutoff := r0 + 2
               var rays [] [2]int
               for _, p := range comp {
                   dx := float64(p[0]) - cx
                   dy := float64(p[1]) - cy
                   d := math.Hypot(dx, dy)
                   if int(d+0.5) > cutoff {
                       rays = append(rays, p)
                   }
               }
               if len(rays) == 0 {
                   rayCounts = append(rayCounts, 0)
                   continue
               }
               // bounding box for rays
               minx, miny := rays[0][0], rays[0][1]
               maxx, maxy := minx, miny
               for _, p := range rays {
                   if p[0] < minx {
                       minx = p[0]
                   }
                   if p[0] > maxx {
                       maxx = p[0]
                   }
                   if p[1] < miny {
                       miny = p[1]
                   }
                   if p[1] > maxy {
                       maxy = p[1]
                   }
               }
               hh := maxx - minx + 1
               ww := maxy - miny + 1
               mask := make([][]bool, hh)
               vis2 := make([][]bool, hh)
               for ii := 0; ii < hh; ii++ {
                   mask[ii] = make([]bool, ww)
                   vis2[ii] = make([]bool, ww)
               }
               for _, p := range rays {
                   mask[p[0]-minx][p[1]-miny] = true
               }
               // count ray connected components
               var rc int
               for ii := 0; ii < hh; ii++ {
                   for jj := 0; jj < ww; jj++ {
                       if mask[ii][jj] && !vis2[ii][jj] {
                           rc++
                           // BFS in mask
                           var q2 [][2]int
                           q2 = append(q2, [2]int{ii, jj})
                           vis2[ii][jj] = true
                           for qi := 0; qi < len(q2); qi++ {
                               x, y := q2[qi][0], q2[qi][1]
                               for _, d := range dirs {
                                   nx, ny := x+d[0], y+d[1]
                                   if nx >= 0 && nx < hh && ny >= 0 && ny < ww && mask[nx][ny] && !vis2[nx][ny] {
                                       vis2[nx][ny] = true
                                       q2 = append(q2, [2]int{nx, ny})
                                   }
                               }
                           }
                       }
                   }
               }
               rayCounts = append(rayCounts, rc)
           }
       }
   }
   sort.Ints(rayCounts)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprintln(out, len(rayCounts))
   for i, v := range rayCounts {
       if i > 0 {
           fmt.Fprint(out, " ")
       }
       fmt.Fprint(out, v)
   }
   fmt.Fprintln(out)
}
