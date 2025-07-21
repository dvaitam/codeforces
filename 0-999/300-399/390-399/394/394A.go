package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   input, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       fmt.Fprintln(os.Stderr, err)
       return
   }
   input = strings.TrimSpace(input)
   // parse counts
   p := strings.Index(input, "+")
   q := strings.Index(input, "=")
   if p < 0 || q < 0 || p > q {
       fmt.Println("Impossible")
       return
   }
   aCount := p
   bCount := q - p - 1
   cCount := len(input) - q - 1
   // already correct
   if aCount + bCount == cCount {
       fmt.Println(input)
       return
   }
   // try moving one stick from i to j
   counts := []int{aCount, bCount, cCount}
   for i := 0; i < 3; i++ {
       for j := 0; j < 3; j++ {
           if i == j {
               continue
           }
           // ensure at least 1 remains
           if counts[i] <= 1 {
               continue
           }
           newCounts := []int{counts[0], counts[1], counts[2]}
           newCounts[i]--
           newCounts[j]++
           if newCounts[0] + newCounts[1] == newCounts[2] {
               // build result
               a := strings.Repeat("|", newCounts[0])
               b := strings.Repeat("|", newCounts[1])
               c := strings.Repeat("|", newCounts[2])
               // output corrected expression
               fmt.Println(a + "+" + b + "=" + c)
               return
           }
       }
   }
   fmt.Println("Impossible")
}
