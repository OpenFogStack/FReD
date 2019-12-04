package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"gitlab.tu-berlin.de/mcc-fred/fred/fred-node/pkg/memorykg"
	"gitlab.tu-berlin.de/mcc-fred/fred/fred-node/pkg/memorysd"
)

// StorageDriver is an interface that abstracts the component that stores actual keygroup data.
type StorageDriver interface {
	Create(kg string, data string) (uint64, error)
	Read(kg string, id uint64) (string, error)
	Update(kg string, id uint64, data string) error
	Delete(kg string, id uint64) error
	CreateKeygroup(kg string) error
	DeleteKeygroup(kg string) error
}

// KeygroupManager is an interface that abstracts the component that stores keygroup metadata.
type KeygroupManager interface {
	Create(kg string) error
	Delete(kg string) error
	Exists(kg string) bool
}

func setupRouter(sd StorageDriver, kgm KeygroupManager) (r *gin.Engine) {
	r = gin.Default()

	r.POST("/keygroup/:kgname", func (context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		err := kgm.Create(kgname)

		if err != nil {
			context.Status(http.StatusConflict)
		} else {
			err = sd.CreateKeygroup(kgname)

			if err != nil {
				context.Status(http.StatusConflict)
			} else {
				context.Status(http.StatusOK)
			}
		}
	})

	r.DELETE("/keygroups/:kgname", func (context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		err := kgm.Delete(kgname)

		if err != nil {
			context.Status(http.StatusConflict)
		} else {
			err = sd.DeleteKeygroup(kgname)

			if err != nil {
				context.Status(http.StatusConflict)
			} else {
				context.Status(http.StatusOK)
			}
		}
	})

	r.POST("/keygroup/:kgname/items", func (context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		if !kgm.Exists(kgname) {
			context.Status(http.StatusNotFound)
		} else {
			var json struct {
				Data string `json:"data" binding:"required"`
			}

			if context.Bind(&json) != nil {
				context.Status(http.StatusBadRequest)
			} else {
				data := json.Data
				id, err := sd.Create(kgname, data)

				if err != nil {
					context.Status(http.StatusConflict)
				} else {
					context.String(http.StatusOK, string(id))
				}
			}
		}

	})

	r.GET("/keygroup/:kgname/items/:id", func (context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		if !kgm.Exists(kgname) {
			context.Status(http.StatusNotFound)
		} else {
			id, err := strconv.Atoi(context.Params.ByName("id"))

			if err == nil {
				data, err := sd.Read(kgname, uint64(id))

				if err != nil {
					context.Status(http.StatusConflict)
				} else {
					context.String(http.StatusOK, data)
				}
			} else {
				context.Status(http.StatusBadRequest)
			}
		}
	})

	r.PUT("/keygroup/:kgname/items/:id", func (context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		if !kgm.Exists(kgname) {
			context.Status(http.StatusNotFound)
		} else {
			id, err := strconv.Atoi(context.Params.ByName("id"))

			if err == nil {
				var json struct {
					Data string `json:"data" binding:"required"`
				}

				if context.Bind(&json) != nil {
					context.Status(http.StatusBadRequest)
				} else {
					data := json.Data
					err := sd.Update(kgname, uint64(id), data)

					if err != nil {
						context.Status(http.StatusConflict)
					} else {
						context.Status(http.StatusOK)
					}
				}
			} else {
				context.Status(http.StatusBadRequest)
			}
		}
	})

	r.DELETE("/keygroup/:kgname/items/:id", func (context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		if !kgm.Exists(kgname) {
			context.Status(http.StatusNotFound)
		} else {
			id, err := strconv.Atoi(context.Params.ByName("id"))

			if err == nil {

				err := sd.Delete(kgname, uint64(id))

				if err == nil {
					context.Status(http.StatusOK)
				} else {
					context.Status(http.StatusNotFound)
				}
			} else {
				context.Status(http.StatusBadRequest)
			}
		}
	})

	return
}

func main() {
	var sd StorageDriver = memorysd.New()
	var kgm KeygroupManager = memorykg.New()

	r := setupRouter(sd, kgm)

	r.Run(":9001")
}
