package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   words := make([]string, 6)
   for i := 0; i < 6; i++ {
       if _, err := fmt.Fscan(in, &words[i]); err != nil {
           return
       }
   }
   // Try partitions: choose 3 words as horizontals
   var best []string
   used := make([]bool, 6)
   // generate masks for horizontals
   var idxH [3]int
   var idxV [3]int
   // combination of 3 indices out of 6
   var combs [][]int
   var dfsComb func(start, cnt int, cur []int)
   dfsComb = func(start, cnt int, cur []int) {
       if cnt == 0 {
           a := make([]int, len(cur))
           copy(a, cur)
           combs = append(combs, a)
           return
       }
       for i := start; i <= 6-cnt; i++ {
           dfsComb(i+1, cnt-1, append(cur, i))
       }
   }
   dfsComb(0, 3, []int{})
   // permutations helper
   var permute func(a []int, l int, fn func([]int))
   permute = func(a []int, l int, fn func([]int)) {
       if l == len(a)-1 {
           b := make([]int, len(a)); copy(b, a)
           fn(b)
       } else {
           for i := l; i < len(a); i++ {
               a[l], a[i] = a[i], a[l]
               permute(a, l+1, fn)
               a[l], a[i] = a[i], a[l]
           }
       }
   }
   // iterate partitions
   for _, comb := range combs {
       // mark horizontals
       hset := make(map[int]bool)
       for _, i := range comb {
           hset[i] = true
       }
       // build verticals indices
       vi := 0
       for i := 0; i < 6; i++ {
           if !hset[i] {
               idxV[vi] = i
               vi++
           }
       }
       // horizontal indices
       for i := 0; i < 3; i++ {
           idxH[i] = comb[i]
       }
       // permute horizontals
       permute(idxH[:], 0, func(hp []int) {
           H1, H2, H3 := words[hp[0]], words[hp[1]], words[hp[2]]
           // permute verticals
           permute(idxV[:], 0, func(vp []int) {
               V1, V2, V3 := words[vp[0]], words[vp[1]], words[vp[2]]
               // all lengths at least 3
               if len(H1) < 3 || len(H2) < 3 || len(H3) < 3 || len(V1) < 3 || len(V2) < 3 || len(V3) < 3 {
                   return
               }
               // search intersections
               // k1: H1[k1]==V2[0] and in H2
               for k1 := 1; k1 <= len(H1)-2; k1++ {
                   if H1[k1] != V2[0] { continue }
                   if k1 >= len(H2)-1 { continue }
                   if k1 < 1 || k1 > len(H2)-2 { continue }
                   // i1: V1[i1]==H2[0] and V2[i1]==H2[k1]
                   for i1 := 1; i1 <= len(V1)-2; i1++ {
                       if V1[i1] != H2[0] { continue }
                       if i1 >= len(V2)-1 { continue }
                       if V2[i1] != H2[k1] { continue }
                       // k2: H2[k2]==V3[0]
                       for k2 := k1+2; k2 <= len(H2)-2; k2++ {
                           if H2[k2] != V3[0] { continue }
                           width := k2 - k1 + 1
                           if len(H3) != width { continue }
                           // i3: V2[i3]==H3[0], meets bottom
                           for i3 := i1+2; i3 <= len(V2)-2; i3++ {
                               if V2[i3] != H3[0] { continue }
                               off := i3 - i1
                               if off < 1 || off > len(V3)-2 { continue }
                               if V3[off] != H3[width-1] { continue }
                               // build grid
                               // determine size
                               bottom := i1 + len(V3) - 1
                               if len(V1)-1 > bottom { bottom = len(V1)-1 }
                               if len(V2)-1 > bottom { bottom = len(V2)-1 }
                               if i3 > bottom { bottom = i3 }
                               rows := bottom + 1
                               right := len(H1)-1
                               if len(H2)-1 > right { right = len(H2)-1 }
                               c3 := k1 + len(H3) - 1
                               if c3 > right { right = c3 }
                               if k2 > right { right = k2 }
                               cols := right + 1
                               // init grid
                               grid := make([][]byte, rows)
                               for r := 0; r < rows; r++ {
                                   grid[r] = make([]byte, cols)
                                   for c := 0; c < cols; c++ {
                                       grid[r][c] = '.'
                                   }
                               }
                               conflict := false
                               // place words
                               // H1 at (0,0)
                               for j := 0; j < len(H1); j++ {
                                   if grid[0][j] != '.' && grid[0][j] != H1[j] { conflict = true; break }
                                   grid[0][j] = H1[j]
                               }
                               if conflict { continue }
                               // H2 at (i1,0)
                               for j := 0; j < len(H2); j++ {
                                   if grid[i1][j] != '.' && grid[i1][j] != H2[j] { conflict = true; break }
                                   grid[i1][j] = H2[j]
                               }
                               if conflict { continue }
                               // H3 at (i3, k1)
                               for j := 0; j < len(H3); j++ {
                                   c := k1 + j
                                   if grid[i3][c] != '.' && grid[i3][c] != H3[j] { conflict = true; break }
                                   grid[i3][c] = H3[j]
                               }
                               if conflict { continue }
                               // V1 at (0,0)
                               for i := 0; i < len(V1); i++ {
                                   if grid[i][0] != '.' && grid[i][0] != V1[i] { conflict = true; break }
                                   grid[i][0] = V1[i]
                               }
                               if conflict { continue }
                               // V2 at (0,k1)
                               for i := 0; i < len(V2); i++ {
                                   if grid[i][k1] != '.' && grid[i][k1] != V2[i] { conflict = true; break }
                                   grid[i][k1] = V2[i]
                               }
                               if conflict { continue }
                               // V3 at (i1, k2)
                               for i := 0; i < len(V3); i++ {
                                   r := i1 + i
                                   if grid[r][k2] != '.' && grid[r][k2] != V3[i] { conflict = true; break }
                                   grid[r][k2] = V3[i]
                               }
                               if conflict { continue }
                               // check blank regions
                               total, inner := countBlanks(grid)
                               if total != 4 || inner != 2 {
                                   continue
                               }
                               // build strings
                               out := make([]string, rows)
                               for r := 0; r < rows; r++ {
                                   out[r] = string(grid[r])
                               }
                               // compare lex
                               if best == nil || lessGrid(out, best) {
                                   best = out
                               }
                           }
                       }
                   }
               }
           })
       })
   }
   if len(best) == 0 {
       fmt.Println("Impossible")
   } else {
       for _, line := range best {
           fmt.Println(line)
       }
   }
}

