package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   const maxN = 20005
   isPrime := make([]bool, maxN)
   for i := 2; i < maxN; i++ {
       isPrime[i] = true
   }
   for i := 2; i*i < maxN; i++ {
       if isPrime[i] {
           for j := i * i; j < maxN; j += i {
               isPrime[j] = false
           }
       }
   }
   a, b := 0, 0
   for i := n; i >= 2; i-- {
       if isPrime[i] {
           a = i
           break
       }
   }
   for i := n; i < maxN; i++ {
       if isPrime[i] {
           b = i
           break
       }
   }
   fmt.Println(a, b)
}
