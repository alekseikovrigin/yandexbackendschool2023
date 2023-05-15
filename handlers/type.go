package handlers

import (
	"github.com/alekseikovrigin/yandexbackendschool2023/models"
	"github.com/alekseikovrigin/yandexbackendschool2023/repository"
	"log"
)

var typeRepository repository.TypeRepository

func init() {
	typeRepository = repository.NewTypeRepository()
}

func TypeInsert() error {
	res := typeRepository.FindAll()

	if len(res) == 0 {
		foot := models.Type{
			ID:          "FOOT",
			Ratio:       2,
			RatingRatio: 3,
			MaxWeight:   10,
			MaxOrders:   2,
			MaxRegions:  1,
		}
		bike := models.Type{
			ID:          "BIKE",
			Ratio:       3,
			RatingRatio: 2,
			MaxWeight:   20,
			MaxOrders:   4,
			MaxRegions:  2,
		}
		auto := models.Type{
			ID:          "AUTO",
			Ratio:       4,
			RatingRatio: 1,
			MaxWeight:   40,
			MaxOrders:   7,
			MaxRegions:  3,
		}
		id1, err := typeRepository.Save(foot)
		id2, err := typeRepository.Save(bike)
		id3, err := typeRepository.Save(auto)
		log.Println(id1, id2, id3)
		log.Println(err)
	}
	return nil
}
