package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	payment_service "github.com/polnaya-katuxa/ds-lab-02/gateway/internal/generated/openapi/clients/payment-service"
	"github.com/polnaya-katuxa/ds-lab-02/gateway/internal/models"
)

type PaymentServiceClient struct {
	c *payment_service.Client
}

func NewPaymentServiceClient(c *payment_service.Client) *PaymentServiceClient {
	return &PaymentServiceClient{
		c: c,
	}
}

func (c *PaymentServiceClient) Create(ctx context.Context, payment payment_service.CreatePaymentRequest) (*payment_service.PaymentInfo, error) {
	resp, err := c.c.Create(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		var validationError models.ValidationError
		err := json.Unmarshal(body, &validationError)
		if err != nil {
			return nil, fmt.Errorf("parse service error: %w", err)
		}

		return nil, validationError
	case http.StatusInternalServerError:
		var internalError models.InternalError
		err := json.Unmarshal(body, &internalError)
		if err != nil {
			return nil, fmt.Errorf("parse service error: %w", err)
		}

		internalError.StatusCode = resp.StatusCode

		return nil, internalError
	case http.StatusOK:
		var paymentInfo payment_service.PaymentInfo
		err := json.Unmarshal(body, &paymentInfo)
		if err != nil {
			return nil, fmt.Errorf("parse payment info: %w", err)
		}

		return &paymentInfo, nil
	default:
		return nil, fmt.Errorf("unknown response %d: %w", resp.StatusCode, models.ErrUnknownResponseStatus)
	}
}

func (c *PaymentServiceClient) Cancel(ctx context.Context, paymentUid uuid.UUID) error {
	resp, err := c.c.Cancel(ctx, paymentUid)
	if err != nil {
		return fmt.Errorf("cancel payment: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}
	resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound, http.StatusInternalServerError:
		var internalError models.InternalError
		err := json.Unmarshal(body, &internalError)
		if err != nil {
			return fmt.Errorf("parse service error: %w", err)
		}

		internalError.StatusCode = resp.StatusCode

		return internalError
	case http.StatusNoContent:
		return nil
	default:
		return fmt.Errorf("unknown response %d: %w", resp.StatusCode, models.ErrUnknownResponseStatus)
	}
}

func (c *PaymentServiceClient) Get(ctx context.Context, paymentUid uuid.UUID) (*payment_service.PaymentInfo, error) {
	resp, err := c.c.Get(ctx, paymentUid)
	if err != nil {
		return nil, fmt.Errorf("get payment: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound, http.StatusInternalServerError:
		var internalError models.InternalError
		err := json.Unmarshal(body, &internalError)
		if err != nil {
			return nil, fmt.Errorf("parse service error: %w", err)
		}

		internalError.StatusCode = resp.StatusCode

		return nil, internalError
	case http.StatusOK:
		var paymentInfo payment_service.PaymentInfo
		err := json.Unmarshal(body, &paymentInfo)
		if err != nil {
			return nil, fmt.Errorf("parse payment info: %w", err)
		}

		return &paymentInfo, nil
	default:
		return nil, fmt.Errorf("unknown response %d: %w", resp.StatusCode, models.ErrUnknownResponseStatus)
	}
}
