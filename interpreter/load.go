package wi

import(
  //"os"
  "fmt"
  "bytes"
  "errors"
  "trickyunits/qff"
)


func Bytes2VM(l []byte) (*VM,error){
  var e error
  ret:=NewVM(nil)
  bt:=bytes.NewReader(l)
  h:=qff.RawReadString(bt,3)
  readchunk:=""
  if h!="WBC" { return nil,errors.New("Given memory bank is NOT Wendicka Byte Code")}
  bits:=qff.ReadInt64(bt)
  if bits!=8 { return nil,errors.New(fmt.Sprintf("This version of the Wendicka interpreter only supports 8 bit instructions and this bytecode contains %d bit instructions",bits))}
  for bt.Len()>0 { //!qff.EOF(*bt) {
    ins:=qff.ReadByte(bt)
    if ins!=0 && readchunk=="" { return nil,errors.New(fmt.Sprintf("Chunkless code: %X",ins))}
    switch ins {
      case 0:
        namechunk:=qff.ReadString(bt)
        if _,exist:=ret.chunks[namechunk];exist { return nil,errors.New("Duplicate chunk: "+namechunk)}
        ret.chunks[namechunk]=tchunk{}
        readchunk=namechunk
      default:
        q,ok:=winstructs[ins]
        if !ok { return nil,errors.New(fmt.Sprintf("Unknown instruction %X",ins))}
        a:=[][]byte {}//*tIdentifier{}
        for i,partype:=range q.needparam {
          switch partype {
            case "string":
              s:= qff.ReadString(bt)
              a = append(a,[]byte(s))
            default:
              return nil,errors.New(fmt.Sprintf("Unknown paramter type for instruction %X, position %d: %s",ins,i,partype))
          }
        }
        cc:=tchunkins{}
        cc.param=a
        cc.ins=ins
    }
  }
  return ret,e
}

func File2VM(filename string) (*VM,error){
  b,e:=qff.EGetFile(filename)
  if e!=nil { return nil,e }
  var r *VM
  r,e=Bytes2VM(b)
  return r,e
}
