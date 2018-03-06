package main

type params struct {
	CID    *int
	User   *string
	Passwd *string
	URL    *string
}

type userParams struct {
	CID        *int
	User       *string
	Email      *string
	Role       *string
	Firstname  *string
	Lastname   *string
	Address    *string
	GroupID    *string
	CustomAttr *string
}

type pwLogin struct {
	CID    int
	User   string
	Passwd string
	IP     string
}

type config struct {
	UserID     string `json:"user_id,omitempty"`
	CustomerID int    `json:"customer_id,omitempty"`
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

type configFile struct {
	DirPath  string
	FilePath string
}

type userNew struct {
	CustomerID int                    `json:"customer_id"`
	UserID     string                 `json:"user_id"`
	Email      string                 `json:"email"`
	Role       string                 `json:"role"`
	Firstname  string                 `json:"firstname,omitempty"`
	Lastname   string                 `json:"lastname,omitempty"`
	Address    map[string]interface{} `json:"address,omitempty"`     // JSON object
	GroupID    []interface{}          `json:"group_id,omitempty"`    // JSON array
	CustomAttr map[string]interface{} `json:"custom_attr,omitempty"` // JSON object
}
