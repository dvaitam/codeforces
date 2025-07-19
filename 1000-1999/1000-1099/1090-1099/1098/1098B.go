package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       grid[i] = make([]byte, m)
       var line []byte
       for len(line) < m {
           part, _, err := reader.ReadLine()
           if err != nil {
               break
           }
           line = append(line, part...)
       }
       copy(grid[i], line[:m])
   }
   letters := []byte{'A', 'G', 'C', 'T'}
   perms := make([][]byte, 0, 24)
   var permute func(a []byte, l int)
   permute = func(a []byte, l int) {
       if l == len(a)-1 {
           p := make([]byte, len(a))
           copy(p, a)
           perms = append(perms, p)
           return
       }
       for i := l; i < len(a); i++ {
           a[l], a[i] = a[i], a[l]
           permute(a, l+1)
           a[l], a[i] = a[i], a[l]
       }
   }
   permute(letters, 0)
   bestScore := -1
   bestType := 0
   bestPerm := make([]byte, 4)
   var bestFlips []bool
   // row-based
   for _, perm := range perms {
       flips := make([]bool, n)
       total := 0
       for i := 0; i < n; i++ {
           var w1, w2 byte
           if i%2 == 0 {
               w1, w2 = perm[0], perm[1]
           } else {
               w1, w2 = perm[2], perm[3]
           }
           s1, s2 := 0, 0
           for j := 0; j < m; j++ {
               if j%2 == 0 {
                   if grid[i][j] == w1 {
                       s1++
                   }
                   if grid[i][j] == w2 {
                       s2++
                   }
               } else {
                   if grid[i][j] == w2 {
                       s1++
                   }
                   if grid[i][j] == w1 {
                       s2++
                   }
               }
           }
           if s1 >= s2 {
               total += s1
               flips[i] = false
           } else {
               total += s2
               flips[i] = true
           }
       }
       if total > bestScore {
           bestScore = total
           bestType = 0
           copy(bestPerm, perm)
           bestFlips = flips
       }
   }
   // column-based
   for _, perm := range perms {
       flips := make([]bool, m)
       total := 0
       for j := 0; j < m; j++ {
           var w1, w2 byte
           if j%2 == 0 {
               w1, w2 = perm[0], perm[1]
           } else {
               w1, w2 = perm[2], perm[3]
           }
           s1, s2 := 0, 0
           for i := 0; i < n; i++ {
               if i%2 == 0 {
                   if grid[i][j] == w1 {
                       s1++
                   }
                   if grid[i][j] == w2 {
                       s2++
                   }
               } else {
                   if grid[i][j] == w2 {
                       s1++
                   }
                   if grid[i][j] == w1 {
                       s2++
                   }
               }
           }
           if s1 >= s2 {
               total += s1
               flips[j] = false
           } else {
               total += s2
               flips[j] = true
           }
       }
       if total > bestScore {
           bestScore = total
           bestType = 1
           copy(bestPerm, perm)
           bestFlips = flips
       }
   }
   // reconstruct
   out := make([][]byte, n)
   for i := 0; i < n; i++ {
       out[i] = make([]byte, m)
   }
   if bestType == 0 {
       for i := 0; i < n; i++ {
           var w1, w2 byte
           if i%2 == 0 {
               w1, w2 = bestPerm[0], bestPerm[1]
           } else {
               w1, w2 = bestPerm[2], bestPerm[3]
           }
           if bestFlips[i] {
               w1, w2 = w2, w1
           }
           for j := 0; j < m; j++ {
               if j%2 == 0 {
                   out[i][j] = w1
               } else {
                   out[i][j] = w2
               }
           }
       }
   } else {
       for j := 0; j < m; j++ {
           var w1, w2 byte
           if j%2 == 0 {
               w1, w2 = bestPerm[0], bestPerm[1]
           } else {
               w1, w2 = bestPerm[2], bestPerm[3]
           }
           if bestFlips[j] {
               w1, w2 = w2, w1
           }
           for i := 0; i < n; i++ {
               if i%2 == 0 {
                   out[i][j] = w1
               } else {
                   out[i][j] = w2
               }
           }
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < n; i++ {
       writer.Write(out[i])
       writer.WriteByte('\n')
   }
}
