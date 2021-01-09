package generationk

import (
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Direction int

//negativeBalanceErr is used to flag a negative carry
var negativeBalanceErr = errors.New("Balance < 0")

//Direction of a trade: long or short
const (
	Long Direction = iota
	Short
)

//The portfolio holds assets: holdings, the portfolio holds a mutext to
//be able to use the same portfolio when testing many assets in parallell but
//updating the account on a single portfolio
type Portfolio struct {
	sync.Mutex
	holdings []Holding
	cash     float64
}

//This is what we are owning, a holding
type Holding struct {
	qty       int
	assetName string
	price     float64
	time      time.Time
}

//Is used to create a new portfolio
func NewPortfolio() *Portfolio {
	portfolio := Portfolio{

		holdings: make([]Holding, 0),
		cash:     0,
	}
	return &portfolio
}

//IsOwning is used to find out if a position is already owned in this asset
func (p *Portfolio) IsOwning(assetName string) bool {
	p.Lock()
	defer p.Unlock()

	for k := range p.holdings {
		if p.holdings[k].assetName == assetName {

			/*log.WithFields(log.Fields{
				"AssetName": assetName,
			}).Debug("Already owned")*/
			return true
		}
	}
	return false
}

//Remove a holding, its sold
func (p *Portfolio) RemoveHolding(position Holding) {
	p.Lock()
	defer p.Unlock()
	log.Info().
		Str("asset", position.assetName).
		Time("time", position.time).
		Float64("price", position.price).
		Int("Qty", position.qty).
		Msg("PORTFOLIO> Removing position from portfolio")

	pos := -1

	for k := range p.holdings {
		if position.assetName == p.holdings[k].assetName {
			pos = k
		}
	}

	p.holdings = remove(pos, p.holdings)
}

func remove(ix int, holdings []Holding) []Holding {
	return append(holdings[:ix], holdings[ix+1:]...)
}

//AddHolding, its been bought
func (p *Portfolio) AddHolding(position Holding) {
	log.Info().
		Str("asset", position.assetName).
		Time("time", position.time).
		Int("Qty", position.qty).
		Msg("PORTFOLIO> Adding position to portfolio")

	p.Lock()
	p.holdings = append(p.holdings, position)
	p.Unlock()
}

//checkBalance is used to check the balance before buying
func (p *Portfolio) checkBalance(cost float64) error {
	p.Lock()
	balance := p.cash + cost
	p.Unlock()
	if balance < 0 {
		return negativeBalanceErr
	}
	return nil
}

//addToBalance is used to add to the account after selling with profit
func (p *Portfolio) addToBalance(value float64) {
	p.Lock()
	p.cash += value
	p.Unlock()
}

//SubtractFromBalance is used to decrease the amount on the account
func (p *Portfolio) subtractFromBalance(cost float64) error {
	err := p.checkBalance(cost)
	if err != nil {
		return err
	}

	p.Lock()
	p.cash -= cost
	p.Unlock()

	return nil
}

//SetBalance is used to set the starting balance of the account
func (p *Portfolio) SetBalance(amount float64) {
	p.cash = amount
}

//GetBalance returns in the balance of the account
func (p *Portfolio) GetBalance() float64 {
	return p.cash
}
