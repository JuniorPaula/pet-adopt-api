package dto

type CreatePetDto struct {
	Name   string   `json:"name"`
	Age    int      `json:"age"`
	Weight float64  `json:"weight"`
	Size   string   `json:"size"`
	Color  string   `json:"color"`
	Images []string `json:"images"`
}
