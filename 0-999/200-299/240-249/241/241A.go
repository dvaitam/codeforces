package main

import (
   "bufio"
   "fmt"
   "os"
)


func main() {
   reader := bufio.NewReader(os.Stdin)
   var m, k int
   if _, err := fmt.Fscan(reader, &m, &k); err != nil {
       return
   }
   d := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &d[i])
   }
   s := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &s[i])
   }
   // Total travel distance and sum of supplies
   totalD := 0
   sumS := 0
   for i := 0; i < m; i++ {
       totalD += d[i]
       sumS += s[i]
   }
   // Find city with maximum supply
   maxS := 0
   maxIdx := 0 // 0-based
   for i := 0; i < m; i++ {
       if s[i] > maxS {
           maxS = s[i]
           maxIdx = i
       }
   }
   // Phase 1: minimal waits to reach city with maxS (index maxIdx)
   fuel := 0
   waitTime := 0
   if m > 0 {
       fuel = s[0]
   }
   // traverse roads 0..maxIdx-1
   for i := 0; i < maxIdx; i++ {
       need := d[i]
       if fuel < need {
           // wait enough batches at city i (supply s[i])
           // t = ceil((need - fuel) / s[i])
           rem := need - fuel
           t := rem / s[i]
           if rem%s[i] != 0 {
               t++
           }
           waitTime += t * k
           fuel += t * s[i]
       }
       // travel to next city
       fuel = fuel - need + s[i+1]
   }
   // Phase 2: wait at city maxIdx for extra fuel
   extra := totalD - sumS
   if extra < 0 {
       extra = 0
   }
   // batches needed at max city
   t2 := 0
   if extra > 0 {
       t2 = extra / maxS
       if extra%maxS != 0 {
           t2++
       }
   }
   waitTime += t2 * k
   // total time is travel time totalD plus waitTime
   fmt.Println(waitTime + totalD)
}
