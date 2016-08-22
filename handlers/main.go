package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/clashr/rund/api"
)

func MainController(w http.ResponseWriter, r *http.Request) {
	maxsize := viper.GetInt64("maxsize")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		logrus.Info("received wrong method")
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxsize))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		logrus.Errorf("bad request: %s", err)
		return
	}
	defer r.Body.Close()

	var runRequest api.RunRequest

	if err := json.Unmarshal(body, &runRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(err)
		logrus.Errorf("bad request: %s", err)
		return
	}

	if len(runRequest.Command) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(errors.New("no command specified"))
		return
	}

	logrus.Info("received good request")
	result, err := executeRequest(runRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		logrus.Errorf("failed process: %s", err)
	}

	logrus.Info("writing response...")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(result)

	logrus.Info("finished processing request")
}
