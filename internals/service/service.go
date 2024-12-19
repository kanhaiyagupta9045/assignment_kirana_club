package service

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kanhaiyagupta9045/kirana_club/internals/models"
	"github.com/kanhaiyagupta9045/kirana_club/internals/process"
	"github.com/kanhaiyagupta9045/kirana_club/internals/repository"
	"github.com/kanhaiyagupta9045/kirana_club/message_broker"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SubmitJobHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		var storeVisit models.StoresVisit

		if err := c.BindJSON(&storeVisit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error:": err.Error()})
			return
		}
		if err := validateData(storeVisit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
			return
		}
		storeVisit.Status = "ongoing"

		srv := repository.NewStoreService()
		jobId, err := srv.InsertStoreVisitService(storeVisit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
			return
		}
		go process.ProcessJob(jobId, storeVisit)
		producer, err := message_broker.NewProducer(os.Getenv("RABBITMQ_URL"), os.Getenv("QUEUE_NAME"))
		if err != nil {
			log.Printf("%v", err)
		}
		data := message_broker.Data{
			JobId:       jobId,
			Store_Visit: storeVisit,
		}
		producer.Publish(data)
		c.JSON(http.StatusCreated, gin.H{"job_id: ": jobId.Hex()})
	}
}

func GetJobInfoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		jobid := c.Query("jobid")
		if jobid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": "please provide the job Id"})
			return
		}

		id, err := primitive.ObjectIDFromHex(jobid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error:": err.Error()})
			return
		}
		svs := repository.NewStoreService()

		status, errMssg, failedStoreID, err := svs.GetStatusAndErrorByID(id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error:": "job id doesn't exist"})
			return
		} else {
			if status == "completed" || status == "ongoing" {
				type Response struct {
					Status string `json:"status"`
					JobID  string `json:"job_id"`
				}

				resp := &Response{
					Status: status,
					JobID:  jobid,
				}
				c.JSON(http.StatusOK, resp)
			} else {
				type ErrStruct struct {
					StoreID string `json:"store_id"`
					Error   string `json:"error"`
				}

				type Response struct {
					Status string    `json:"status"`
					JobID  string    `json:"job_id"`
					Error  ErrStruct `json:"error"`
				}

				resp := &Response{
					Status: status,
					JobID:  jobid,
					Error: ErrStruct{
						StoreID: failedStoreID,
						Error:   errMssg,
					},
				}
				c.JSON(http.StatusOK, resp)
			}
		}
	}
}
func validateData(sv models.StoresVisit) error {
	if sv.Count < 0 {
		return fmt.Errorf("count should not be less than zero")
	}
	if sv.Count != len(sv.Visits) {
		return fmt.Errorf("count should be equal to len of visits")
	}
	if len(sv.Visits) == 0 {
		return fmt.Errorf("visits should not be emtpy")
	}

	for _, item := range sv.Visits {
		if item.StoreID == "" {
			return fmt.Errorf("store_id is required")
		}
		if item.VisitTime == "" {
			return fmt.Errorf("visit_time is required")
		}
		if len(item.ImageURLs) == 0 {
			return fmt.Errorf("image_url cannot be empty for store_id: %s", item.StoreID)
		}
	}
	return nil
}
