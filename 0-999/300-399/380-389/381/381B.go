package main

import "fmt"

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   // Frequency array, cap at 2
   const maxA = 5000
   freq := make([]int, maxA+1)
   maxV := 0
   for i := 0; i < n; i++ {
       var a int
       if _, err := fmt.Scan(&a); err != nil {
           return
       }
       if a > maxV {
           maxV = a
       }
       if freq[a] < 2 {
           freq[a]++
       }
   }
   ascCount := 0
   for i := 0; i <= maxV; i++ {
       if freq[i] > 0 {
           ascCount++
       }
   }
   dupCount := 0
   for i := 0; i < maxV; i++ {
       if freq[i] > 1 {
           dupCount++
       }
   }
   total := ascCount + dupCount
   fmt.Println(total)

   // Build sequence: ascending unique, then descending duplicates
   // Print ascending unique elements
   for i := 0; i < maxV; i++ {
       if freq[i] > 0 {
           fmt.Printf("%d ", i)
       }
   }
   // Print maxV
   fmt.Printf("%d", maxV)
   // Print descending duplicates (excluding maxV)
   for i := maxV - 1; i >= 0; i-- {
       if freq[i] > 1 {
           fmt.Printf(" %d", i)
       }
   }
   fmt.Println()
}
