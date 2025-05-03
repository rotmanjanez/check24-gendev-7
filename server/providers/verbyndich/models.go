package verbyndich

import (
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
)

type Response struct {
	Product     string `json:"product"`
	Description string `json:"description"`
	Last        bool   `json:"last"`
	Valid       bool   `json:"valid"`
}

type UnitValue struct {
	Value int32
	Unit  string
}

// all int32 in DescriptionData are positive values and checked in the parser
// having int32 here helps to avoid separate checks when conterting to output format
type DescriptionData struct {
	Price                   *int32
	ConnectionType          m.ConnectionType
	Speed                   *UnitValue
	MinimalContractDuration *UnitValue
	PercentageDiscount      *m.PercentageDiscount
	AbsoluteDiscount        *m.AbsoluteDiscount
	SubsequentCost          *m.SubsequentCost
	UnthrottledCapacity     *UnitValue
	IncludedTVSender        string
	MaxAge                  *int32
	MinAge                  *int32
	MinOrderValue           *int32
	InstallationIncluded    bool
}
