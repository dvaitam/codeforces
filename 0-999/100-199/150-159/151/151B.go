package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   scanner := bufio.NewScanner(reader)
   scanner.Split(bufio.ScanWords)
   next := func() string {
       if scanner.Scan() {
           return scanner.Text()
       }
       return ""
   }

   nStr := next()
   n, err := strconv.Atoi(nStr)
   if err != nil {
       return
   }
   names := make([]string, 0, n)
   taxiCounts := make([]int, 0, n)
   pizzaCounts := make([]int, 0, n)
   girlCounts := make([]int, 0, n)

   for i := 0; i < n; i++ {
       siStr := next()
       si, err := strconv.Atoi(siStr)
       if err != nil {
           si = 0
       }
       name := next()
       names = append(names, name)
       tCnt, pCnt, gCnt := 0, 0, 0
       for j := 0; j < si; j++ {
           num := next()
           // extract digits
           var d [6]int
           idx := 0
           for _, c := range num {
               if c >= '0' && c <= '9' {
                   if idx < 6 {
                       d[idx] = int(c - '0')
                       idx++
                   }
               }
           }
           // check taxi: all equal
           allEqual := true
           for k := 1; k < 6; k++ {
               if d[k] != d[0] {
                   allEqual = false
                   break
               }
           }
           if allEqual {
               tCnt++
               continue
           }
           // check pizza: strictly decreasing
           dec := true
           for k := 1; k < 6; k++ {
               if d[k] >= d[k-1] {
                   dec = false
                   break
               }
           }
           if dec {
               pCnt++
           } else {
               gCnt++
           }
       }
       taxiCounts = append(taxiCounts, tCnt)
       pizzaCounts = append(pizzaCounts, pCnt)
       girlCounts = append(girlCounts, gCnt)
   }

   maxTaxi, maxPizza, maxGirl := 0, 0, 0
   for i := 0; i < n; i++ {
       if taxiCounts[i] > maxTaxi {
           maxTaxi = taxiCounts[i]
       }
       if pizzaCounts[i] > maxPizza {
           maxPizza = pizzaCounts[i]
       }
       if girlCounts[i] > maxGirl {
           maxGirl = girlCounts[i]
       }
   }

   var taxiNames, pizzaNames, girlNames []string
   for i := 0; i < n; i++ {
       if taxiCounts[i] == maxTaxi {
           taxiNames = append(taxiNames, names[i])
       }
       if pizzaCounts[i] == maxPizza {
           pizzaNames = append(pizzaNames, names[i])
       }
       if girlCounts[i] == maxGirl {
           girlNames = append(girlNames, names[i])
       }
   }

   fmt.Printf("If you want to call a taxi, you should call: %s.\n", strings.Join(taxiNames, ", "))
   fmt.Printf("If you want to order a pizza, you should call: %s.\n", strings.Join(pizzaNames, ", "))
   fmt.Printf("If you want to go to a cafe with a wonderful girl, you should call: %s.\n", strings.Join(girlNames, ", "))
}
