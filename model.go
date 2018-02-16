package main

type params struct {
	CID    *int
	User   *string
	Passwd *string
	URL    *string
}

type config struct {
	UserID     string `json:"user_id,omitempty"`
	CustomerID int    `json:"customer_id,omitempty"`
	URL        string `json:"url,omitempty"`
	Atoken     string `json:"access_token,omitempty"`
	Rtoken     string `json:"refresh_token,omitempty"`
}

type login struct {
	CustomerID int    `json:"customer_id,omitempty"`
	IP         string `json:"ip_address,omitempty"`
}

type tokenReturn struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type returnStatus struct {
	Code    int    `json:"code,omitempty"`
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}

type paging struct {
	Cursor string `json:"next_cursor"`
	Qty    int    `json:"quantity"`
}

type returnMsg struct {
	Status returnStatus  `json:"status,omitempty"`
	Data   []tokenReturn `json:"data,omitempty"`
}
