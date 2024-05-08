package emailshdl

import "main/internal/core/domain"

type sendEmail struct {
	ID           string                 `json:"id"`
	From         *domain.EmailAddress   `json:"from" validate:"nil=false"`
	Destinations *[]domain.EmailAddress `json:"destinations" validate:"nil=false > empty=false"`
	PlainContent string                 `json:"plainContent" validate:"empty=false"`
	HtmlContent  string                 `json:"htmlContent"`
	Subject      string                 `json:"subject" validate:"empty=false"`
}
