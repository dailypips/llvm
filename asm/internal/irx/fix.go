package irx

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/instruction"
	"github.com/llir/llvm/ir/value"
	"github.com/mewkiz/pkg/errutil"
)

// === [ Modules ] =============================================================

// dummyMap maps local variable names to their corresponding values.
type dummyMap map[string]value.Value

// set maps name to the given value.
func (m dummyMap) set(name string, val value.Value) {
	if old, ok := m[name]; ok {
		panic(fmt.Sprintf("mapping for %q already present; old value %v, new value %v", name, old, val))
	}
	m[name] = val
}

// get returns the value for the given name.
func (m dummyMap) get(name string) value.Value {
	val, ok := m[name]
	if !ok {
		panic(fmt.Sprintf("unable to locate mapping for %q", name))
	}
	return val
}

// fixModule replaces dummy values within the given module with their
// corresponding local variables.
func fixModule(module *ir.Module) *ir.Module {
	for i, oldFunc := range module.Funcs {
		// Allocate function specific mapping from local variables names to their
		// corresponding values.
		m := make(dummyMap)
		f := m.fixFunc(oldFunc)
		module.Funcs[i] = f
	}
	return module
}

// === [ Functions ] ===========================================================

// fixFunc replaces dummy values within the given function with their
// corresponding local variables.
func (m dummyMap) fixFunc(oldFunc *ir.Function) *ir.Function {
	f := ir.NewFunction(oldFunc.Name(), oldFunc.Sig())
	for _, oldBlock := range oldFunc.Blocks() {
		block := m.fixBlock(oldBlock)
		f.AppendBlock(block)
	}
	if err := f.AssignIDs(); err != nil {
		panic(errutil.Err(err))
	}
	return f
}

// === [ Basic blocks ] ========================================================

// fixBlock replaces dummy values within the given basic block with their
// corresponding local variables.
func (m dummyMap) fixBlock(oldBlock *ir.BasicBlock) *ir.BasicBlock {
	block := ir.NewBasicBlock(oldBlock.Name())
	for _, oldInst := range oldBlock.Insts() {
		inst := m.fixInst(oldInst)
		block.AppendInst(inst)
	}
	term := m.fixTerm(oldBlock.Term())
	block.SetTerm(term)
	return block
}

// === [ Instructions ] ========================================================

// fixInst replaces dummy values within the given instruction with their
// corresponding local variables.
func (m dummyMap) fixInst(oldInst instruction.Instruction) instruction.Instruction {
	switch oldInst := oldInst.(type) {
	case *instruction.LocalVarDef:
		return m.fixLocalVarDefInst(oldInst)
	case *instruction.Store:
		return m.fixStoreInst(oldInst)
	case *instruction.Fence:
		return m.fixFenceInst(oldInst)
	default:
		panic(fmt.Sprintf("support for instruction type %T not yet implemented", oldInst))
	}
}

// fixLocalVarDefInst replaces dummy values within the given LocalVarDef
// instruction with their corresponding local variables.
func (m dummyMap) fixLocalVarDefInst(oldInst *instruction.LocalVarDef) *instruction.LocalVarDef {
	name := oldInst.Name()
	oldValInst := oldInst.ValInst()
	valInst := m.fixValueInst(oldValInst)
	inst, err := instruction.NewLocalVarDef(name, valInst)
	if err != nil {
		panic(errutil.Err(err))
	}
	m.set(name, inst)
	return inst
}

// fixStoreInst replaces dummy values within the given Store instruction with
// their corresponding local variables.
func (m dummyMap) fixStoreInst(oldInst *instruction.Store) *instruction.Store {
	oldVal := oldInst.Val()
	val := m.fixValue(oldVal)
	oldAddr := oldInst.Addr()
	addr := m.fixValue(oldAddr)
	inst, err := instruction.NewStore(val, addr)
	if err != nil {
		panic(errutil.Err(err))
	}
	return inst
}

// fixFenceInst replaces dummy values within the given Fence instruction with
// their corresponding local variables.
func (m dummyMap) fixFenceInst(oldInst *instruction.Fence) *instruction.Fence {
	panic("irx.dummyMap.fixFenceInst: not yet implemented")
}

// === [ Value instructions ] ==================================================

