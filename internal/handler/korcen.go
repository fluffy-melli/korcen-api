// internal/handler/korcen.go

package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/fluffy-melli/korcen-api/pkg/check"
	"github.com/gin-gonic/gin"
)

// @Summary     Process Korcen Request
// @Description Processes a Korcen request and returns the result
// @Tags        korcen
// @Accept      json,xml
// @Produce     json,xml
// @Param       input  body  check.Header  true  "Korcen Input"
// @Success     200    {object}  check.Respond    "Korcen Result"
// @Failure     400    {object}  map[string]interface{}  "Invalid Request"
// @Failure     500    {object}  map[string]interface{}  "Internal Server Error"
// @Router      /api/v1/korcen [post]
func KorcenV1(c *gin.Context, system *actor.ActorSystem, korcenPID *actor.PID) {
	var header check.Header
	isXML := false

	switch c.ContentType() {
	case "text/xml", "application/xml":
		if err := c.ShouldBindXML(&header); err != nil {
			respond(c, http.StatusBadRequest, true, gin.H{"error": "Invalid XML request"})
			return
		}
		isXML = true
	default:
		if err := c.ShouldBindJSON(&header); err != nil {
			respond(c, http.StatusBadRequest, false, gin.H{"error": "Invalid JSON request"})
			return
		}
	}

	if strings.TrimSpace(header.Input) == "" {
		if isXML {
			c.XML(http.StatusBadRequest, gin.H{"error": "Invalid request: empty input"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: empty input"})
		}
		return
	}

	future := system.Root.RequestFuture(korcenPID, &check.KorcenRequest{Header: &header}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		if isXML {
			c.XML(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	korcenResp, ok := result.(*check.KorcenResponse)
	if !ok {
		if isXML {
			c.XML(http.StatusInternalServerError, gin.H{"error": "Invalid actor response"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid actor response"})
		}
		return
	}

	if korcenResp.Err != nil {
		if isXML {
			c.XML(http.StatusBadRequest, gin.H{"error": korcenResp.Err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": korcenResp.Err.Error()})
		}
		return
	}

	response := korcenResp.Respond
	respond(c, http.StatusOK, isXML, response)
}
