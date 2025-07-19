package main

import "fmt"

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   a := make([]int, 0)
   for i := 0; i < n; i++ {
       var p int
       fmt.Scan(&p)
       if p == 0 {
           sz := len(a)
           if sz < 3 {
               if sz == 2 {
                   fmt.Println("pushStack")
                   fmt.Println("pushQueue")
                   fmt.Println("2 popStack popQueue")
               } else if sz == 1 {
                   fmt.Println("pushStack")
                   fmt.Println("1 popStack")
               } else {
                   fmt.Println("0")
               }
               a = a[:0]
               continue
           }
           // sz >= 3: find three largest via sentinel
           a = append(a, -1)
           m1, m2, m3 := sz, sz, sz
           for idx, val := range a {
               if a[m1] < val {
                   m1 = idx
               }
           }
           for idx, val := range a {
               if idx != m1 && a[m2] < val {
                   m2 = idx
               }
           }
           for idx, val := range a {
               if idx != m1 && idx != m2 && a[m3] < val {
                   m3 = idx
               }
           }
           // perform pushes for first sz elements
           for idx := 0; idx < sz; idx++ {
               switch idx {
               case m1:
                   fmt.Println("pushStack")
               case m2:
                   fmt.Println("pushQueue")
               case m3:
                   fmt.Println("pushFront")
               default:
                   fmt.Println("pushBack")
               }
           }
           fmt.Println("3 popStack popQueue popFront")
           a = a[:0]
       } else {
           a = append(a, p)
       }
   }
   // remaining elements: pushStack
   for range a {
       fmt.Println("pushStack")
   }
}