// fixValueInst replaces dummy values within the given value instruction with
// their corresponding local variables.
func (m dummyMap) fixValueInst(oldValInst instruction.ValueInst) instruction.ValueInst {
	switch oldValInst := oldValInst.(type) {
	// Binary Operations
	case *instruction.Add:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewAdd(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst
	case *instruction.FAdd:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewFAdd(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst
	case *instruction.Sub:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewSub(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst
	case *instruction.FSub:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewFSub(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst
	case *instruction.Mul:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewMul(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst
	case *instruction.FMul:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewFMul(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst
	case *instruction.UDiv:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewUDiv(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst
	case *instruction.SDiv:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewSDiv(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst
	case *instruction.FDiv:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewFDiv(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst
	case *instruction.URem:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewURem(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst
	case *instruction.SRem:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewSRem(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst
	case *instruction.FRem:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewFRem(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst

	// Bitwise Binary Operations
	case *instruction.ShL:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewShL(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst
	case *instruction.LShR:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewLShR(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst
	case *instruction.AShR:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewAShR(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst
	case *instruction.And:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewAnd(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst
	case *instruction.Or:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewOr(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst
	case *instruction.Xor:
		x := m.fixValue(oldValInst.X())
		y := m.fixValue(oldValInst.Y())
		inst, err := instruction.NewXor(x, y)
		if err != nil {
			panic(errutil.Err(err))
		}
		return inst

	// Vector Operations
	case *instruction.ExtractElement:
		panic("irx.dummyMap.fixValueInst: ExtractElement not yet implemented")
	case *instruction.InsertElement:
		panic("irx.dummyMap.fixValueInst: InsertElement not yet implemented")
	case *instruction.ShuffleVector:
		panic("irx.dummyMap.fixValueInst: ShuffleVector not yet implemented")

	// Aggregate Operations
	case *instruction.ExtractValue:
		panic("irx.dummyMap.fixValueInst: ExtractValue not yet implemented")
	case *instruction.InsertValue:
		panic("irx.dummyMap.fixValueInst: InsertValue not yet implemented")

	// Memory Access and Addressing Operations
	case *instruction.Alloca:
		panic("irx.dummyMap.fixValueInst: Alloca not yet implemented")
	case *instruction.Load:
		panic("irx.dummyMap.fixValueInst: Load not yet implemented")
	case *instruction.CmpXchg:
		panic("irx.dummyMap.fixValueInst: CmpXchg not yet implemented")
	case *instruction.AtomicRMW:
		panic("irx.dummyMap.fixValueInst: AtomicRMW not yet implemented")
	case *instruction.GetElementPtr:
		panic("irx.dummyMap.fixValueInst: GetElementPtr not yet implemented")

	// Conversion Operations
	case *instruction.Trunc:
		panic("irx.dummyMap.fixValueInst: Trunc not yet implemented")
	case *instruction.ZExt:
		panic("irx.dummyMap.fixValueInst: ZExt not yet implemented")
	case *instruction.SExt:
		panic("irx.dummyMap.fixValueInst: SExt not yet implemented")
	case *instruction.FPTrunc:
		panic("irx.dummyMap.fixValueInst: FPTrunc not yet implemented")
	case *instruction.FPExt:
		panic("irx.dummyMap.fixValueInst: FPExt not yet implemented")
	case *instruction.FPToUI:
		panic("irx.dummyMap.fixValueInst: FPToUI not yet implemented")
	case *instruction.FPToSI:
		panic("irx.dummyMap.fixValueInst: FPToSI not yet implemented")
	case *instruction.UIToFP:
		panic("irx.dummyMap.fixValueInst: UIToFP not yet implemented")
	case *instruction.SIToFP:
		panic("irx.dummyMap.fixValueInst: SIToFP not yet implemented")
	case *instruction.PtrToInt:
		panic("irx.dummyMap.fixValueInst: PtrToInt not yet implemented")
	case *instruction.IntToPtr:
		panic("irx.dummyMap.fixValueInst: IntToPtr not yet implemented")
	case *instruction.BitCast:
		panic("irx.dummyMap.fixValueInst: BitCast not yet implemented")
	case *instruction.AddrSpaceCast:
		panic("irx.dummyMap.fixValueInst: AddrSpaceCast not yet implemented")

	// Other Operations
	case *instruction.ICmp:
		panic("irx.dummyMap.fixValueInst: ICmp not yet implemented")
	case *instruction.FCmp:
		panic("irx.dummyMap.fixValueInst: FCmp not yet implemented")
	case *instruction.PHI:
		panic("irx.dummyMap.fixValueInst: PHI not yet implemented")
	case *instruction.Select:
		panic("irx.dummyMap.fixValueInst: Select not yet implemented")
	case *instruction.Call:
		panic("irx.dummyMap.fixValueInst: Call not yet implemented")
	case *instruction.VAArg:
		panic("irx.dummyMap.fixValueInst: VAArg not yet implemented")
	case *instruction.LandingPad:
		panic("irx.dummyMap.fixValueInst: LandingPad not yet implemented")
	default:
		panic("irx.dummyMap.fixValueInst: not yet implemented")
	}
}

// --- [ Binary Operations ] ---------------------------------------------------

//    *Add
//    *FAdd
//    *Sub
//    *FSub
//    *Mul
//    *FMul
//    *UDiv
//    *SDiv
//    *FDiv
//    *URem
//    *SRem
//    *FRem

// --- [ Bitwise Binary Operations ] -------------------------------------------

//    *ShL
//    *LShR
//    *AShR
//    *And
//    *Or
//    *Xor

// --- [ Vector Operations ] ---------------------------------------------------

//    *ExtractElement
//    *InsertElement
//    *ShuffleVector

// --- [ Aggregate Operations ] ------------------------------------------------

//    *ExtractValue
//    *InsertValue

// --- [ Memory Access and Addressing Operations ] -----------------------------

//    *Alloca
//    *Load
//    *CmpXchg
//    *AtomicRMW
//    *GetElementPtr

// --- [ Conversion Operations ] -----------------------------------------------

//    *Trunc
//    *ZExt
//    *SExt
//    *FPTrunc
//    *FPExt
//    *FPToUI
//    *FPToSI
//    *UIToFP
//    *SIToFP
//    *PtrToInt
//    *IntToPtr
//    *BitCast
//    *AddrSpaceCast

// --- [ Other Operations ] ----------------------------------------------------

//    *ICmp
//    *FCmp
//    *PHI
//    *Select
//    *Call
//    *VAArg
//    *LandingPad

// === [ Terminators ] =========================================================

//    *Ret
//    *Jmp
//    *Br
//    *Switch
//    *IndirectBr
//    *Invoke
//    *Resume
//    *Unreachable

// fixTerm replaces dummy values within the given terminator with their
// corresponding local variables.
func (m dummyMap) fixTerm(oldTerm instruction.Terminator) instruction.Terminator {
	switch oldTerm := oldTerm.(type) {
	case *instruction.Ret:
		oldVal := oldTerm.Value()
		var val value.Value
		if oldVal != nil {
			val = m.fixValue(oldVal)
		}
		term, err := instruction.NewRet(val)
		if err != nil {
			panic(errutil.Err(err))
		}
		return term
	case *instruction.Jmp:
		panic("irx.dummyMap.fixTerm: Jmp not yet implemented")
	case *instruction.Br:
		panic("irx.dummyMap.fixTerm: Br not yet implemented")
	case *instruction.Switch:
		panic("irx.dummyMap.fixTerm: Switch not yet implemented")
	case *instruction.IndirectBr:
		panic("irx.dummyMap.fixTerm: IndirectBr not yet implemented")
	case *instruction.Invoke:
		panic("irx.dummyMap.fixTerm: Invoke not yet implemented")
	case *instruction.Resume:
		panic("irx.dummyMap.fixTerm: Resume not yet implemented")
	case *instruction.Unreachable:
		panic("irx.dummyMap.fixTerm: Unreachable not yet implemented")
	default:
		panic(fmt.Sprintf("support for terminator type %T not yet implemented", oldTerm))
	}
}

// === [ Values ] ==============================================================

// fixValue replaces dummy values within the given value with their
// corresponding local variables.
func (m dummyMap) fixValue(oldVal value.Value) value.Value {
	switch oldVal := oldVal.(type) {
	case *LocalDummy:
		return m.get(oldVal.name)
	case constant.Constant:
		return oldVal
	default:
		panic(fmt.Sprintf("support for value type %T not yet implemented", oldVal))
	}
}
