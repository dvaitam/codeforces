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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   v := make([]int, n)
   for i := 1; i < n; i++ {
       fmt.Fscan(reader, &v[i])
   }

   // Print the first element as root
   fmt.Fprintln(writer, v[1])

   findPath := make([]bool, n+1)
   findPath[v[1]] = true

   p := 2
   last := n
   ans := make([][2]int, 0, n-1)

   for {
       for last > 0 && findPath[last] {
           last--
       }
       if last <= 0 {
           break
       }
       findPath[last] = true

       for p < n && !findPath[v[p]] {
           ans = append(ans, [2]int{v[p-1], v[p]})
           findPath[v[p]] = true
           p++
       }
       // connect the last missing node
       ans = append(ans, [2]int{last, v[p-1]})
       p++
   }

   // Output edges
   for i := 0; i < len(ans); i++ {
       fmt.Fprintln(writer, ans[i][0], ans[i][1])
   }
}
