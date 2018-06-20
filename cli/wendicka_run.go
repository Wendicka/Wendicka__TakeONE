package main

import _"wendicka/compiler"
import _"wendicka/interpreter"
import "os"
import "fmt"

func main(){
  if len(os.Args)<2 {
    fmt.Println("usage: wendicka run <script1> [<script2> [<script3> [...]]]")
    fmt.Println("\nThis is a very minimalist interpreter tool, meant for quick prototyping and testing only")
    fmt.Println("Yes, it should technically be able to do all the CORE features of Wendicka, but it only has a few core APIs in the background, so don't expect much of this tool.")
    os.Exit(0)
  }
}
