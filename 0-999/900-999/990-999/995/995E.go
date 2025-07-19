package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "time"
)

// random returns a random int64 in [0, p)
func random(p int64) int64 {
   return rand.Int63() % p
}

// step computes the total number of operations for given a, b
func step(a, b int64) int {
   ret := 0
   for b > 0 {
       ret += int(a/b + 1)
       tmp := a % b
       a = b
       b = tmp
   }
   return ret
}

// solve generates a sequence of operations to reduce (k*u%p, k) to (0, _)
// encoded as 0 for subtraction and 1 for swap
func solve(u, p int64) []int {
   for {
       k := random(p)
       if k == 0 {
           continue
       }
       a := (k * u) % p
       b := k
       if step(a, b) >= 100 {
           continue
       }
       var ret []int
       for a > 0 {
           if a >= b {
               ret = append(ret, 0)
               a -= b
           } else {
               ret = append(ret, 1)
               a, b = b, a
           }
       }
       return ret
   }
}

func main() {
   rand.Seed(time.Now().UnixNano())
   var u, v, p int64
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &u, &v, &p)
   a1 := solve(u, p)
   a2 := solve(v, p)
   // reverse a2
   for i, j := 0, len(a2)-1; i < j; i, j = i+1, j-1 {
       a2[i], a2[j] = a2[j], a2[i]
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, len(a1)+len(a2))
   for _, x := range a1 {
       if x == 1 {
           fmt.Fprint(writer, 3, " ")
       } else {
           fmt.Fprint(writer, 2, " ")
       }
   }
   for _, x := range a2 {
       if x == 1 {
           fmt.Fprint(writer, 3, " ")
       } else {
           fmt.Fprint(writer, 1, " ")
       }
   }
   fmt.Fprintln(writer)
}
