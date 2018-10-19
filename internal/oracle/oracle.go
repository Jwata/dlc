package oracle

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil/hdkeychain"
)

// TimeFormat is a format of settlement time
const TimeFormat = "20060102"

// Oracle is a struct
type Oracle struct {
	name     string                  // display name
	nRpoints int                     // number of commited R-points
	extKey   *hdkeychain.ExtendedKey // extended key
}

// New creates a oracle
func New(name string, params chaincfg.Params, nRpoints int) (*Oracle, error) {
	if isMainNet(params) {
		return nil, fmt.Errorf("mainnet isn't supported yet")
	}

	extKey, err := randomExtKey(name, params)
	if err != nil {
		return nil, err
	}

	// TODO: define path for oracle's HD keys
	// See also bip44, bip47

	oracle := &Oracle{name: name, nRpoints: nRpoints, extKey: extKey}
	return oracle, nil
}

func isMainNet(params chaincfg.Params) bool {
	return params.Net == chaincfg.MainNetParams.Net
}

// randomExtKey creates oracle's random master key
func randomExtKey(name string, params chaincfg.Params) (*hdkeychain.ExtendedKey, error) {
	// TODO: add random logic
	seed := chainhash.DoubleHashB([]byte(name))
	return hdkeychain.NewMaster(seed, &params)
}