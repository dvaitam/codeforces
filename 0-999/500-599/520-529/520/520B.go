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
   // If starting number is greater or equal, only subtract operations are needed
   if n >= m {
       fmt.Println(n - m)
       return
   }
   // Greedy reverse operations: from m back to n
   steps := 0
   for m > n {
       if m%2 == 0 {
           m /= 2
       } else {
           m++
       }
       steps++
   }
   // If m is now less than n, add remaining subtractions
   steps += n - m
   fmt.Println(steps)
}
