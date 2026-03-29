package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/schodevio/trellgo/internal/platform/apierrors"
	"github.com/schodevio/trellgo/internal/platform/validator"
)

const refreshTokenCookie = "refresh_token"

type handler struct {
	service Service
}

func newHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) Refresh(ctx fiber.Ctx) error {
	token := ctx.Cookies(refreshTokenCookie)
	if token == "" {
		return apierrors.Unauthorized("missing refresh token")
	}

	req := RefreshRequest{Token: token}
	req.UserAgent = ctx.Get(fiber.HeaderUserAgent)
	req.IpAddress = ctx.IP()

	resp, err := h.service.RefreshToken(&req)
	if err != nil {
		return err
	}

	h.setRefreshTokenCookie(ctx, resp.RefreshToken)

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (h *handler) SignIn(ctx fiber.Ctx) error {
	var req SignInRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return apierrors.BadRequest("invalid request body")
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	req.UserAgent = ctx.Get(fiber.HeaderUserAgent)
	req.IpAddress = ctx.IP()

	resp, err := h.service.SignInUser(&req)
	if err != nil {
		return err
	}

	h.setRefreshTokenCookie(ctx, resp.RefreshToken)

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (h *handler) SignUp(ctx fiber.Ctx) error {
	var req SignUpRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return apierrors.BadRequest("invalid request body")
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	resp, err := h.service.SignUpUser(&req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

// private

func (h *handler) setRefreshTokenCookie(ctx fiber.Ctx, token string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     refreshTokenCookie,
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
		MaxAge:   int(refreshTokenTTL.Seconds()),
	})
}
