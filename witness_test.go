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

func TestScriptEngine(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)

	_, pub := randKeys()
	txVersion := int32(2)
	tx := wire.NewMsgTx(txVersion)

	// prepare fake txin
	prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
	txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0}, nil)
	tx.AddTxIn(txIn)
	pkScript, _ := P2WPKHpkScript(pub)
	txOut := wire.NewTxOut(100000000, pkScript)
	tx.AddTxOut(txOut)
	txHash := tx.TxHash()

	redeemTx := wire.NewMsgTx(txVersion)
	prevOut = wire.NewOutPoint(&txHash, 0)
	txIn = wire.NewTxIn(prevOut, nil, nil)
	redeemTx.AddTxIn(txIn)
	txOut = wire.NewTxOut(0, nil)
	redeemTx.AddTxOut(txOut)

	flags := txscript.ScriptBip16 | txscript.ScriptVerifyWitness

	vm, err := txscript.NewEngine(pkScript, redeemTx, 0, flags, nil, nil, -1)
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
