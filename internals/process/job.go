package process

import (
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/kanhaiyagupta9045/kirana_club/internals/image"
	"github.com/kanhaiyagupta9045/kirana_club/internals/models"
	"github.com/kanhaiyagupta9045/kirana_club/internals/repository"
	"github.com/kanhaiyagupta9045/kirana_club/internals/store"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ProcessJob(id primitive.ObjectID, storevisit models.StoresVisit) {
	fmt.Printf("Starting new job for id %v\n", id.Hex())

	sm, err := store.NewStoreManager()
	if err != nil {
		log.Println(err)
		return
	}

	svs := repository.NewStoreService()

	for index, store := range storevisit.Visits {
		if !sm.CheckStoreIDExist(store.StoreID) {
			log.Printf("Failed job for id %v\n", id.Hex())
			svs.UpdateStoreVisitServiceStatus(id, "failed", "store ID does not exist", store.StoreID)
			return
		}

		new_image_uuids := make([]string, len(store.ImageURLs))
		new_image_perims := make([]int64, len(store.ImageURLs))

		copy(new_image_uuids, store.ImageUUIDs)
		copy(new_image_perims, store.Perimeters)

		initial_done := len(store.ImageUUIDs)

		for i, img_url := range store.ImageURLs {

			// to resume an ongoing but failed in between job
			// skips images already processed
			if i < initial_done {
				continue
			}

			img_holder, err := image.DownloadImage(img_url)

			if err != nil {
				svs.UpdateStoreVisitServiceStatus(id, "failed", err.Error(), store.StoreID)
				return
			}

			err = img_holder.SaveImage(id.Hex(), store.StoreID)

			if err != nil {
				svs.UpdateStoreVisitServiceStatus(id, "failed", err.Error(), store.StoreID)
				return
			}

			// calculate perimeter
			perim := int64(img_holder.Width) * int64(img_holder.Height)

			// gpu processing simulation
			ms := 100 + rand.IntN(301)
			time.Sleep(time.Duration(ms) * time.Millisecond)

			new_image_perims[i] = perim
			new_image_uuids[i] = fmt.Sprintf("%s.%s", img_holder.ID, img_holder.Format)
		}
		err = svs.UpdateVisitInfo(id, index, new_image_perims, new_image_uuids)
		if err != nil {
			log.Printf("Failed job for id %v\n", id.Hex())
			return
		}
	}
	svs.UpdateStoreVisitServiceStatus(id, "completed", "", "")
	log.Printf("Completed job for id %v", id.Hex())
}
