package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

// isPrime checks if n is a prime number
func isPrime(n int64) bool {
   if n < 2 {
       return false
   }
   limit := int64(math.Sqrt(float64(n)))
   for i := int64(2); i <= limit; i++ {
       if n%i == 0 {
           return false
       }
   }
   return true
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       arr := make([]int64, n)
       var sum int64
       oddIdx := -1
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &arr[i])
           sum += arr[i]
           if arr[i]%2 != 0 {
               oddIdx = i
           }
       }
       if isPrime(sum) {
           // exclude one odd element
           fmt.Fprintln(writer, n-1)
           for i := 0; i < n; i++ {
               if i == oddIdx {
                   continue
               }
               fmt.Fprint(writer, i+1)
               if i != n-1 {
                   fmt.Fprint(writer, " ")
               }
           }
           fmt.Fprintln(writer)
       } else {
           fmt.Fprintln(writer, n)
           for i := 0; i < n; i++ {
               fmt.Fprint(writer, i+1)
               if i != n-1 {
                   fmt.Fprint(writer, " ")
               }
           }
           fmt.Fprintln(writer)
       }
   }
}
