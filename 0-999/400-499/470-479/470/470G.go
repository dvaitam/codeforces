package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   var s [2]string
   i := 0
   for scanner.Scan() {
       if i < 2 {
           s[i] = scanner.Text()
       }
       i++
       if i >= 2 {
           break
       }
   }
   // Compute Hamming distance
   n := len(s[0])
   dist := 0
   for j := 0; j < n; j++ {
       if s[0][j] != s[1][j] {
           dist++
       }
   }
   fmt.Println(dist)
}
