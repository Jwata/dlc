package main

import (
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/stretchr/testify/assert"
)

func TestP2WPKHpkScript(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)

	pri, pub := randKeys()
	txVersion := int32(2)
	tx := wire.NewMsgTx(txVersion)

	// script
	pkScript, _ := P2WPKHpkScript(pub)

	// prepare fake txin
	prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
	txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0}, nil)
	tx.AddTxIn(txIn)
	utxoAmt := int64(100000000)
	tx.AddTxOut(wire.NewTxOut(utxoAmt, pkScript))
	txHash := tx.TxHash()

	redeemTx := wire.NewMsgTx(txVersion)
	prevOut = wire.NewOutPoint(&txHash, 0)
	redeemTx.AddTxIn(wire.NewTxIn(prevOut, nil, nil))
	redeemTx.AddTxOut(wire.NewTxOut(0, nil))

	sighash := txscript.NewTxSigHashes(redeemTx)
	sign, err := txscript.RawTxInWitnessSignature(
		redeemTx, sighash, 0, utxoAmt, tx.TxOut[0].PkScript, txscript.SigHashAll, pri)
	assert.Nil(err)
	tw := wire.TxWitness{}
	tw = append(tw, sign)
	tw = append(tw, pub.SerializeCompressed())
	redeemTx.TxIn[0].Witness = tw

	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(pkScript, redeemTx, 0, flags, nil, nil, utxoAmt)
	assert.Nil(err)

	err = vm.Execute()
	assert.Nil(err)
}

func randKeys() (*btcec.PrivateKey, *btcec.PublicKey) {
	seed, _ := hdkeychain.GenerateSeed(hdkeychain.MinSeedBytes)
	extKey, _ := hdkeychain.NewMaster(seed, &chaincfg.RegressionNetParams)
	pub, _ := extKey.ECPubKey()
	priv, _ := extKey.ECPrivKey()
	return priv, pub
}

// P2WPKHpkScript creates pk script is OP_0 + HASH160(<public key>)
func P2WPKHpkScript(pub *btcec.PublicKey) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_0)
	builder.AddData(btcutil.Hash160(pub.SerializeCompressed()))
	return builder.Script()
}
