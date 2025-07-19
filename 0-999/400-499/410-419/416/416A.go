package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)

   var l, r int64 = -2000000000, 2000000000
   for i := 0; i < n; i++ {
       var op string
       var number int64
       var verify string
       fmt.Fscan(reader, &op, &number, &verify)
       if op == ">" {
           if verify == "Y" {
               if number+1 > l {
                   l = number + 1
               }
           } else {
               if number < r {
                   r = number
               }
           }
       } else if op == "<" {
           if verify == "Y" {
               if number-1 < r {
                   r = number - 1
               }
           } else {
               if number > l {
                   l = number
               }
           }
       } else if op == ">=" {
           if verify == "Y" {
               if number > l {
                   l = number
               }
           } else {
               if number-1 < r {
                   r = number - 1
               }
           }
       } else if op == "<=" {
           if verify == "Y" {
               if number < r {
                   r = number
               }
           } else {
               if number+1 > l {
                   l = number + 1
               }
           }
       }
       if l > r {
           fmt.Println("Impossible")
           return
       }
   }
   fmt.Println(l)
}
