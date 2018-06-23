/*
***********************************************************
wendicka_run.go
This particular file has been released in the public domain
and is therefore free of any restriction. You are allowed
to credit me as the original author, but this is not 
required.
This file was setup/modified in: 
2018
If the law of your country does not support the concept
of a product being released in the public domain, while
the original author is still alive, or if his death was
not longer than 70 years ago, you can deem this file
"(c) Jeroen Broks - licensed under the CC0 License",
with basically comes down to the same lack of
restriction the public domain offers. (YAY!)
*********************************************************** 
Version 18.06.23
*/
package main

import _"wendicka/compiler"
import "wendicka/interpreter"
import "os"
import "fmt"

func main(){
  if len(os.Args)<2 {
    fmt.Println("usage: wendicka run <script1> [<script2> [<script3> [...]]]")
    fmt.Println("\nThis is a very minimalist interpreter tool, meant for quick prototyping and testing only")
    fmt.Println("Yes, it should technically be able to do all the CORE features of Wendicka, but it only has a few core APIs in the background, so don't expect much of this tool.")
    os.Exit(0)
  }

  for i:=1;i<len(os.Args);i++{
    w,e:=wi.File2VM(os.Args[i])
    if e!=nil {
      fmt.Println("ERROR!")
      fmt.Println("\t"+e.Error())
      os.Exit(1)
    }
    ok,err:=w.Call("MAIN")
    if !ok {
		fmt.Println("Runtime error:")
		fmt.Println("\t"+err)
		os.Exit(2)
	}
  }
}
