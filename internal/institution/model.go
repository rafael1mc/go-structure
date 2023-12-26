package institution

type InstitutionResponse struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name_"`
}
