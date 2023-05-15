package handlers

import (
	"github.com/alekseikovrigin/yandexbackendschool2023/models"
	"github.com/alekseikovrigin/yandexbackendschool2023/repository"
	"github.com/alekseikovrigin/yandexbackendschool2023/util"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"strconv"
)

var courierRepository repository.CourierRepository

func init() {
	courierRepository = repository.NewCourierRepository()
}

// GetCourierRating Gets courier rating
func GetCourierRating(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 0)

	startDate := c.Query("start_date", "2023-01-01")
	endDate := c.Query("end_date", "2023-01-01")

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.Response{})
	}

	courier, err := courierRepository.FindByID(uint(id))
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(models.Response{})
	}

	if courier == nil {
		return c.Status(http.StatusNotFound).JSON(models.Response{})
	}

	orders := orderRepository.FindCompleteByCourierId(uint(id)) //TODO
	var sum int
	for _, val := range orders {
		sum = sum + (val.Cost * courier.Type.Ratio)
	}

	count := len(orders)
	hoursBetweenDates := util.DiffTimeHours(startDate, endDate)
	if hoursBetweenDates == 0 {
		hoursBetweenDates = 1
	}

	var rating = (count / hoursBetweenDates) * courier.Type.RatingRatio
	log.Println(count)

	resp := models.ResponseMetaInfo{
		Earnings: sum,
		Rating:   rating,
	}

	var regions []int
	for _, region := range courier.Regions {
		regions = append(regions, int(region))
	}

	resp.CourierId = courier.ID
	resp.CourierType = courier.Type.ID
	resp.Regions = regions
	resp.WorkingHours = courier.WorkingHours

	return c.Status(http.StatusOK).JSON(resp)
}

// GetAllCouriers gets all repository information
func GetAllCouriers(c *fiber.Ctx) error {
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

	couriers := courierRepository.FindAll(offset, limit)

	var responseCouriers []models.ResponseCourier
	for _, val := range couriers {

		var regions []int
		for _, region := range val.Regions {
			regions = append(regions, int(region))
		}

		tempCourier := models.ResponseCourier{
			CourierId:    val.ID,
			CourierType:  val.Type.ID,
			Regions:      regions,
			WorkingHours: val.WorkingHours,
		}
		responseCouriers = append(responseCouriers, tempCourier)
	}

	resp := models.ResponseCouriers{
		Couriers: responseCouriers,
		Limit:    limit,
		Offset:   offset,
	}

	return c.Status(http.StatusOK).JSON(resp)
}

// GetSingleCourier Gets single courier information
func GetSingleCourier(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 0)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.Response{})
	}

	courier, err := courierRepository.FindByID(uint(id))
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(models.Response{})
	}

	if courier == nil {
		return c.Status(http.StatusNotFound).JSON(models.Response{})
	}

	var regions []int
	for _, region := range courier.Regions {
		regions = append(regions, int(region))
	}

	resp := models.ResponseCourier{
		CourierId:    courier.ID,
		CourierType:  courier.Type.ID,
		Regions:      regions,
		WorkingHours: courier.WorkingHours,
	}

	return c.Status(http.StatusOK).JSON(resp)
}

// AddNewCouriers adds new order
func AddNewCouriers(c *fiber.Ctx) error {
	couriers := models.Couriers{}

	errParse := c.BodyParser(&couriers)
	log.Println(couriers.Couriers)
	addedCouriers, err := courierRepository.Add(couriers.Couriers)

	if errParse != nil {
		return c.Status(http.StatusBadRequest).JSON(models.Response{})
	}

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.Response{})
	}

	log.Println(addedCouriers)

	var responseCouriers []models.ResponseCourier
	for _, val := range couriers.Couriers {

		var regions []int
		for _, region := range val.Regions {
			regions = append(regions, int(region))
		}

		tempCourier := models.ResponseCourier{
			CourierId:    val.ID,
			CourierType:  *val.TypeID,
			Regions:      regions,
			WorkingHours: val.WorkingHours,
		}
		responseCouriers = append(responseCouriers, tempCourier)
	}

	resp := models.ResponseCouriers{
		Couriers: responseCouriers,
	}

	return c.Status(http.StatusOK).JSON(resp)
}
