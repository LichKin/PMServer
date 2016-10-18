package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"

	"github.com/gorilla/mux"
)

var Log *logrus.Logger = logrus.New()

func init() {
	Custom_MysqlManager.Connect()
}
func main() {

	dispatcher := mux.NewRouter()
	dispatcher.HandleFunc("/transfer", TransferHandler)
	dispatcher.HandleFunc("/pull/{uname}", PullHandler)
	http.Handle("/", dispatcher)

	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 1000
	http.ListenAndServe(":5000", nil)

}
