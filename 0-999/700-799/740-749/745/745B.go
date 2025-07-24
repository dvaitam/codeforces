package main

import (
   "bufio"
   "fmt"
   "os"
)

// pt represents a cell of the shape relative to its bounding box
type pt struct{ x, y int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &grid[i])
   }
   var pts []pt
   minx, miny := n, m
   maxx, maxy := -1, -1
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == 'X' {
               pts = append(pts, pt{i, j})
               if i < minx {
                   minx = i
               }
               if i > maxx {
                   maxx = i
               }
               if j < miny {
                   miny = j
               }
               if j > maxy {
                   maxy = j
               }
           }
       }
   }
   // relative positions
   h0 := maxx - minx + 1
   w0 := maxy - miny + 1
   k := len(pts)
   // boolean grid of shape
   grid0 := make([]bool, h0*w0)
   rel := make([]pt, 0, k)
   for _, p := range pts {
       rx := p.x - minx
       ry := p.y - miny
       grid0[rx*w0+ry] = true
       rel = append(rel, pt{rx, ry})
   }
   t := 2 * k
   // try all rectangle sizes H x W = 2*k
   for H := 1; H*H <= t; H++ {
       if t%H != 0 {
           continue
       }
       W := t / H
       if tryRect(H, W, h0, w0, k, grid0, rel) {
           fmt.Println("YES")
           return
       }
       if H != W {
           // try swapped dims
           if tryRect(W, H, h0, w0, k, grid0, rel) {
               fmt.Println("YES")
               return
           }
       }
   }
   fmt.Println("NO")
}

// tryRect checks if rectangle HxW can be formed with two copies of shape
func tryRect(H, W, h0, w0, k int, grid0 []bool, rel []pt) bool {
   if H < h0 || W < w0 {
       return false
   }
   // dx options
   var dxs []int
   if H == h0 {
       dxs = []int{0}
   } else {
       dxs = []int{0, H - h0}
   }
   // dy options
   var dys []int
   if W == w0 {
       dys = []int{0}
   } else {
       dys = []int{0, W - w0}
   }
   // area match automatically H*W == 2*k
   for _, dx := range dxs {
       for _, dy := range dys {
           if dx == 0 && dy == 0 {
               continue
           }
           // check overlap
           // if shift moves shape entirely out of original bounds, no overlap
           if dx >= h0 || dy >= w0 {
               return true
           }
           ok := true
           for _, p := range rel {
               i2 := p.x + dx
               j2 := p.y + dy
               if i2 >= 0 && i2 < h0 && j2 >= 0 && j2 < w0 {
                   if grid0[i2*w0+j2] {
                       ok = false
                       break
                   }
               }
           }
           if ok {
               return true
           }
       }
   }
   return false
}