// count blank components and number of inner (not touching border)
func countBlanks(grid [][]byte) (total, inner int) {
   rows := len(grid)
   cols := len(grid[0])
   vis := make([][]bool, rows)
   for i := range vis {
       vis[i] = make([]bool, cols)
   }
   var dfs func(i, j int) bool
   dfs = func(i, j int) bool {
       stack := [][2]int{{i, j}}
       vis[i][j] = true
       touches := false
       for len(stack) > 0 {
           x, y := stack[len(stack)-1][0], stack[len(stack)-1][1]
           stack = stack[:len(stack)-1]
           if x == 0 || y == 0 || x == rows-1 || y == cols-1 {
               touches = true
           }
           for _, d := range [][2]int{{1,0},{-1,0},{0,1},{0,-1}} {
               nx, ny := x+d[0], y+d[1]
               if nx>=0 && nx<rows && ny>=0 && ny<cols && !vis[nx][ny] && grid[nx][ny]=='.' {
                   vis[nx][ny] = true
                   stack = append(stack, [2]int{nx, ny})
               }
           }
       }
       total++
       if !touches {
           inner++
       }
       return true
   }
   for i := 0; i < rows; i++ {
       for j := 0; j < cols; j++ {
           if grid[i][j]=='.' && !vis[i][j] {
               dfs(i, j)
           }
       }
   }
   return total, inner
}

// lex compare grids
func lessGrid(a, b []string) bool {
   for i := 0; i < len(a) && i < len(b); i++ {
       if a[i] < b[i] {
           return true
       } else if a[i] > b[i] {
           return false
       }
   }
   return len(a) < len(b)
}
