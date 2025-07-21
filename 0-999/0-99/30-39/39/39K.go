package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   // read grid and find object bounding boxes
   r1 := make([]int, k)
   r2 := make([]int, k)
   c1 := make([]int, k)
   c2 := make([]int, k)
   for i := 0; i < k; i++ {
       r1[i] = n + 1
       c1[i] = m + 1
   }
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &grid[i])
   }
   // map of object id per cell
   id := make([][]int, n)
   for i := range id {
       id[i] = make([]int, m)
   }
   // assign object ids by scanning connected components
   // Since k objects are rectangles not touching, we can assign by flood fill
   var dirs = [4][2]int{{1,0},{-1,0},{0,1},{0,-1}}
   obj := 0
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == '*' && id[i][j] == 0 {
               obj++
               // flood fill
               stack := [][2]int{{i,j}}
               id[i][j] = obj
               for len(stack) > 0 {
                   x, y := stack[len(stack)-1][0], stack[len(stack)-1][1]
                   stack = stack[:len(stack)-1]
                   // update bounds
                   idx := obj - 1
                   if x+1 < r1[idx] { r1[idx] = x+1 }
                   if x+1 > r2[idx] { r2[idx] = x+1 }
                   if y+1 < c1[idx] { c1[idx] = y+1 }
                   if y+1 > c2[idx] { c2[idx] = y+1 }
                   for _, d := range dirs {
                       nx, ny := x+d[0], y+d[1]
                       if nx>=0 && nx<n && ny>=0 && ny<m && grid[nx][ny]=='*' && id[nx][ny]==0 {
                           id[nx][ny] = obj
                           stack = append(stack, [2]int{nx,ny})
                       }
                   }
               }
           }
       }
   }
   // compute answer by inclusion-exclusion up to size 3
   var ans int64
   // helper to process a subset
   process := func(idx []int, sign int64) {
       // bounding box of S
       rl, rh := n+1, 0
       cl, ch := m+1, 0
       inS := make([]bool, k)
       for _, i := range idx {
           inS[i] = true
           if r1[i] < rl { rl = r1[i] }
           if r2[i] > rh { rh = r2[i] }
           if c1[i] < cl { cl = c1[i] }
           if c2[i] > ch { ch = c2[i] }
       }
       // vertical limits for c-overlap objs
       T_low, B_high := 1, n
       // horizontal blockers
       maxC2 := 0
       minC1 := m + 1
       for j := 0; j < k; j++ {
           if inS[j] { continue }
           // check column relation
           if c1[j] <= ch && c2[j] >= cl {
               // c-overlap
               if r2[j] < rl {
                   // above
                   if r2[j]+1 > T_low { T_low = r2[j] + 1 }
               } else if r1[j] > rh {
                   // below
                   if r1[j]-1 < B_high { B_high = r1[j] - 1 }
               } else {
                   // intersects S-range: invalid
                   return
               }
           } else if c2[j] < cl && r1[j] <= rh && r2[j] >= rl {
               // left overlap
               if c2[j] > maxC2 { maxC2 = c2[j] }
           } else if c1[j] > ch && r1[j] <= rh && r2[j] >= rl {
               // right overlap
               if c1[j] < minC1 { minC1 = c1[j] }
           }
       }
       if T_low > rl || B_high < rh { return }
       horiz := int64(cl - maxC2) * int64(minC1 - ch)
       if horiz <= 0 { return }
       vert := int64(rl - T_low + 1) * int64(B_high - rh + 1)
       ans += sign * horiz * vert
   }
   // size 1
   for i := 0; i < k; i++ {
       process([]int{i}, 1)
   }
   // size 2
   for i := 0; i < k; i++ {
       for j := i + 1; j < k; j++ {
           process([]int{i, j}, -1)
       }
   }
   // size 3
   for i := 0; i < k; i++ {
       for j := i + 1; j < k; j++ {
           for l := j + 1; l < k; l++ {
               process([]int{i, j, l}, 1)
           }
       }
   }
   fmt.Println(ans)
}
