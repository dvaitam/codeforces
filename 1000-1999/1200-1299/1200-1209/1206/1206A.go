package main

import "fmt"

func main() {
   var n int
   if fmt.Scan(&n) != 1 {
       return
   }
   a := make([]int, n)
   for i := range a {
       fmt.Scan(&a[i])
   }

   var m int
   fmt.Scan(&m)
   b := make([]int, m)
   for i := range b {
       fmt.Scan(&b[i])
   }


	ma := make(map[int]bool, len(a))
	mb := make(map[int]bool, len(b))
	for _, v := range a {
		ma[v] = true
	}
	for _, v := range b {
		mb[v] = true
	}

   for _, x := range a {
       for _, y := range b {
           s := x + y
           if !ma[s] && !mb[s] {
               fmt.Println(x, y)
               return
           }
       }
   }
}
