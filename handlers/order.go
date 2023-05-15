package handlers

import (
	"fmt"
	"github.com/alekseikovrigin/yandexbackendschool2023/models"
	"github.com/alekseikovrigin/yandexbackendschool2023/repository"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var orderRepository repository.OrderRepository

func init() {
	orderRepository = repository.NewOrderRepository()
}

// GetAllOrders gets all repository information
func GetAllOrders(c *fiber.Ctx) error {
	queryOffset := c.Query("offset", "0")
	queryLimit := c.Query("limit", "1")

	limit, err := strconv.Atoi(queryLimit)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.Response{})
	}

	offset, err := strconv.Atoi(queryOffset)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.Response{})
	}

	orders := orderRepository.FindAll(offset, limit)

	var responseOrders []models.ResponseOrder
	for _, val := range orders {
		tempOrder := models.ResponseOrder{
			Cost:          val.Cost,
			DeliveryHours: val.DeliveryHours,
			OrderID:       val.ID,
			Regions:       val.Regions,
			Weight:        val.Weight,
			CompletedTime: val.CompletedTime,
		}
		responseOrders = append(responseOrders, tempOrder)
	}
	return c.Status(http.StatusOK).JSON(responseOrders)
}

// GetSingleOrder Gets single order information
func GetSingleOrder(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 0)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.Response{})
	}

	order, err := orderRepository.FindByID(uint(id))
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(models.Response{})
	}

	if order == nil {
		return c.Status(http.StatusNotFound).JSON(models.Response{})
	}

	resp := models.ResponseOrder{
		Cost:          order.Cost,
		DeliveryHours: order.DeliveryHours,
		OrderID:       order.ID,
		Regions:       order.Regions,
		Weight:        order.Weight,
		CompletedTime: order.CompletedTime,
	}

	return c.Status(http.StatusOK).JSON(resp)
}

// MarkAsComplete adds new order
func MarkAsComplete(c *fiber.Ctx) error {
	orders := models.CompleteOrders{}

	errParse := c.BodyParser(&orders)

	if errParse != nil {
		return c.Status(http.StatusBadRequest).JSON(models.Response{})
	}

	values := make(map[uint][]uint)
	times := make(map[uint]string)

	for _, val := range orders.Orders {
		values[val.ID] = []uint{val.ID, val.CourierID}
		times[val.ID] = (val.CompletedTime).Format(time.RFC3339)
	}
	orders1 := orderRepository.FindAllBy(values)

	if len(orders1) == len(orders.Orders) {
		for _, val := range orders1 {
			val.CompletedTime, _ = time.Parse(time.RFC3339, times[val.ID])
			err := orderRepository.Update(*val)

			log.Println(err)
		}
	} else {
		return c.Status(http.StatusBadRequest).JSON(models.Response{})
	}

	var responseOrders []models.ResponseOrder
	for _, val := range orders1 {
		tempOrder := models.ResponseOrder{
			Cost:          val.Cost,
			DeliveryHours: val.DeliveryHours,
			OrderID:       val.ID,
			Regions:       val.Regions,
			Weight:        val.Weight,
			CompletedTime: val.CompletedTime,
		}
		responseOrders = append(responseOrders, tempOrder)
	}

	resp := responseOrders
	return c.Status(http.StatusOK).JSON(resp)
}

// AddNewOrders adds new order
func AddNewOrders(c *fiber.Ctx) error {
	orders := models.Orders{}

	errParse := c.BodyParser(&orders)
	addedOrders, err := orderRepository.Add(orders.Orders)

	if errParse != nil {
		return c.Status(http.StatusBadRequest).JSON(models.Response{})
	}

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.Response{})
	}

	var responseOrders []models.ResponseOrder
	for _, val := range addedOrders {
		tempOrder := models.ResponseOrder{
			Cost:          val.Cost,
			DeliveryHours: val.DeliveryHours,
			OrderID:       val.ID,
			Regions:       val.Regions,
			Weight:        val.Weight,
			//CompletedTime: "",
		}
		log.Println(val)
		responseOrders = append(responseOrders, tempOrder)
	}

	resp := responseOrders
	return c.Status(http.StatusOK).JSON(resp)
}

//AssignOrders gets all
func AssignOrders(c *fiber.Ctx) error {
	queryOffset := c.Query("offset", "0")
	queryLimit := c.Query("limit", "10000")

	limit, err := strconv.Atoi(queryLimit)
	if err != nil {
		return c.Status(500).JSON("Invalid limit option")
	}

	offset, err := strconv.Atoi(queryOffset)
	if err != nil {
		return c.Status(500).JSON("Invalid offset option")
	}

	ordersAll := orderRepository.FindAllAssign(offset, limit)

	couriers := []models.Courier1{
		{CourierType: models.Foot, MaxWeight: 10, MaxOrders: 2, Regions: []int{1, 2}},
		{CourierType: models.Bike, MaxWeight: 20, MaxOrders: 4, Regions: []int{2, 3, 4}},
		{CourierType: models.Car, MaxWeight: 40, MaxOrders: 7, Regions: []int{1, 2, 3, 4}},
	}

	// Sort the orders by descending order of cost
	sort.Sort(models.OrdersByCost(ordersAll))

	// Assign orders to couriers
	for _, order := range ordersAll {
		assigned := false
		for i := range couriers {
			courier := &couriers[i]
			if courier.CanTakeOrder(order) {
				courier.TakeOrder(order)
				assigned = true
				break
			}
		}
		if !assigned {
			fmt.Printf("Unable to assign order with weight %.2f to any courier\n", order.Weight)
		}
	}

	// Print the details of each courier
	for i := range couriers {
		courier := &couriers[i]
		fmt.Printf("%s: %d orders, %.2f kg, $%.2f, %.1f minutes\n",
			courier.CourierType, len(courier.Orders), courier.TotalWeight, courier.TotalCost, courier.TotalTime)
		for j := range courier.Orders {
			order := &courier.Orders[j]
			fmt.Printf("  - Order with weight %.2f, region %d\n", order.Weight, order.Regions)
		}
	}

	var AssignCouriers []models.AssignCourier
	var AssignCourier []models.AssignOrders

	var responseOrders []models.ResponseOrder
	for _, val := range ordersAll {
		tempOrder := models.ResponseOrder{
			Cost:          val.Cost,
			DeliveryHours: val.DeliveryHours,
			OrderID:       val.ID,
			Regions:       val.Regions,
			Weight:        val.Weight,
			CompletedTime: val.CompletedTime,
		}
		responseOrders = append(responseOrders, tempOrder)
	}

	resp2 := models.AssignOrders{
		Orders: responseOrders,
	}

	resp1 := models.AssignCourier{
		Orders: AssignCourier,
	}

	resp3 := models.AssignCouriers{
		Date:     time.Now(),
		Couriers: AssignCouriers,
	}

	log.Println(resp1)
	log.Println(resp2)
	log.Println(resp3)
	resp := responseOrders

	return c.Status(http.StatusOK).JSON(resp)
}
