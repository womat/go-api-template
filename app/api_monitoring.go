package app

import (
	"github.com/womat/go-api-template/app/service/monitoring"
	"github.com/womat/go-api-template/pkg/web"
	"net/http"
)

// HandleMonitoring
//
//	@Summary		Get monitoring data.
//	@Description	Get monitoring data for WATCHIT.
//	@Tags			info
//	@Success		200	{object}	[]monitoring.Model
//	@Failure		400	{object}	ApiError
//	@Failure		403	{object}	ApiError
//	@Failure		422	{object}	ApiError
//	@Router			/monitoring [get]
//	@Security		APIKeyAuth
func (app *App) HandleMonitoring() http.Handler {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			resp, err := monitoring.Monitoring(r.Host, VERSION)
			if err != nil {
				web.Encode(w, http.StatusInternalServerError, web.NewApiError(err))
				return
			}

			web.Encode(w, http.StatusOK, resp)
		},
	)
}
