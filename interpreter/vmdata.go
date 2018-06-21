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
   ncall:=&tcall{}
   ncall.pos=0
   ncall.chunk=chunk
   ncall.params=self.nextcallparam
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
