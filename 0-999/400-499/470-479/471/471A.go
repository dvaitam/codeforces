package main

import "fmt"

func main() {
   var a [6]int
   for i := 0; i < 6; i++ {
       if _, err := fmt.Scan(&a[i]); err != nil {
           return
       }
   }
   var freq [10]int
   for _, v := range a {
       if v >= 1 && v <= 9 {
           freq[v]++
       }
   }
   // find leg length appearing at least 4 times
   legLen := 0
   for v := 1; v <= 9; v++ {
       if freq[v] >= 4 {
           legLen = v
           break
       }
   }
   if legLen == 0 {
       fmt.Println("Alien")
       return
   }
   // collect remaining two sticks
   rem := make([]int, 0, 2)
   removed := 0
   for _, v := range a {
       if v == legLen && removed < 4 {
           removed++
       } else {
           rem = append(rem, v)
       }
   }
   if len(rem) != 2 {
       fmt.Println("Alien")
       return
   }
   if rem[0] == rem[1] {
       fmt.Println("Elephant")
   } else {
       fmt.Println("Bear")
   }
}
