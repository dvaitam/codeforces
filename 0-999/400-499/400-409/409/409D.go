package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var idx int
   if _, err := fmt.Fscan(reader, &idx); err != nil {
       return
   }
   facts := []int{
       8848,  // Mount Everest height in meters
       958,   // participants in largest board game tournament
       12766, // participants in largest online maths competition
       6695,  // Nile length in kilometers
       1100,  // Amazon width in kilometers
       807,   // Angel Falls drop in meters
       31962, // Hotel Everest View height in meters
       146,   // neutrons in common uranium isotope
       -68,   // lowest recorded temperature in Oymyakon
       25,    // length in feet of longest snake held in captivity
       134,   // longest fur on a cat in centimeters
       10000, // hairs per square inch on sea otters
       663268, // area of Alaska in square miles
       154103, // coastline length of Alaska in miles
       1642,   // depth of Lake Baikal in meters
       106,    // colors in Turkmenistan flag
   }
   if idx >= 1 && idx <= len(facts) {
       fmt.Println(facts[idx-1])
   }
}
