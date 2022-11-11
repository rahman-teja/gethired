package model

type ToDo struct {
	ActivityGroupID *int64  `json:"activity_group_id"`
	Title           *string `json:"title"`
	IsActive        *bool   `json:"is_active"`
	Priority        *string `json:"priority"`
}

func (t ToDo) Validate() error {
	return nil
}
