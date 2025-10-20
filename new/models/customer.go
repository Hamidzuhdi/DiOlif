package models

type Customer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"` // TK, SD, SMP, Kelompok Tadarus
	Contact   string `json:"contact"`
	Address   string `json:"address"`
	CreatedAt string `json:"created_at"`
	Uniforms []CustomerUniform `json:"uniforms"`
}

type CustomerUniform struct {
    ID         int     `json:"id"`
    CustomerID int     `json:"customer_id"`
    UniformName string `json:"uniform_name"`
    Size       string  `json:"size"`
    Price      float64 `json:"price"`
    Notes      string  `json:"notes"`
    CreatedAt  string  `json:"created_at"`
}
