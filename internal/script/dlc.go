package script

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
)

// ContractExecutionScript returns a contract execution script.
//
// Script Code:
//  OP_IF
//    <public key A + message pub key>
//  OP_ELSE
//    delay(fix 144)
//    OP_CHECKSEQUENCEVERIFY
//    OP_DROP
//    <public key B>
//  OP_ENDIF
//  OP_CHECKSIG
//
// The if block can be passed when the contractor A has a valid oracle's sign to the message.
// But if the contractor sends this transaction without the oracle's valid sign, the else block will be used by the other party B after the delay time (1 day approximately).
// Please check the original paper for more details.
//
// https://adiabat.github.io/dlc.pdf
func ContractExecutionScript(pub1, pub2 *btcec.PublicKey) ([]byte, error) {
	delay := uint16(144)
	csvflg := uint32(0x00000000)
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_IF)
	builder.AddData(pub1.SerializeCompressed())
	builder.AddOp(txscript.OP_ELSE)
	builder.AddInt64(int64(delay) + int64(csvflg))
	builder.AddOp(txscript.OP_CHECKSEQUENCEVERIFY)
	builder.AddOp(txscript.OP_DROP)
	builder.AddData(pub2.SerializeCompressed())
	builder.AddOp(txscript.OP_ENDIF)
	builder.AddOp(txscript.OP_CHECKSIG)
	return builder.Script()
}
