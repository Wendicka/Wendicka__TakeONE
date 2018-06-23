/*
  vmdata.go
  
  version: 18.06.23
  Copyright (C) 2018 Jeroen P. Broks
  This software is provided 'as-is', without any express or implied
  warranty.  In no event will the authors be held liable for any damages
  arising from the use of this software.
  Permission is granted to anyone to use this software for any purpose,
  including commercial applications, and to alter it and redistribute it
  freely, subject to the following restrictions:
  1. The origin of this software must not be misrepresented; you must not
     claim that you wrote the original software. If you use this software
     in a product, an acknowledgment in the product documentation would be
     appreciated but is not required.
  2. Altered source versions must be plainly marked as such, and must not be
     misrepresented as being the original software.
  3. This notice may not be removed or altered from any source distribution.
*/
package wi

import "fmt"

type apient struct {
  f (func(v *VM) (bool,string))
}

type API struct{
 fs *(map[string] *apient)
}

func (self *API) Init(){
   self.fs= &map[string] *apient{}
   self.minapi()
}

func (self *API) Register(name string,f (func(v *VM) (bool,string))) {
  q:=self.fs
  (*q)[name] = &apient{}
  qf:=(*q)[name]
  qf.f=f
}

type tIdentifier struct {
  itype string
  vint int64
  vfloat float64
  vstring string
  vchunk string // Also used to point to tables and apis
}

func str2identifier(s string) *tIdentifier{
  r:=&tIdentifier {}
  r.itype = "string"
  r.vstring = s
  return r
}

type identifiermap struct {
    i map[string] *tIdentifier
}

type argquery struct {
  i [] *tIdentifier
}


type tchunkins struct {
  ins byte // refers to instructions
  param [][]byte
}

type tcall struct{
  pos int
  chunk string
  achunk *tchunk
  params *argquery
  returns *argquery
  ended bool
}

type tchunk struct {
  instruction []*tchunkins
  labels map[string] int
}

type VM struct{
    //chunks map[string] []byte
    chunks map[string] *tchunk
    vapi *API
    identifiers *identifiermap
    calls [] *tcall
    ccall int
    nextcallparam *argquery
    lastcallreturn *argquery
    lastcompare bool
}

func (self *VM) callChunk(chunk string) bool{
   if MaxCallStack>0 && len(self.calls)>=MaxCallStack {
	   wError("Call stack has reached the limit!")
	   return false
   }
   ncall:=&tcall{}
   ncall.pos=0
   ncall.chunk=chunk
   ncall.params=self.nextcallparam
   ncall.returns=&argquery{}
   self.calls = append(self.calls,ncall)
   self.ccall = len(self.calls)-1
   achunk,ok:=self.chunks[chunk]
   if ok { ncall.achunk=achunk}
   return ok
}

func (self *VM) Call(chunk string) (bool,string){
  if !self.callChunk(chunk) { return false,"Call to chunk "+chunk+" failed!"; }
  rsuccess:=true
  err:=""
  for len(self.calls)>0{
	  nc:=self.calls[len(self.calls)-1]
      if len(nc.achunk.instruction)<=nc.pos { panic(fmt.Sprintf("Position past chunk end: %d of %d",nc.pos,len(nc.achunk.instruction)))}
      insl:=nc.achunk.instruction[nc.pos]
      insn:=insl.ins
      insa:=insl.param
      insd,insf:=winstructs[insn]
      if !insf { panic( fmt.Sprintf("FATAL INTERNAL ERROR! Unknown instruction code in execution %X/%d\nPlease report",insn,insn));}
      chat("Executing instruction: ",fmt.Sprintf("%d",insn))
      insd.do_it(self,insa)
      if !nc.ended{
		nc.pos++
		if nc.pos>=len(nc.achunk.instruction) { return false,"Chunk not properly ended";}
	  }
  }
  return rsuccess,err
}

func (self *VM) callAPI(chunk string) (bool,string){
  rsuccess:=true
  err:=""
  a:=self.vapi.fs
  f,ok:=(*a)[chunk]
  ncall:=&tcall{}
  ncall.chunk=chunk
  ncall.params=self.nextcallparam
  ncall.returns=&argquery{}
  self.calls = append(self.calls,ncall)
  self.ccall = len(self.calls)-1
  if ncall.params==nil { chat("WARNING!","Param field is nil")}
  if !ok { return false,"Unknown API called: "+chunk;}
  if f.f==nil { return false,"API nilfunc in "+chunk}
  rsuccess,err = f.f(self)
  self.calls = self.calls[:len(self.calls)-1]
  self.ccall = len(self.calls)-1
  self.lastcallreturn = ncall.returns
  self.nextcallparam = &argquery{}
  return rsuccess,err
}


func (self *VM) Init(a *API){
  if a==nil {
    self.vapi = &API{}
    self.vapi.Init()
    self.vapi.minapi() // Minimal features all APIs MUST have.
  } else {
    self.vapi = a
  }
  ti:=map[string] *tIdentifier {}
  self.identifiers = &identifiermap{}
  self.identifiers.i = ti
  self.chunks = map[string] *tchunk {}
  self.calls = []*tcall{}
  self.nextcallparam = &argquery{}
  self.lastcallreturn = &argquery{}
}

// Normally every VM creates its own API register, but if you have multiple VMS all using the same API, why bother?
func (self *VM) TieAPI(api *API) {
  self.vapi=api
}

func (self *VM) Register(name string, f (func(v *VM) (bool,string))){
  self.vapi.Register(name,f)
}

// Creates a new VM for Wendicka
// If you leave a *API to nil, one will be generated
func NewVM(a *API) *VM{
  ret:=&VM{}
  ret.Init(a)
  return ret
}
