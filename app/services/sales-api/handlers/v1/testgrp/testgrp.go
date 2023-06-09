package testgrp

import (
	"context"
	"net/http"
	"github.com/mjanes71/ardan-labs-learning/foundation/web"
)

// Status represents a test handler for now.
func Status(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}
	return web.Respond(ctx, w, status, http.StatusOK)
}