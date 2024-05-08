package emailssqs

import "main/internal/core/domain"

type sendEmailSqs struct {
	Email       *domain.Email        `json:"email" validate:"nil=false"`
	Provider    *int                 `json:"provider" validate:"nil=false"`
	Destination *domain.EmailAddress `json:"destination" validate:"nil=false"`
}
