package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
)

const (
	Extract   = "/extract"
	Info      = "/info"
	Aggregate = "/aggregate"
)

type Server struct {
	extractor extractors.IExtractor
}

func New(extractor extractors.IExtractor) *Server {
	return &Server{
		extractor: extractor,
	}
}
func (controller *Server) Start(port string) error {
	http.HandleFunc(Info, controller.Info)
	http.HandleFunc(Extract, controller.Extract)
	http.HandleFunc(Aggregate, controller.Aggregate)

	fmt.Printf("Listening on port: %s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		return err
	}
	return nil
}
func (controller *Server) Extract(w http.ResponseWriter, req *http.Request) {
	b, err := controller.extract()
	if err != nil && err != extractors.NotFoundErr {
		http.Error(w, err.Error(), 400)
	} else if err == extractors.NotFoundErr {
		http.Error(w, err.Error(), 404)
	}

	fmt.Fprint(w, string(b))
}
func (controller *Server) extract() ([]byte, error) {
	ctx := context.Background()
	data, err := controller.extractor.Extract(ctx)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (controller *Server) Info(w http.ResponseWriter, req *http.Request) {
	b, err := json.Marshal(controller.extractor.Info())
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	fmt.Fprint(w, string(b))
}

func (controller *Server) Aggregate(w http.ResponseWriter, req *http.Request) {
	b, err := controller.extract()
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	fmt.Fprint(w, string(b))
}

func (controller *Server) aggregate(req *http.Request) ([]byte, error) {
	var values []extractors.Data

	if err := json.NewDecoder(req.Body).Decode(&values); err != nil {
		return nil, err
	}
	value, err := controller.extractor.Aggregate(values)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(&value)
	if err != nil {
		return nil, err
	}

	return b, nil
}
