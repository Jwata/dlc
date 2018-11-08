package wallet

import (
	"fmt"
	"path/filepath"
)

const (
	defaultConfigFilename = "bitcoin.conf"
	walletDbName          = "wallet.db"
)

// get location of bitcoind folder
// make sure bitcoin.conf file exists
// read parameters from conf file
// pass parameters into wallet?

func loadConfig() {
	configFilePath := filepath.Join(appDataDir, defaultConfigFilename)
	fmt.Println(configFilePath)

}
