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

const TxVersion = int32(2)

func TestP2WPKHpkScript(t *testing.T) {
	var err error
	assert := assert.New(t)

	pri, pub := randKeys()
	amt := int64(10000)
	pkScript, _ := P2WPKHpkScript(pub)

	// prepare source transaction
	sourceTx := testSourceTx()

	// append P2WPKH script
	witOut := wire.NewTxOut(amt, pkScript)
	sourceTx.AddTxOut(witOut)

	// create redeem tx from source tx
	redeemTx := createRedeemTx(sourceTx)

	// append witness signature to redeem tx
	redeemTx, err = appendWitnessSignature(witOut, redeemTx, pri, pub)
	assert.Nil(err)

	// execute script
	err = executeScript(pkScript, redeemTx, amt)
	assert.Nil(err)
}

// P2WPKHpkScript creates pk script is OP_0 + HASH160(<public key>)
func P2WPKHpkScript(pub *btcec.PublicKey) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_0)
	builder.AddData(btcutil.Hash160(pub.SerializeCompressed()))
	return builder.Script()
}

// test utilities
func randKeys() (*btcec.PrivateKey, *btcec.PublicKey) {
	seed, _ := hdkeychain.GenerateSeed(hdkeychain.MinSeedBytes)
	extKey, _ := hdkeychain.NewMaster(seed, &chaincfg.RegressionNetParams)
	pub, _ := extKey.ECPubKey()
	priv, _ := extKey.ECPrivKey()
	return priv, pub
}

func testSourceTx() *wire.MsgTx {
	tx := wire.NewMsgTx(TxVersion)
	prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
	txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0}, nil)
	tx.AddTxIn(txIn)
	return tx
}

func createRedeemTx(sourceTx *wire.MsgTx) *wire.MsgTx {
	txHash := sourceTx.TxHash()
	outPt := wire.NewOutPoint(&txHash, 0)

	tx := wire.NewMsgTx(TxVersion)
	tx.AddTxIn(wire.NewTxIn(outPt, nil, nil))

	return tx
}

func appendWitnessSignature(
	txOut *wire.TxOut, tx *wire.MsgTx, pri *btcec.PrivateKey, pub *btcec.PublicKey,
) (*wire.MsgTx, error) {
	sighash := txscript.NewTxSigHashes(tx)
	sign, err := txscript.RawTxInWitnessSignature(
		tx, sighash, 0, txOut.Value, txOut.PkScript, txscript.SigHashAll, pri)
	if err != nil {
		return nil, err
	}
	tx.TxIn[0].Witness = witnessForP2WPKH(sign, pub)

	return tx, nil
}

func witnessForP2WPKH(sign []byte, pub *btcec.PublicKey) [][]byte {
	tw := wire.TxWitness{}
	tw = append(tw, sign)
	tw = append(tw, pub.SerializeCompressed())
	return tw
}

func executeScript(pkScript []byte, tx *wire.MsgTx, amt int64) error {
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(pkScript, tx, 0, flags, nil, nil, amt)
	if err != nil {
		return err
	}

	return vm.Execute()
}
