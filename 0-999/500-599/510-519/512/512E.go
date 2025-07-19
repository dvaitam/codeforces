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

   var n int
   fmt.Fscan(reader, &n)
   // Read original and goal diagonals
   orig := make([][2]int, n-3)
   for i := 0; i < n-3; i++ {
       fmt.Fscan(reader, &orig[i][0], &orig[i][1])
   }
   goal := make([][2]int, n-3)
   for i := 0; i < n-3; i++ {
       fmt.Fscan(reader, &goal[i][0], &goal[i][1])
   }
   // Initialize adjacency for original and goal
   adjOrig := make([][]bool, n+1)
   adjGoal := make([][]bool, n+1)
   for i := 0; i <= n; i++ {
       adjOrig[i] = make([]bool, n+1)
       adjGoal[i] = make([]bool, n+1)
   }
   // add polygon edges
   for i := 1; i <= n; i++ {
       j := i%n + 1
       adjOrig[i][j], adjOrig[j][i] = true, true
       adjGoal[i][j], adjGoal[j][i] = true, true
   }
   // add diagonals
   for _, d := range orig {
       x, y := d[0], d[1]
       adjOrig[x][y], adjOrig[y][x] = true, true
   }
   for _, d := range goal {
       x, y := d[0], d[1]
       adjGoal[x][y], adjGoal[y][x] = true, true
   }

   // Simulation on original
   a1 := make([][]bool, n+1)
   for i := 0; i <= n; i++ {
       a1[i] = make([]bool, n+1)
       copy(a1[i], adjOrig[i])
   }
   var ans1 [][2]int
   for {
       moved := false
       for i := 3; i < n; i++ {
           if !a1[1][i] {
               // find k >= i with a1[1][k]
               k := i
               for ; k <= n; k++ {
                   if a1[1][k] {
                       break
                   }
               }
               // find x: common neighbor of i-1 and k
               var x int
               for j := 2; j <= n; j++ {
                   if a1[i-1][j] && a1[k][j] {
                       x = j
                       break
                   }
               }
               ans1 = append(ans1, [2]int{i - 1, k})
               // perform flip: remove (i-1,k), add (1,x)
               a1[i-1][k], a1[k][i-1] = false, false
               a1[1][x], a1[x][1] = true, true
               moved = true
               break
           }
       }
       if !moved {
           break
       }
   }

   // Simulation on goal
   a2 := make([][]bool, n+1)
   for i := 0; i <= n; i++ {
       a2[i] = make([]bool, n+1)
       copy(a2[i], adjGoal[i])
   }
   var ans2 [][2]int
   for {
       moved := false
       for i := 3; i < n; i++ {
           if !a2[1][i] {
               k := i
               for ; k <= n; k++ {
                   if a2[1][k] {
                       break
                   }
               }
               var x int
               for j := 2; j <= n; j++ {
                   if a2[i-1][j] && a2[k][j] {
                       x = j
                       break
                   }
               }
               // record added diagonal (1,x)
               ans2 = append(ans2, [2]int{1, x})
               // perform flip: remove (i-1,k), add (1,x)
               a2[i-1][k], a2[k][i-1] = false, false
               a2[1][x], a2[x][1] = true, true
               moved = true
               break
           }
       }
       if !moved {
           break
       }
   }

   // Output combined moves
   total := len(ans1) + len(ans2)
   fmt.Fprintln(writer, total)
   for _, p := range ans1 {
       fmt.Fprintln(writer, p[0], p[1])
   }
   // reverse ans2
   for i := len(ans2) - 1; i >= 0; i-- {
       fmt.Fprintln(writer, ans2[i][0], ans2[i][1])
   }
}
