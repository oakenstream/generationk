package generationk

import (
	"time"

	log "github.com/sirupsen/logrus"
)

//OrderType is used to describe an order
type OrderType int

const (
	//Buy order
	Buy OrderType = iota
	//Sell order
	Sell
	//SellShort order
	SellShort
	//Cover short order
	Cover
)

//Broker is used to send orders
type Broker struct {
	portfolio Portfolio
	channel   chan Event
}

//PlaceOrder is used to place an order with the broker
func (b *Broker) PlaceOrder(order Order) {
	log.WithFields(log.Fields{
		"ordertype": order.Ordertype,
		"asset":     (*order.Asset).Name,
		"time":      order.Time,
		"amount":    order.Amount,
	}).Debug("BROKER>PLACE BUY ORDER")
	if order.Ordertype == Buy {
		b.buy(order.Asset, order.Time, order.Amount)
	}
}

func (b *Broker) buy(asset *Asset, time time.Time, amount float64) {
	//How many are we buying
	qty := int(amount / asset.Close())
	//b.portfolio.Add(*pos)
	log.WithFields(log.Fields{
		"Amount": amount,
	}).Info("BROKER> FILLED")
	b.channel <- Fill{Qty: qty, AssetName: (*asset).Name, Time: time}
	log.Info("BROKER> Put FILL EVENT in queue")
}

func (b *Broker) sell(asset *Asset, amount int) {

}