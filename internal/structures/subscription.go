package structures

type Subscription struct {
	ID          int    `json:"id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date, omitempty"`
}

type Counting struct {
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	UserID      string `json:"user_id"`
	ServiceName string `json:"service_name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
