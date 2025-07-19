package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read six counts
   A := make([]int, 7)
   total := 0
   for i := 1; i <= 6; i++ {
       var v int
       if _, err := fmt.Fscan(reader, &v); err != nil {
           return
       }
       A[i] = v
       total += v
   }
   // Read number of customers
   var N int
   if _, err := fmt.Fscan(reader, &N); err != nil {
       return
   }
   // consume end of line
   reader.ReadString('\n')
   // Prepare
   ans := make([]int, N)
   type P struct{ x, id int }
   C := make([]P, 0, N)

   // Helper to map size string to int
   size := func(s string) int {
       switch s {
       case "S":
           return 1
       case "M":
           return 2
       case "L":
           return 3
       case "XL":
           return 4
       case "XXL":
           return 5
       case "XXXL":
           return 6
       }
       return 0
   }
   // Read customer preferences
   for i := 0; i < N; i++ {
       line, err := reader.ReadString('\n')
       if err != nil {
           return
       }
       line = strings.TrimSpace(line)
       parts := strings.Split(line, ",")
       sz1 := size(parts[0])
       if len(parts) == 1 {
           // single size
           ans[i] = sz1
           A[sz1]--
       } else {
           sz2 := size(parts[1])
           x := sz1
           if sz2 < x {
               x = sz2
           }
           C = append(C, P{x: x, id: i})
       }
   }
   // Check total shirts
   if total < N {
       fmt.Println("NO")
       return
   }
   // Sort flexible by smaller size
   sort.Slice(C, func(i, j int) bool {
       return C[i].x < C[j].x
   })
   // Check for negative inventory after direct assigns
   for i := 1; i <= 6; i++ {
       if A[i] < 0 {
           fmt.Println("NO")
           return
       }
   }
   // Assign flexible
   for _, p := range C {
       nx := p.x
       if A[nx] > 0 {
           A[nx]--
           ans[p.id] = nx
       } else if nx+1 <= 6 && A[nx+1] > 0 {
           A[nx+1]--
           ans[p.id] = nx + 1
       } else {
           fmt.Println("NO")
           return
       }
   }
   // Output
   fmt.Println("YES")
   tab := []string{"S", "M", "L", "XL", "XXL", "XXXL"}
   for i := 0; i < N; i++ {
       idx := ans[i]
       if idx >= 1 && idx <= 6 {
           fmt.Println(tab[idx-1])
       } else {
           fmt.Println()
       }
   }
}
