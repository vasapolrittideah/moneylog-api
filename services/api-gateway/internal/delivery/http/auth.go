package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/vasapolrittideah/moneylog-api/services/api-gateway/internal/payload"
	"github.com/vasapolrittideah/moneylog-api/services/api-gateway/internal/validator"
	authclient "github.com/vasapolrittideah/moneylog-api/services/auth-service/pkg/client"
	"github.com/vasapolrittideah/moneylog-api/shared/contract"
	authpbv1 "github.com/vasapolrittideah/moneylog-api/shared/protos/auth/v1"
	"google.golang.org/grpc/status"
)

type AuthHTTPHandler struct {
	authServiceClient *authclient.AuthServiceClient
	router            fiber.Router
	logger            *zerolog.Logger
}

func NewAuthHTTPHandler(
	authServiceClient *authclient.AuthServiceClient,
	router fiber.Router,
	logger *zerolog.Logger,
) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authServiceClient: authServiceClient,
		router:            router,
		logger:            logger,
	}
}

func (h *AuthHTTPHandler) RegisterRoutes() {
	router := h.router.Group("/auth")
	router.Post("/login", h.login)
	router.Post("/signup", h.signUp)
}

func (h *AuthHTTPHandler) login(c *fiber.Ctx) error {
	var req payload.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			contract.NewErrorResponse(contract.ErrorCodeValidation, err.Error()),
		)
	}

	if errs := validator.ValidateStruct(req); len(errs) != 0 {
		return c.Status(http.StatusBadRequest).JSON(
			contract.NewValidationErrorResponse(errs),
		)
	}

	grpcResp, err := h.authServiceClient.Client.Login(c.Context(), &authpbv1.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		st := status.Convert(err)
		h.logger.Error().Err(st.Err()).Msg("Failed to login")

		errorCode := contract.ErrorCodeFromGRPCCode(st.Code())
		httpStatus := contract.HTTPStatusFromGRPCCode(st.Code())

		return c.Status(httpStatus).JSON(
			contract.NewErrorResponse(errorCode, "failed to login"),
		)
	}

	apiResp := contract.NewSuccessResponse(&payload.LoginResponse{
		AccessToken:  grpcResp.GetAccessToken(),
		RefreshToken: grpcResp.GetRefreshToken(),
	})

	return c.Status(http.StatusOK).JSON(apiResp)
}

func (h *AuthHTTPHandler) signUp(c *fiber.Ctx) error {
	var req payload.SignUpRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			contract.NewErrorResponse(contract.ErrorCodeValidation, err.Error()),
		)
	}

	if errs := validator.ValidateStruct(req); len(errs) != 0 {
		return c.Status(http.StatusBadRequest).JSON(
			contract.NewValidationErrorResponse(errs),
		)
	}

	grpcResp, err := h.authServiceClient.Client.SignUp(c.Context(), &authpbv1.SignUpRequest{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	})
	if err != nil {
		st := status.Convert(err)
		h.logger.Error().Err(st.Err()).Msg("Failed to sign up")

		errorCode := contract.ErrorCodeFromGRPCCode(st.Code())
		httpStatus := contract.HTTPStatusFromGRPCCode(st.Code())

		return c.Status(httpStatus).JSON(
			contract.NewErrorResponse(errorCode, "failed to sign up"),
		)
	}

	apiResp := contract.NewSuccessResponse(&payload.SignUpResponse{
		AccessToken:  grpcResp.GetAccessToken(),
		RefreshToken: grpcResp.GetRefreshToken(),
	})

	return c.Status(http.StatusOK).JSON(apiResp)
}
