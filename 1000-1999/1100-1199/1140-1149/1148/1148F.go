package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "math/rand"
   "os"
   "time"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   vals := make([]int64, n)
   masks := make([]uint64, n)
   var S int64
   for i := 0; i < n; i++ {
       var v int64
       var m uint64
       fmt.Fscan(reader, &v, &m)
       vals[i] = v
       masks[i] = m
       S += v
   }
   rand.Seed(time.Now().UnixNano())
   var s uint64
   for {
       s = rand.Uint64() & ((1<<62)-1)
       if s == 0 {
           s = 1
       }
       var T int64
       for i := 0; i < n; i++ {
           if bits.OnesCount64(s & masks[i])&1 == 1 {
               T += vals[i]
           }
       }
       // Check if flipping yields opposite sign and non-zero
       if S > 0 {
           if 2*T > S && 2*T != S {
               fmt.Println(s)
               return
           }
       } else {
           if 2*T < S && 2*T != S {
               fmt.Println(s)
               return
           }
       }
   }
