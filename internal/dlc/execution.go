package dlc

import (
	"fmt"

	"github.com/btcsuite/btcd/wire"
	"github.com/dgarage/dlc/internal/script"
)

// ContractExecutionTx constracts a contract execution tx using pubkeys and given condition.
// Both parties have different transactions signed by the other side.
// input:
//   [0]:fund transaction output[0]
// output:
//   [0]:settlement script
//   [1]:p2wpkh (option)
func (d *DLC) ContractExecutionTx() (*wire.MsgTx, error) {

	// tx := wire.NewMsgTx(2)
	// txid := d.FundTx().TxHash()
	// tx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(&txid, 0), nil, nil))
	tx, err := d.newRedeemTx()
	if err != nil {
		return nil, err
	}

	// TODO: switching script
	// var val1 int64
	// var val2 int64
	// var pub1 *btcec.PublicKey
	// var pub2 *btcec.PublicKey
	// if isA {
	// 	val1 = rate.amta
	// 	val2 = rate.amtb
	// 	pub1 = d.puba
	// 	pub2 = d.pubb
	// } else {
	// 	val1 = rate.amtb
	// 	val2 = rate.amta
	// 	pub1 = d.pubb
	// 	pub2 = d.puba
	// }
	// if val1 <= 0 {
	// 	return nil
	// }
	// TODO: Add first party's pubkey to oracle's pubkey
	// pub := &btcec.PublicKey{}
	// pub.X, pub.Y = btcec.S256().Add(rate.key.X, rate.key.Y, pub1.X, pub1.Y)

	// TODO: create script
	sc, err := script.ContractExecutionScript(pub, pub2)
	if err != nil {
		return nil, err
	}
	pkScript, err := script.P2WSHpkScript(sc)
	if err != nil {
		return nil, err
	}
	fmt.Println(pkScript)

	// TODO: set txout
	// txout1 := wire.NewTxOut(val1, pkScript)
	// tx.AddTxOut(txout1)
	// if val2 > 0 {
	// 	txout2 := wire.NewTxOut(val2, P2WPKHpkScript(pub2))
	// 	tx.AddTxOut(txout2)
	// }
	return tx, nil
}
