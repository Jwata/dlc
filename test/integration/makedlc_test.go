package integration

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/dgarage/dlc/internal/dlc"
	"github.com/stretchr/testify/assert"
)

func TestContractorMakeDLC(t *testing.T) {
	// Given an oracle "Olivia" who publishes random a 5-digit number everday like lottery
	nDigit := 5
	olivia, _ := newOracle("Olivia", nDigit)

	// And next announcement time
	fixingTime := nextLotteryAnnouncement()

	// And a contractor "Alice"
	alice, _ := newContractor("Alice")

	// And a contractor "Bob"
	bob, _ := newContractor("Bob")

	// And Alice and Bob bed on all cases that will be publsed
	contractorsBetOnAllDigitPatters(t, alice, bob, nDigit, fixingTime)

	fmt.Println(olivia)
}

func nextLotteryAnnouncement() time.Time {
	tomorrow := time.Now().AddDate(0, 0, 1)
	year, month, day := tomorrow.Date()
	return time.Date(year, month, day, 12, 0, 0, 0, tomorrow.Location())
}

func contractorsBetOnAllDigitPatters(
	t *testing.T, c1, c2 *Contractor, nDigit int, fixingTime time.Time) {

	var onebtc btcutil.Amount = 1 * btcutil.SatoshiPerBitcoin
	famta, famtb := onebtc, onebtc
	deals := randomDealsForAllDigitPatterns(nDigit, int(famta+famtb))
	conds, err := dlc.NewConditions(fixingTime, famta, famtb, 1, 1, 1, deals)
	assert.NoError(t, err)

	c1.createDLCBuilder(conds, dlc.FirstParty)
	c2.createDLCBuilder(conds, dlc.SecondParty)
}

func randomDealsForAllDigitPatterns(nDigit, famt int) []*dlc.Deal {
	var deals []*dlc.Deal
	for d := 0; d < int(math.Pow10(nDigit)); d++ {
		// randomly decide deals
		damta, damtb := randomDealAmts(famt)
		deal := dlc.NewDeal(damta, damtb, nDigitToBytes(d, nDigit))
		deals = append(deals, deal)
	}
	return deals
}

func randomDealAmts(famt int) (btcutil.Amount, btcutil.Amount) {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	a := r.Intn(famt + 1)
	b := famt - a
	return iToAmount(a), iToAmount(b)
}

func iToAmount(n int) btcutil.Amount {
	return btcutil.Amount(float64(n))
}

// convert a n-digit number to byte
// e.g. 123 -> [][]byte{{1}, {2}, {3}}
func nDigitToBytes(d int, n int) [][]byte {
	b := make([][]byte, n)
	for i := 0; i < n; i++ {
		b[i] = []byte{byte(d % 10)}
		d = d / 10
	}
	return b
}
