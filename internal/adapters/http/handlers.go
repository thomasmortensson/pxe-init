package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/thomasmortensson/pxe-init/internal/usecases"
)

// bootHandler handles the /boot.ipxe endpoint returning a static chainloader
func bootHandler(c *gin.Context) {
	c.String(http.StatusOK, "#!ipxe\nchain ipxe?uuid=${uuid}&mac=${mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}\n")
}

// statusHandler handles the /ready and /health endpoints. Can be extended to provide application DB health
func statusHandler(c *gin.Context) {
	// TODO - this obviously needs implementing with DB checking logic. I've run out of time so leaving as purely a liveness check
	c.String(http.StatusOK, "")
}

// ipxeHandler handles incoming requests of the form:
// /ipxe?uuid=${uuid}&mac=${mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}
// If an entry exists for the incoming MAC address, an ipxe script is generated to chainload with
// the specified kernel arguments for first run
func ipxeHandler(c *gin.Context) {
	// /ipxe?uuid=2f1214fe-59ba-9f42-a5c5-1af6f124aaf7&mac=08-00-27-c3-62-83&domain=&hostname=&serial=0

	serverCtx, ok := c.MustGet(paramServerContext).(*ServerContext)
	if !ok {
		fmt.Println("Failed to get server context")
		// Return 500
		c.String(http.StatusInternalServerError, "")
		return
	}

	// Get variables from context
	logger := serverCtx.Logger
	db := serverCtx.Datastore
	forwardAddress := serverCtx.ForwardServer.Host
	mac := c.Query("mac")

	usecase := usecases.NewPxeChainload(logger, db)

	pxeScript, err := usecase.Execute(mac, forwardAddress)
	if db.IsNotFound(err) {
		logger.Info(
			"Not Found",
			zap.String("mac-address", mac),
		)
		c.String(http.StatusNotFound, "")
		return
	} else if err != nil {
		logger.Error(
			"Failed in PxeChainload usecase",
			zap.Error(err),
		)
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.String(http.StatusOK, pxeScript)
}
