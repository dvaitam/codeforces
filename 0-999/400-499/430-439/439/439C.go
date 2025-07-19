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

   var n, k, p int
   if _, err := fmt.Fscan(reader, &n, &k, &p); err != nil {
       return
   }
   odd := make([]int, 0, n)
   even := make([]int, 0, n)
   for i := 0; i < n; i++ {
       var v int
       fmt.Fscan(reader, &v)
       if v%2 == 0 {
           even = append(even, v)
       } else {
           odd = append(odd, v)
       }
   }
   oddCnt := len(odd)
   evenCnt := len(even)
   needOdd := k - p
   // Check feasibility
   if oddCnt < needOdd || (oddCnt-needOdd)%2 != 0 || ((oddCnt-needOdd)/2+evenCnt) < p {
       fmt.Fprintln(writer, "NO")
       return
   }
   fmt.Fprintln(writer, "YES")
   if p > 0 {
       // Build k-p subsets with odd sum
       for i := 0; i < needOdd; i++ {
           fmt.Fprintf(writer, "1 %d\n", odd[i])
       }
       // merge even numbers and remaining odds
       remOdd := odd[needOdd:]
       merged := make([]int, 0, len(even)+len(remOdd))
       merged = append(merged, even...)
       merged = append(merged, remOdd...)
       idx := 0
       // first p-1 even-sum subsets
       for j := 0; j < p-1; j++ {
           if merged[idx]%2 != 0 {
               fmt.Fprintf(writer, "2 %d %d\n", merged[idx], merged[idx+1])
               idx += 2
           } else {
               fmt.Fprintf(writer, "1 %d\n", merged[idx])
               idx++
           }
       }
       // last subset: all remaining
       rem := merged[idx:]
       fmt.Fprintf(writer, "%d", len(rem))
       for _, v := range rem {
           fmt.Fprintf(writer, " %d", v)
       }
       fmt.Fprintln(writer)
   } else {
       // p == 0: all subsets are odd-sum
       // merge all odds and evens
       merged := make([]int, 0, len(odd)+len(even))
       merged = append(merged, odd...)
       merged = append(merged, even...)
       // first k-1 odd-sum subsets
       for i := 0; i < k-1; i++ {
           fmt.Fprintf(writer, "1 %d\n", merged[i])
       }
       // last subset: all remaining
       rem := merged[k-1:]
       fmt.Fprintf(writer, "%d", len(rem))
       for _, v := range rem {
           fmt.Fprintf(writer, " %d", v)
       }
       fmt.Fprintln(writer)
   }
}
