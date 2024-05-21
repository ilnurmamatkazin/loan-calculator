package http

import (
	"net/http"
)

func (s *Server) cache(w http.ResponseWriter, _ *http.Request) error {
	list := s.stor.Cache()

	if len(list) == 0 {
		return errorStatus{error: ErrEmptyCache, status: http.StatusBadRequest}
	}

	writeResponse(w, http.StatusOK, list)

	return nil
}
