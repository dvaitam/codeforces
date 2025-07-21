package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   scanner.Split(bufio.ScanLines)
   if !scanner.Scan() {
       return
   }
   var n int
   fmt.Sscanf(scanner.Text(), "%d", &n)
   for b := 0; b < n; b++ {
       // Read board lines, skipping empty separators
       var grid [10]string
       i := 0
       for i < 10 && scanner.Scan() {
           line := scanner.Text()
           if len(line) == 0 {
               continue
           }
           grid[i] = line
           i++
       }
       if isValid(grid) {
           fmt.Println("YES")
       } else {
           fmt.Println("NO")
       }
   }
}

// isValid checks if the given 10x10 board satisfies Battleship rules
func isValid(grid [10]string) bool {
   // Convert to fixed grid
   var g [10][10]byte
   for i := 0; i < 10; i++ {
       for j := 0; j < 10; j++ {
           if j < len(grid[i]) && grid[i][j] == '*' {
               g[i][j] = '*'
           } else {
               g[i][j] = '0'
           }
       }
   }
   // Check diagonal adjacency
   diag := [4][2]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
   for i := 0; i < 10; i++ {
       for j := 0; j < 10; j++ {
           if g[i][j] != '*' {
               continue
           }
           for _, d := range diag {
               ni, nj := i+d[0], j+d[1]
               if ni >= 0 && ni < 10 && nj >= 0 && nj < 10 && g[ni][nj] == '*' {
                   return false
               }
           }
       }
   }
   // Track visited cells
   var vis [10][10]bool
   counts := map[int]int{1: 0, 2: 0, 3: 0, 4: 0}
   dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
   for i := 0; i < 10; i++ {
       for j := 0; j < 10; j++ {
           if g[i][j] == '*' && !vis[i][j] {
               // DFS stack
               stack := [][2]int{{i, j}}
               vis[i][j] = true
               var cells [][2]int
               cells = append(cells, [2]int{i, j})
               for len(stack) > 0 {
                   cur := stack[len(stack)-1]
                   stack = stack[:len(stack)-1]
                   ci, cj := cur[0], cur[1]
                   for _, d := range dirs {
                       ni, nj := ci+d[0], cj+d[1]
                       if ni >= 0 && ni < 10 && nj >= 0 && nj < 10 && g[ni][nj] == '*' && !vis[ni][nj] {
                           vis[ni][nj] = true
                           stack = append(stack, [2]int{ni, nj})
                           cells = append(cells, [2]int{ni, nj})
                       }
                   }
               }
               size := len(cells)
               if size < 1 || size > 4 {
                   return false
               }
               // Determine bounding box
               minI, maxI, minJ, maxJ := 10, -1, 10, -1
               setC := make(map[[2]int]bool)
               for _, c := range cells {
                   setC[c] = true
                   if c[0] < minI {
                       minI = c[0]
                   }
                   if c[0] > maxI {
                       maxI = c[0]
                   }
                   if c[1] < minJ {
                       minJ = c[1]
                   }
                   if c[1] > maxJ {
                       maxJ = c[1]
                   }
               }
               if minI == maxI {
                   // horizontal line
                   if maxJ-minJ+1 != size {
                       return false
                   }
                   for jj := minJ; jj <= maxJ; jj++ {
                       if !setC[[2]int{minI, jj}] {
                           return false
                       }
                   }
               } else if minJ == maxJ {
                   // vertical line
                   if maxI-minI+1 != size {
                       return false
                   }
                   for ii := minI; ii <= maxI; ii++ {
                       if !setC[[2]int{ii, minJ}] {
                           return false
                       }
                   }
               } else {
                   // bent shape
                   return false
               }
               counts[size]++
           }
       }
   }
   // Required ships: 4 of size1, 3 of size2, 2 of size3, 1 of size4
   required := map[int]int{1: 4, 2: 3, 3: 2, 4: 1}
   for k, v := range required {
       if counts[k] != v {
           return false
       }
   }
   return true
}
