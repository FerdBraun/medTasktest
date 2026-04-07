package ui

import (
	"embed"
	"net/http"
)

//go:embed index.html
var uiFS embed.FS

type Handler struct {
	content []byte
}

func NewHandler() *Handler {
	content, err := uiFS.ReadFile("index.html")
	if err != nil {
		panic(err)
	}

	return &Handler{content: content}
}

func (h *Handler) ServeUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.content)
}
