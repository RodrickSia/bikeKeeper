package card

type Card struct {
	CardUID  string  `json:"cardUid"`
	CardType string  `json:"cardType"`
	MemberID *string `json:"memberId"`
	IsInside bool    `json:"isInside"`
	Status   string  `json:"status"`
}
