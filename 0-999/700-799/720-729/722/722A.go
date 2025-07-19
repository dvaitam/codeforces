package main

import "fmt"

func main() {
   var n int
   var s string
   fmt.Scan(&n, &s)
   b := []byte(s)
   // Fix minute tens digit (must be 0-5)
   if b[3] > '5' {
       b[3] = '0'
   }
   if n == 12 {
       // 12-hour format: hours must be 01-12
       if b[0] != '1' && b[1] == '0' {
           b[0] = '1'
       } else if b[0] > '1' || (b[0] == '1' && b[1] > '2') {
           b[0] = '0'
       }
   } else {
       // 24-hour format: hours must be 00-23
       if b[0] > '2' || (b[0] == '2' && b[1] > '3') {
           b[0] = '0'
       }
   }
   fmt.Println(string(b))
}
