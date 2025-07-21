package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var sFinal, sBirth string
   if _, err := fmt.Fscan(reader, &sFinal, &sBirth); err != nil {
       return
   }
   var fD, fM, fY int
   fmt.Sscanf(sFinal, "%d.%d.%d", &fD, &fM, &fY)
   var b0, b1, b2 int
   fmt.Sscanf(sBirth, "%d.%d.%d", &b0, &b1, &b2)
   tokens := []int{b0, b1, b2}
   // all permutations of three tokens
   perms := [6][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
   // month lengths
   monthDays := []int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
   for _, p := range perms {
       d := tokens[p[0]]
       m := tokens[p[1]]
       y := tokens[p[2]]
       // valid month
       if m < 1 || m > 12 {
           continue
       }
       // determine days in month, considering leap years
       days := monthDays[m]
       yearBirth := 2000 + y
       if m == 2 && yearBirth%4 == 0 {
           days = 29
       }
       // valid day
       if d < 1 || d > days {
           continue
       }
       // compute age
       yearFinal := 2000 + fY
       age := yearFinal - yearBirth
       if fM < m || (fM == m && fD < d) {
           age--
       }
       if age >= 18 {
           fmt.Println("YES")
           return
       }
   }
   fmt.Println("NO")
}
