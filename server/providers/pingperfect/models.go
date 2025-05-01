package pingperfect

type Request struct {
	Street      string `json:"street"`
	PostalCode  string `json:"plz"`
	HouseNumber int32  `json:"houseNumber"`
	City        string `json:"city"`
	WantsFiber  bool   `json:"wantsFiber"`
}

type Response struct {
	Offers []InternetProduct
}

type InternetProduct struct {
	ProviderName   string         `json:"providerName"`
	ProductInfo    ProductInfo    `json:"productInfo"`
	PricingDetails PricingDetails `json:"pricingDetails"`
}

type ProductInfo struct {
	Speed                    int32  `json:"speed"`
	ContractDurationInMonths int32  `json:"contractDurationInMonths"`
	ConnectionType           string `json:"connectionType"`
	Tv                       string `json:"tv"`
	LimitFrom                int32  `json:"limitFrom"`
	MaxAge                   int32  `json:"maxAge"`
}

type PricingDetails struct {
	MonthlyCostInCent   int32  `json:"monthlyCostInCent"`
	InstallationService string `json:"installationService"`
}
