package main

import "fmt"

func main() {
   var a, b, c, d int
   if _, err := fmt.Scan(&a, &b, &c, &d); err != nil {
       return
   }
   // Case 1
   if c == 0 && d == 0 {
       if a == b+1 {
           fmt.Println("YES")
           for i := 0; i < b; i++ {
               fmt.Print("0 1 ")
           }
           fmt.Print("0")
           return
       }
   }
   // Case 2
   if a == 0 && b == 0 {
       if c+1 == d {
           fmt.Println("YES")
           for i := 0; i < c; i++ {
               fmt.Print("3 2 ")
           }
           fmt.Print("3")
           return
       }
   }
   // Case 3
   if b-a+d == c && b-a >= 0 {
       fmt.Println("YES")
       x := a
       y := b - a
       z := d
       for i := 0; i < x; i++ {
       fmt.Print("0 1 ")
       }
       for i := 0; i < y; i++ {
       fmt.Print("2 1 ")
       }
       for i := 0; i < z; i++ {
       fmt.Print("2 3 ")
       }
       return
   }
   // Case 4
   if b-a+d-1 == c && b-a-1 >= 0 {
       fmt.Println("YES")
       x := a
       y := b - a - 1
       z := d
       fmt.Print("1 ")
       for i := 0; i < x; i++ {
           fmt.Print("0 1 ")
       }
       for i := 0; i < y; i++ {
           fmt.Print("2 1 ")
       }
       for i := 0; i < z; i++ {
           fmt.Print("2 3 ")
       }
       return
   }
   // Case 5
   if b-a+d+1 == c && b-a >= 0 {
       fmt.Println("YES")
       x := a
       y := b - a
       z := d
       for i := 0; i < x; i++ {
       fmt.Print("0 1 ")
       }
       for i := 0; i < y; i++ {
       fmt.Print("2 1 ")
       }
       for i := 0; i < z; i++ {
       fmt.Print("2 3 ")
       }
       fmt.Print("2 ")
       return
   }
   // Case 6
   if b-a+d == c && b-a-1 >= 0 {
       fmt.Println("YES")
       x := a
       y := b - a - 1
       z := d
       fmt.Print("1 ")
       for i := 0; i < x; i++ {
           fmt.Print("0 1 ")
       }
       for i := 0; i < y; i++ {
           fmt.Print("2 1 ")
       }
       for i := 0; i < z; i++ {
           fmt.Print("2 3 ")
       }
       fmt.Print("2 ")
       return
   }
   // No valid sequence
   fmt.Println("NO")
}
