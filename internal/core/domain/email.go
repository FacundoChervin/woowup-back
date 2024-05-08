package domain

type EmailAddress struct {
	Name  string `json:"name" validate:"empty=false"`
	Email string `json:"email" validate:"empty=false & format=email"`
}

type EmailBatch struct {
	ID         string         `json:"id"`
	Email      Email          `json:"email"`
	Recipìents []EmailAddress `json:"destinations"`
	Provider   int            `json:"provider" validate:"gte=0"`
}

type Email struct {
	ID           string       `json:"id"`
	From         EmailAddress `json:"from"`
	Subject      string       `json:"subject"`
	PlainContent string       `json:"plainContent"`
	HtmlContent  string       `json:"htmlContent"`
}

type EmailView struct {
	ID string `json:"id"`
}

type SendEmail struct {
	Email     Email        `json:"email"`
	Provider  int          `json:"provider"`
	Recipient EmailAddress `json:"destination"`
}

type EmailEntity struct {
	ID        string       `json:"batchId"`
	From      EmailAddress `json:"from"`
	Subject   string       `json:"subject"`
	Recipient string       `json:"recipient"`
	Sent      bool         `json:"status"`
	Retries   int          `json:"retries"`
	Content   string       `json:"content"`
}
