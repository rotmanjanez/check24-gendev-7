package byteme

type Request struct {
	Street      string `url:"street"`
	HouseNumber string `url:"houseNumber"`
	City        string `url:"city"`
	PostalCode  string `url:"plz"`
}

type ResponseRow struct {
	Id                             string `csv:"productId"`
	ProviderName                   string `csv:"providerName"`
	Speed                          int32  `csv:"speed"`
	MonthlyCostInCent              int32  `csv:"monthlyCostInCent"`
	AfterTwoYearsMonthlyCostInCent int32  `csv:"afterTwoYearsMonthlyCost"`
	DurationInMonths               int32  `csv:"durationInMonths"`
	ConnectionType                 string `csv:"connectionType"`
	InstallationService            string `csv:"installationService"`
	TV                             string `csv:"tv"`
	LimitFrom                      int32  `csv:"limitFrom"`
	MaxAge                         int32  `csv:"maxAge"`
	VoucherType                    string `csv:"voucherType"`
	VoucherValue                   int32  `csv:"voucherValue"`
}

type Response struct {
	Offers []ResponseRow
}
