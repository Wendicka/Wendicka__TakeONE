package wi

type apient struct {
  f (func(v *VM))
}

type API struct{
 fs *(map[string] *apient)
}

func (self *API) Init(){
   self.fs= &map[string] *apient{}
   self.minapi()
}

func (self *API) Register(name string,f (func(v *VM))) {
  q:=self.fs
  (*q)[name] = &apient{}

}

type tIdentifier struct {
  itype string
  vint int64
  vfloat float64
  vstring string
  vchunk string // Also used to point to tables and apis
}



type identifiermap struct {
    i map[string] *tIdentifier
}

type argquery struct {
  i [] *tIdentifier
}

type tcall struct{
  pos int
  chunk string
  params *argquery
  returns *argquery
}

type tchunkins struct {
  ins byte // refers to instructions
  param [][]byte
}

type tchunk struct {
  instruction []*tchunkins
}

type VM struct{
    //chunks map[string] []byte
    chunks map[string] tchunk
    vapi *API
    identifiers *identifiermap
    calls [] *tcall
    ccall int
    nextcallparam *argquery
    lastcallreturn *argquery
}

func (self *VM) Init(a *API){
  if a==nil {
    self.vapi = &API{}
    self.vapi.minapi() // Minimal features all APIs MUST have.
  } else {
    self.vapi = a
  }
  ti:=map[string] *tIdentifier {}
  self.identifiers = &identifiermap{}
  self.identifiers.i = ti
}

// Normally every VM creates its own API register, but if you have multiple VMS all using the same API, why bother?
func (self *VM) TieAPI(api *API) {
  self.vapi=api
}

func (self *VM) Register(name string, f (func(v *VM))){
  self.vapi.Register(name,f)
}

// Creates a new VM for Wendicka
// If you leave a *API to nil, one will be generated
func NewVM(a *API) *VM{
  ret:=&VM{}
  ret.Init(a)
  return ret
}
