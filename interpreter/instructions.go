package wi

import(
  "trickyunits/qint"
  "errors"
)

type winstruct struct {
  do_it func(w *VM,args [][]byte) bool
  needparam []string
}

var winstructs = map[byte] *winstruct {}

func igidentifier(w *VM,b[]byte) (*tIdentifier,error){
  i:=b[0]
  p:=b[1:]
  switch i{
    case 0:
      r:=tIdentifier{}
      r.itype="string"
      r.vstring=string(p)
      return &r,nil
    case 1:
      r:=w.identifiers.i
      return r[string(p)],nil
    case 2:
      var e error
      r:=tIdentifier{}
      r.itype="integer"
      r.vint,e=qint.BytesToInt64(p)
      return &r,e
    case 3:
      var e error
      r:=tIdentifier{}
      r.itype="float"
      r.vfloat,e=qint.BytesToFloat(p)
      return &r,e
    default:
      return nil,errors.New("Unknown Identifier Code. Either the code is compiled incorrectly, or the code is corrupt, or your version of Wendicka is too old for this particular code.")
  }
  return nil,errors.New("I'm completely at a loss. This error could only happen to either a bug or an outdated version of Wendicka.")
}

func init(){
  winstructs[1] = &winstruct{
    // MOV
    func (w *VM,args [][]byte) bool{
      v:=string(args[0])
      if args[0][0]!='$' {
        wError("Invalid MOV! "+v)
        return false
      }
      if len(args)<=1 { wError("No values received to appoint to: "+string(args[0])); return false}
      id,e:=igidentifier(w,args[1])
      if e!=nil { wError(e.Error()); return false;}
      switch v {
        case "$__ARG":
          w.nextcallparam.i = append(w.nextcallparam.i,id)
          chat("Argument ",id.itype,string(args[1])," added for next call")
          return true
        default:
          kut:=&w.identifiers.i
          (*kut)[v]=id
          chat("Value ",string(args[1]),"("+id.itype+")"," assigned to variable ",v)
          return true
        }
    },
    []string{"string","identifier"},
  }

  // CALL
  winstructs[10] = &winstruct{
    func (w *VM,args [][]byte) bool{
          ch:=string(args[0])
          if args[0][0]=='$' {
            id:=w.identifiers.i[ch]
            if id.itype=="chunk" || id.itype=="api" { ch=id.vstring } else {wError("Call identifiers must refer to apis or chunks"); return false}
          }
          afs:=w.vapi.fs
          if _,ok:=w.chunks[ch];ok {
            w.callChunk(ch)
          } else if _,aok:=(*afs)[ch];aok {
            s,e:=w.callAPI(ch)
            if !s {wError("API Call error: "+e)}
          }
      return true
    },
    []string{"string"},
  }

  // END CALL
  winstructs[0xff] = &winstruct{
    func (w *VM,args [][]byte) bool{
		cll:=len(w.calls)-1
		cl:=w.calls[cll]
		cl.ended=true
		w.calls=w.calls[:cll]
      return true
    },
    []string{},

  }
}
