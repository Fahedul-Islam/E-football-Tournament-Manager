package domain

type Group struct {
	GroupID     int            `json:"group_id"`
	GroupName   string         `json:"group_name"`
	Participants []*Participant `json:"participants"`
}
