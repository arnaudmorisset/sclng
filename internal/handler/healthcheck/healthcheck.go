package healthcheck

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"
)

func NewPongHandler() handlers.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		log := logger.Get(r.Context())

		res, err := json.Marshal(map[string]string{"status": "pong"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			e := fmt.Errorf("fail to encode JSON: %s", err.Error())
			log.WithError(e).Error(e.Error())
			return e
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(res)
		return nil
	}
}
