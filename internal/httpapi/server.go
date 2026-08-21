package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"github.com/VanceMichael/go-base-airbridge/internal/middleware"
	"github.com/VanceMichael/go-base-airbridge/internal/service"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	app *service.Application
	log *slog.Logger
}

func New(app *service.Application, log *slog.Logger) http.Handler {
	s := &Server{app: app, log: log}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(log))
	r.Get("/healthz", s.health)
	r.Get("/readyz", s.ready)
	r.Post("/v1/auth/login", s.login)
	r.With(s.auth).Post("/v1/auth/logout", s.logout)
	r.With(s.auth).Post("/v1/shipments", s.createShipment)
	r.With(s.auth).Get("/v1/shipments", s.listShipments)
	r.With(s.auth).Post("/v1/shipments/{id}/book", s.bookShipment)
	r.With(s.auth).Post("/v1/shipments/{id}/transition", s.transitionShipment)
	r.With(s.auth).Post("/v1/flight-legs", s.createLeg)
	r.With(s.auth).Post("/v1/flight-legs/{id}/open", s.openLeg)
	r.With(s.auth).Post("/v1/customs/{shipmentID}", s.customs)
	r.With(s.auth).Post("/v1/security/{shipmentID}", s.security)
	r.With(s.auth).Get("/v1/operations/summary", s.summary)
	return r
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	write(w, 200, map[string]any{"ok": true})
}
func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	if err := s.app.Ping(r.Context()); err != nil {
		fail(w, 503, err)
		return
	}
	write(w, 200, map[string]any{"ready": true})
}
func (s *Server) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		u, _, err := s.app.Authenticate(r.Context(), token)
		if err != nil {
			fail(w, 401, domain.ErrForbidden)
			return
		}
		next.ServeHTTP(w, r.WithContext(middleware.WithUser(r.Context(), u)))
	})
}
func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if !decode(r, &in) {
		fail(w, 400, domain.ErrInvalid)
		return
	}
	out, err := s.app.Login(r.Context(), in.Email, in.Password)
	if err != nil {
		fail(w, status(err), err)
		return
	}
	write(w, 200, map[string]any{"token": out.Token, "expires_at": out.ExpiresAt, "user": out.User})
}
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if err := s.app.Logout(r.Context(), token); err != nil {
		fail(w, status(err), err)
		return
	}
	write(w, 204, nil)
}
func (s *Server) createShipment(w http.ResponseWriter, r *http.Request) {
	u := mustUser(r)
	var in domain.Shipment
	if !decode(r, &in) {
		fail(w, 400, domain.ErrInvalid)
		return
	}
	v, err := s.app.CreateShipment(r.Context(), u, in, r.Header.Get("Idempotency-Key"))
	if err != nil {
		fail(w, status(err), err)
		return
	}
	write(w, 201, v)
}
func (s *Server) listShipments(w http.ResponseWriter, r *http.Request) {
	u := mustUser(r)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	v, err := s.app.ListShipments(r.Context(), u, domain.PageRequest{Limit: limit, Cursor: r.URL.Query().Get("cursor")})
	if err != nil {
		fail(w, status(err), err)
		return
	}
	write(w, 200, v)
}
func (s *Server) bookShipment(w http.ResponseWriter, r *http.Request) {
	u := mustUser(r)
	var in struct {
		LegID string `json:"leg_id"`
	}
	if !decode(r, &in) {
		fail(w, 400, domain.ErrInvalid)
		return
	}
	v, err := s.app.BookShipment(r.Context(), u, chi.URLParam(r, "id"), in.LegID)
	if err != nil {
		fail(w, status(err), err)
		return
	}
	write(w, 200, v)
}
func (s *Server) transitionShipment(w http.ResponseWriter, r *http.Request) {
	u := mustUser(r)
	var in struct {
		Status domain.ShipmentStatus `json:"status"`
	}
	if !decode(r, &in) {
		fail(w, 400, domain.ErrInvalid)
		return
	}
	v, err := s.app.TransitionShipment(r.Context(), u, chi.URLParam(r, "id"), in.Status)
	if err != nil {
		fail(w, status(err), err)
		return
	}
	write(w, 200, v)
}
func (s *Server) createLeg(w http.ResponseWriter, r *http.Request) {
	u := mustUser(r)
	var in domain.FlightLeg
	if !decode(r, &in) {
		fail(w, 400, domain.ErrInvalid)
		return
	}
	v, err := s.app.CreateLeg(r.Context(), u, in)
	if err != nil {
		fail(w, status(err), err)
		return
	}
	write(w, 201, v)
}
func (s *Server) openLeg(w http.ResponseWriter, r *http.Request) {
	u := mustUser(r)
	v, err := s.app.OpenLeg(r.Context(), u, chi.URLParam(r, "id"))
	if err != nil {
		fail(w, status(err), err)
		return
	}
	write(w, 200, v)
}
func (s *Server) customs(w http.ResponseWriter, r *http.Request) {
	u := mustUser(r)
	var in domain.CustomsCase
	if !decode(r, &in) {
		fail(w, 400, domain.ErrInvalid)
		return
	}
	in.ShipmentID = chi.URLParam(r, "shipmentID")
	v, err := s.app.PutCustoms(r.Context(), u, in)
	if err != nil {
		fail(w, status(err), err)
		return
	}
	write(w, 200, v)
}
func (s *Server) security(w http.ResponseWriter, r *http.Request) {
	u := mustUser(r)
	var in domain.SecurityCheck
	if !decode(r, &in) {
		fail(w, 400, domain.ErrInvalid)
		return
	}
	in.ShipmentID = chi.URLParam(r, "shipmentID")
	v, err := s.app.PutSecurity(r.Context(), u, in)
	if err != nil {
		fail(w, status(err), err)
		return
	}
	write(w, 200, v)
}
func (s *Server) summary(w http.ResponseWriter, r *http.Request) {
	v, err := s.app.Summary(r.Context(), mustUser(r))
	if err != nil {
		fail(w, status(err), err)
		return
	}
	write(w, 200, v)
}
func mustUser(r *http.Request) domain.User {
	v, _ := middleware.UserFrom(r.Context())
	return v.(domain.User)
}
func decode(r *http.Request, v any) bool {
	defer r.Body.Close()
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	return d.Decode(v) == nil
}
func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}
func fail(w http.ResponseWriter, status int, err error) {
	code := "internal_error"
	switch {
	case errors.Is(err, domain.ErrInvalid):
		code = "invalid_request"
	case errors.Is(err, domain.ErrForbidden):
		code = "forbidden"
	case errors.Is(err, domain.ErrNotFound):
		code = "not_found"
	case errors.Is(err, domain.ErrConflict):
		code = "conflict"
	case errors.Is(err, domain.ErrCapacity):
		code = "capacity_unavailable"
	case errors.Is(err, domain.ErrState):
		code = "invalid_state"
	}
	write(w, status, map[string]string{"code": code, "message": fmt.Sprint(err)})
}
func status(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalid):
		return 400
	case errors.Is(err, domain.ErrForbidden), errors.Is(err, domain.ErrExpired), errors.Is(err, domain.ErrRevoked):
		return 401
	case errors.Is(err, domain.ErrNotFound):
		return 404
	case errors.Is(err, domain.ErrConflict), errors.Is(err, domain.ErrCapacity), errors.Is(err, domain.ErrState):
		return 409
	default:
		return 500
	}
}

var _ = time.UTC
