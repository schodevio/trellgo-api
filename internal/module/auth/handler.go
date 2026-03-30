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

// Refresh godoc
// @Summary      Refresh access token
// @Tags         auth
// @Produce      json
// @Success      200  {object}  SignInResponse
// @Failure      401  {object}  apierrors.ErrorResponse
// @Router       /auth/refresh [post]
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

// SignIn godoc
// @Summary      Sign in
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      SignInRequest  true  "Credentials"
// @Success      200   {object}  SignInResponse
// @Failure      400   {object}  apierrors.ErrorResponse
// @Failure      401   {object}  apierrors.ErrorResponse
// @Router       /auth/signin [post]
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

// SignUp godoc
// @Summary      Sign up
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      SignUpRequest  true  "Registration data"
// @Success      201   {object}  SignUpResponse
// @Failure      400   {object}  apierrors.ErrorResponse
// @Failure      409   {object}  apierrors.ErrorResponse
// @Failure      422   {object}  apierrors.ErrorResponse
// @Router       /auth/signup [post]
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

// SignOut godoc
// @Summary      Sign out
// @Tags         auth
// @Produce      json
// @Success      204
// @Failure      401  {object}  apierrors.ErrorResponse
// @Router       /auth/signout [delete]
func (h *handler) SignOut(ctx fiber.Ctx) error {
	token := ctx.Cookies(refreshTokenCookie)
	if token == "" {
		return apierrors.Unauthorized("missing refresh token")
	}

	if err := h.service.SignOut(token); err != nil {
		return err
	}

	ctx.ClearCookie(refreshTokenCookie)

	return ctx.SendStatus(fiber.StatusNoContent)
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
