package dlc

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
)

// Deal contains information about the distributed amounts, commitment messages, and signatures of fixed messages
type Deal struct {
	amts          map[Contractor]btcutil.Amount
	msgs          [][]byte
	msgCommitment *btcec.PublicKey
	msgSign       []byte // oracle's message sign
	cpSign        []byte // conterparty's sign
}

// NewDeal creates a new deal
func NewDeal(amt1, amt2 btcutil.Amount, msgs [][]byte) *Deal {
	amts := make(map[Contractor]btcutil.Amount)
	amts[FirstParty] = amt1
	amts[SecondParty] = amt2
	return &Deal{
		amts: amts,
		msgs: msgs,
	}
}

// AddDeal adds a deal to DLC
func (b *Builder) AddDeal(deal *Deal) int {
	b.dlc.deals = append(b.dlc.deals, deal)
	return len(b.dlc.deals) - 1
}

// Deal gets a deal by id
func (d *DLC) Deal(idx int) (*Deal, error) {
	if len(d.deals) < idx+1 {
		errmsg := fmt.Sprintf("Invalid deal id. id: %d", idx)
		return nil, errors.New(errmsg)
	}

	deal := d.deals[idx]
	return deal, nil
}

// FixedDeal returns a fixed deal
func (d *DLC) FixedDeal() (*Deal, error) {
	for _, deal := range d.deals {
		if deal.msgSign != nil {
			return deal, nil
		}
	}
	return nil, errors.New("no fixed deal found")
}

// SetMsgCommitmentToDeal sets a message commitment received from oracle
func (b *Builder) SetMsgCommitmentToDeal(
	idx int, mc *btcec.PublicKey) error {

	d, err := b.dlc.Deal(idx)
	if err != nil {
		return err
	}

	d.msgCommitment = mc
	return nil
}

// FixDeal fixes a deal by setting the signature provided by oracle
func (b *Builder) FixDeal(idx int, sign []byte) error {
	d, err := b.dlc.Deal(idx)
	if err != nil {
		return err
	}

	// TODO: verify the given sign

	d.msgSign = sign
	return nil
}