package dto

type CreatePetDto struct {
	Name        string   `json:"name"`
	Age         string   `json:"age"`
	Weight      string   `json:"weight"`
	Size        string   `json:"size"`
	Color       string   `json:"color"`
	Images      []string `json:"images"`
	Description string   `json:"description"`
}

type UpdatePetDto struct {
	Name        string `json:"name"`
	Age         string `json:"age"`
	Weight      string `json:"weight"`
	Size        string `json:"size"`
	Color       string `json:"color"`
	Available   *bool  `json:"available"`
	Description string `json:"description"`
}
