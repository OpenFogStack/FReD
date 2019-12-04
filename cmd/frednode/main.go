package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorykg"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorysd"
)

var addr = flag.String("addr", ":9001", "http service address")

// Storage is an interface that abstracts the component that stores actual Keygroups data.
type Storage interface {
	Create(kg string, data string) (uint64, error)
	Read(kg string, id uint64) (string, error)
	Update(kg string, id uint64, data string) error
	Delete(kg string, id uint64) error
	CreateKeygroup(kg string) error
	DeleteKeygroup(kg string) error
}

// Keygroups is an interface that abstracts the component that stores Keygroups metadata.
type Keygroups interface {
	Create(kg string) error
	Delete(kg string) error
	Exists(kg string) bool
}

func setupRouter(sd Storage, kg Keygroups) (r *gin.Engine) {
	r = gin.Default()

	r.POST("/keygroup/:kgname", func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		err := kg.Create(kgname)

		if err != nil {
			context.Status(http.StatusConflict)
			return
		}

		err = sd.CreateKeygroup(kgname)

		if err != nil {
			context.Status(http.StatusConflict)
			return
		}

		context.Status(http.StatusOK)
		return
	})

	r.DELETE("/keygroup/:kgname", func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		err := kg.Delete(kgname)

		if err != nil {
			context.Status(http.StatusConflict)
			return
		}
		err = sd.DeleteKeygroup(kgname)

		if err != nil {
			context.Status(http.StatusConflict)
			return
		}

		context.Status(http.StatusOK)
		return
	})

	r.POST("/keygroup/:kgname/items", func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		if !kg.Exists(kgname) {
			context.Status(http.StatusNotFound)
			return
		}
		var json struct {
			Data string `json:"data" binding:"required"`
		}

		if err := context.ShouldBindJSON(&json) ; err != nil {
			context.Status(http.StatusBadRequest)
			return
		}

		data := json.Data
		id, err := sd.Create(kgname, data)

		if err != nil {
			context.Status(http.StatusConflict)
			return
		}

		context.String(http.StatusOK, string(id))
		return

	})

	r.GET("/keygroup/:kgname/items/:id", func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		if !kg.Exists(kgname) {
			context.Status(http.StatusNotFound)
			return
		}
		id, err := strconv.Atoi(context.Params.ByName("id"))

		if err != nil {
			context.Status(http.StatusBadRequest)
			return
		}

		data, err := sd.Read(kgname, uint64(id))

		if err != nil {
			context.Status(http.StatusNotFound)
			return
		}

		context.String(http.StatusOK, data)
		return
	})

	r.PUT("/keygroup/:kgname/items/:id", func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		if !kg.Exists(kgname) {
			context.Status(http.StatusNotFound)
			return
		}
		id, err := strconv.Atoi(context.Params.ByName("id"))

		if err != nil {
			context.Status(http.StatusBadRequest)
			return
		}

		var json struct {
			Data string `json:"data" binding:"required"`
		}

		if err := context.ShouldBindJSON(&json) ; err != nil {
			log.Print(err)
			context.Status(http.StatusBadRequest)
			return
		}

		data := json.Data
		err = sd.Update(kgname, uint64(id), data)

		if err != nil {
			context.Status(http.StatusConflict)
			return
		}

		context.Status(http.StatusOK)
		return
	})

	r.DELETE("/keygroup/:kgname/items/:id", func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		if !kg.Exists(kgname) {
			context.Status(http.StatusNotFound)
			return
		}

		id, err := strconv.Atoi(context.Params.ByName("id"))

		if err != nil {
			context.Status(http.StatusBadRequest)
			return
		}
		err = sd.Delete(kgname, uint64(id))

		if err != nil {
			context.Status(http.StatusNotFound)
			return
		}

		context.Status(http.StatusOK)
		return
	})

	return
}

func main() {
	var sd Storage = memorysd.New()
	var kg Keygroups = memorykg.New()

	r := setupRouter(sd, kg)

	log.Fatal(r.Run(*addr))
}
