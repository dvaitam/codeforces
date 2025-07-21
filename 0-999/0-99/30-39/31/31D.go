package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var W, H, n int
   if _, err := fmt.Fscan(in, &W, &H, &n); err != nil {
       return
   }
   type seg struct{ x1, y1, x2, y2 int }
   segs := make([]seg, n)
   xs := make([]int, 0, 2*n+2)
   ys := make([]int, 0, 2*n+2)
   xs = append(xs, 0, W)
   ys = append(ys, 0, H)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &segs[i].x1, &segs[i].y1, &segs[i].x2, &segs[i].y2)
       xs = append(xs, segs[i].x1, segs[i].x2)
       ys = append(ys, segs[i].y1, segs[i].y2)
   }
   sort.Ints(xs)
   sort.Ints(ys)
   ux := xs[:1]
   for i := 1; i < len(xs); i++ {
       if xs[i] != xs[i-1] {
           ux = append(ux, xs[i])
       }
   }
   uy := ys[:1]
   for i := 1; i < len(ys); i++ {
       if ys[i] != ys[i-1] {
           uy = append(uy, ys[i])
       }
   }
   nx, ny := len(ux), len(uy)
   // barriers between cells
   br := make([][]bool, nx-1)
   bt := make([][]bool, nx-1)
   for i := 0; i < nx-1; i++ {
       br[i] = make([]bool, ny-1) // right barrier from (i,j) to (i+1,j)
       bt[i] = make([]bool, ny-1) // top barrier from (i,j) to (i,j+1)
   }
   // map coord to index
   xIndex := make(map[int]int, nx)
   yIndex := make(map[int]int, ny)
   for i, v := range ux {
       xIndex[v] = i
   }
   for j, v := range uy {
       yIndex[v] = j
   }
   for _, s := range segs {
       if s.x1 == s.x2 {
           // vertical segment at x = s.x1
           x := s.x1
           if x <= 0 || x >= W {
               continue
           }
           xi := xIndex[x]
           y1i := yIndex[s.y1]
           y2i := yIndex[s.y2]
           if y1i > y2i {
               y1i, y2i = y2i, y1i
           }
           // barrier between cells at xi-1 and xi
           if xi-1 >= 0 {
               for j := y1i; j < y2i; j++ {
                   br[xi-1][j] = true
               }
           }
       } else if s.y1 == s.y2 {
           // horizontal segment at y = s.y1
           y := s.y1
           if y <= 0 || y >= H {
               continue
           }
           yi := yIndex[y]
           x1i := xIndex[s.x1]
           x2i := xIndex[s.x2]
           if x1i > x2i {
               x1i, x2i = x2i, x1i
           }
           // barrier between cells at yi-1 and yi
           if yi-1 >= 0 {
               for i := x1i; i < x2i; i++ {
                   bt[i][yi-1] = true
               }
           }
       }
   }
   visited := make([][]bool, nx-1)
   for i := range visited {
       visited[i] = make([]bool, ny-1)
   }
   var areas []int
   // neighbor offsets: right, left, up, down
   for i := 0; i < nx-1; i++ {
       for j := 0; j < ny-1; j++ {
           if visited[i][j] {
               continue
           }
           // BFS
           area := 0
           stack := [][2]int{{i, j}}
           visited[i][j] = true
           for len(stack) > 0 {
               ci, cj := stack[len(stack)-1][0], stack[len(stack)-1][1]
               stack = stack[:len(stack)-1]
               // add cell area
               dx := ux[ci+1] - ux[ci]
               dy := uy[cj+1] - uy[cj]
               area += dx * dy
               // right
               if ci+1 < nx-1 && !visited[ci+1][cj] && !br[ci][cj] {
                   visited[ci+1][cj] = true
                   stack = append(stack, [2]int{ci + 1, cj})
               }
               // left
               if ci-1 >= 0 && !visited[ci-1][cj] && !br[ci-1][cj] {
                   visited[ci-1][cj] = true
                   stack = append(stack, [2]int{ci - 1, cj})
               }
               // up
               if cj+1 < ny-1 && !visited[ci][cj+1] && !bt[ci][cj] {
                   visited[ci][cj+1] = true
                   stack = append(stack, [2]int{ci, cj + 1})
               }
               // down
               if cj-1 >= 0 && !visited[ci][cj-1] && !bt[ci][cj-1] {
                   visited[ci][cj-1] = true
                   stack = append(stack, [2]int{ci, cj - 1})
               }
           }
           areas = append(areas, area)
       }
   }
   sort.Ints(areas)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i, a := range areas {
       if i > 0 {
           fmt.Fprint(out, " ")
       }
       fmt.Fprint(out, a)
   }
   fmt.Fprintln(out)
}
