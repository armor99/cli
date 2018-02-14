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
