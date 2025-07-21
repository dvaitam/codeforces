package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"

   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read n
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   sticks := make([]int, n)
   for i := 0; i < n; i++ {
       if _, err := fmt.Fscan(reader, &sticks[i]); err != nil {
           return
       }
   }
   // Sort in descending order
   sort.Slice(sticks, func(i, j int) bool { return sticks[i] > sticks[j] })
   sides := make([]int, 0, n/2)
   // Greedy pair selection
   for i := 0; i < n-1; {
       if sticks[i]-sticks[i+1] <= 1 {
           // can form a pair of length sticks[i+1]
           sides = append(sides, sticks[i+1])
           i += 2
       } else {
           i++
       }
   }
   // Combine sides into rectangles
   var totalArea int64
   m := len(sides)
   for i := 0; i+1 < m; i += 2 {
       a := int64(sides[i])
       b := int64(sides[i+1])
       totalArea += a * b
   }
   // Output result
   // Use Print instead of Fprintln to avoid extra spaces
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   out.WriteString(strconv.FormatInt(totalArea, 10))
   out.WriteByte('\n')
}
