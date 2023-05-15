package models

type CourierType string

const (
	Foot CourierType = "Foot"
	Bike CourierType = "Bike"
	Car  CourierType = "Car"
)

type Order1 struct {
	Weight int
	Region int
	Cost   float64
}

type Courier1 struct {
	CourierType CourierType
	MaxWeight   float32
	MaxOrders   int
	Regions     []int
	Orders      []Order
	TotalWeight float32
	TotalCost   float32
	TotalTime   float32
}

type OrdersByCost []Order

func (o OrdersByCost) Len() int           { return len(o) }
func (o OrdersByCost) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o OrdersByCost) Less(i, j int) bool { return o[i].Cost > o[j].Cost }

func (c *Courier1) CanTakeOrder(order Order) bool {
	return c.CheckWeight(order) && c.CheckRegion(order)
}
func (c *Courier1) CheckWeight(order Order) bool {
	if c.TotalWeight+order.Weight > c.MaxWeight {
		return false
	}
	return true
}

func (c *Courier1) CheckRegion(order Order) bool {
	for _, region := range c.Regions {
		if region == order.Regions {
			return true
		}
	}
	return false
}

func (c *Courier1) TakeOrder(order Order) {
	c.Orders = append(c.Orders, order)
	c.TotalWeight += order.Weight
	c.TotalCost += float32(order.Cost)

	// Calculate delivery time and cost
	firstOrder := len(c.Orders) == 1
	switch c.CourierType {
	case Foot:
		if firstOrder {
			c.TotalTime += 25
		} else {
			c.TotalTime += 10
		}
		c.TotalCost += float32(order.Cost)
	case Bike, Car:
		if firstOrder {
			c.TotalTime += 12
		} else {
			c.TotalTime += 8
		}
		if !firstOrder {
			c.TotalCost += float32(order.Cost) * 0.8
		} else {
			c.TotalCost += float32(order.Cost)
		}
	}
}
