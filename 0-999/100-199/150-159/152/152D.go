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
   // grid with padding
   s := make([][]byte, n+3)
   for i := range s {
       s[i] = make([]byte, m+3)
       for j := range s[i] {
           s[i][j] = '.'
       }
   }
   ur := make([]bool, n+3)
   uc := make([]bool, m+3)
   cp := 0
   // read grid and detect potential strips
   for i := 1; i <= n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       for j := 1; j <= m; j++ {
           s[i][j] = line[j-1]
       }
       for j := m; j >= 1; j-- {
           if s[i][j] == '#' {
               cp++
               if s[i][j+1] == '#' && s[i][j+2] == '#' {
                   ur[i] = true
               }
               if s[i-1][j] == '#' && s[i-2][j] == '#' {
                   uc[j] = true
               }
           }
       }
   }
   // candidate rows and columns
   var r, c []int
   for i := 1; i <= n; i++ {
       if ur[i] {
           r = append(r, i)
       }
   }
   for j := 1; j <= m; j++ {
       if uc[j] {
           c = append(c, j)
       }
   }
   if len(r) > 4 {
       r = []int{r[0], r[1], r[len(r)-2], r[len(r)-1]}
   }
   if len(c) > 4 {
       c = []int{c[0], c[1], c[len(c)-2], c[len(c)-1]}
   }
   const inf = -123456789
   // function to check two rectangles share all '#'
   var Gans func(A, B, C, D, a, b, c2, d int) bool
   Gans = func(A, B, C, D, a, b, c2, d int) bool {
       cnt := 0
       // first rectangle border
       for i := A; i <= B; i++ {
           if s[i][C] == '.' || s[i][D] == '.' {
               return false
           }
           s[i][C], s[i][D] = 'A', 'A'
           cnt += 2
       }
       for i := C + 1; i < D; i++ {
           if s[A][i] == '.' || s[B][i] == '.' {
               // revert first side
               for x := A; x <= B; x++ {
                   if s[x][C] == 'A' {
                       s[x][C] = '#'
                   }
                   if s[x][D] == 'A' {
                       s[x][D] = '#'
                   }
               }
               return false
           }
           s[A][i], s[B][i] = 'A', 'A'
           cnt += 2
       }
       // second rectangle border
       for i := a; i <= b; i++ {
           if s[i][c2] == '.' || s[i][d] == '.' {
               // revert all A
               for x := A; x <= B; x++ {
                   if s[x][C] == 'A' {
                       s[x][C] = '#'
                   }
                   if s[x][D] == 'A' {
                       s[x][D] = '#'
                   }
               }
               for x := C + 1; x < D; x++ {
                   if s[A][x] == 'A' {
                       s[A][x] = '#'
                   }
                   if s[B][x] == 'A' {
                       s[B][x] = '#'
                   }
               }
               return false
           }
           if s[i][c2] == '#' {
               cnt++
           }
           if s[i][d] == '#' {
               cnt++
           }
           s[i][c2], s[i][d] = 'A', 'A'
       }
       for i := c2 + 1; i < d; i++ {
           if s[a][i] == '.' || s[b][i] == '.' {
               // revert all A
               for x := A; x <= B; x++ {
                   if s[x][C] == 'A' {
                       s[x][C] = '#'
                   }
                   if s[x][D] == 'A' {
                       s[x][D] = '#'
                   }
               }
               for x := C + 1; x < D; x++ {
                   if s[A][x] == 'A' {
                       s[A][x] = '#'
                   }
                   if s[B][x] == 'A' {
                       s[B][x] = '#'
                   }
               }
               for x := a; x <= b; x++ {
                   if s[x][c2] == 'A' {
                       s[x][c2] = '#'
                   }
                   if s[x][d] == 'A' {
                       s[x][d] = '#'
                   }
               }
               return false
           }
           if s[a][i] == '#' {
               cnt++
           }
           if s[b][i] == '#' {
               cnt++
           }
           s[a][i], s[b][i] = 'A', 'A'
       }
       // revert all A to '#'
       for i := A; i <= B; i++ {
           if s[i][C] == 'A' {
               s[i][C] = '#'
           }
           if s[i][D] == 'A' {
               s[i][D] = '#'
           }
       }
       for i := C + 1; i < D; i++ {
           if s[A][i] == 'A' {
               s[A][i] = '#'
           }
           if s[B][i] == 'A' {
               s[B][i] = '#'
           }
       }
       for i := a; i <= b; i++ {
           if s[i][c2] == 'A' {
               s[i][c2] = '#'
           }
           if s[i][d] == 'A' {
               s[i][d] = '#'
           }
       }
       for i := c2 + 1; i < d; i++ {
           if s[a][i] == 'A' {
               s[a][i] = '#'
           }
           if s[b][i] == 'A' {
               s[b][i] = '#'
           }
       }
       return cnt == cp
   }
   // search two rectangles
   for i := 0; i < len(r); i++ {
       for j := i + 1; j < len(r); j++ {
           if r[i]+1 >= r[j] {
               continue
           }
           for k := 0; k < len(c); k++ {
               for l := k + 1; l < len(c); l++ {
                   if c[k]+1 >= c[l] {
                       continue
                   }
                   for A := 0; A < len(r); A++ {
                       for B := A + 1; B < len(r); B++ {
                           if r[A]+1 >= r[B] {
                               continue
                           }
                           for C2 := 0; C2 < len(c); C2++ {
                               for D := C2 + 1; D < len(c); D++ {
                                   if c[C2]+1 >= c[D] {
                                       continue
                                   }
                                   if Gans(r[i], r[j], c[k], c[l], r[A], r[B], c[C2], c[D]) {
                                       fmt.Fprintln(writer, "YES")
                                       fmt.Fprintf(writer, "%d %d %d %d\n", r[i], c[k], r[j], c[l])
                                       fmt.Fprintf(writer, "%d %d %d %d\n", r[A], c[C2], r[B], c[D])
                                       return
                                   }
                               }
                           }
                       }
                   }
               }
           }
       }
   }
   fmt.Fprintln(writer, "NO")
}
