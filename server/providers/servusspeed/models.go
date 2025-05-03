package servusspeed

type AvailableProductsRequest struct {
	Address Address `json:"address"`
}

type Address struct {
	Street      string `json:"strasse"`
	HouseNumber string `json:"hausnummer"`
	PostalCode  string `json:"postleitzahl"`
	City        string `json:"stadt"`
	Country     string `json:"land"`
}

type AvailableProductsResponse struct {
	Products []string `json:"availableProducts"`
}

type ProductDetailsRequest struct {
	Address Address `json:"address"`
}

type ProductDetailsResponse struct {
	ServusSpeedProduct ServusSpeedProduct `json:"servusSpeedProduct"`
}

type ServusSpeedProduct struct {
	ProviderName   string         `json:"providerName"`
	ProductInfo    ProductInfo    `json:"productInfo"`
	PricingDetails PricingDetails `json:"pricingDetails"`
	Discount       int32          `json:"discount"`
}

type ProductInfo struct {
	Speed                    int32  `json:"speed"`
	ContractDurationInMonths int32  `json:"contractDurationInMonths"`
	ConnectionType           string `json:"connectionType"`
	TV                       string `json:"tv"`
	LimitFrom                int32  `json:"limitFrom"`
	MaxAge                   int32  `json:"maxAge"`
}

type PricingDetails struct {
	MonthlyCostInCent   int32 `json:"monthlyCostInCent"`
	InstallationService bool  `json:"installationService"`
}
