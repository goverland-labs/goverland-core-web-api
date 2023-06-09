package dao

type TopCategory struct {
	TotalCount uint64 `json:"total_count"`
	List       []Dao  `json:"list"`
}

type TopCategories map[string]TopCategory
