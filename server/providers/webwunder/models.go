package webwunder

import "encoding/xml"

// Types defined in the WSDL

// SOAP envelope structures
type SoapRequestEnvelope struct {
	XMLName xml.Name        `xml:"soapenv:Envelope"`
	Soapenv string          `xml:"xmlns:soapenv,attr"`
	Gs      string          `xml:"xmlns:gs,attr"`
	Body    SoapRequestBody `xml:"soapenv:Body"`
}

type SoapRequestBody struct {
	XMLName xml.Name `xml:"soapenv:Body"`
	Content interface{}
}

type LegacyGetInternetOffers struct {
	XMLName xml.Name `xml:"gs:legacyGetInternetOffers"`
	Input   Input    `xml:"gs:input"`
}

// Input represents the request parameters for the API
type Input struct {
	XMLName      xml.Name `xml:"gs:input"`
	Installation bool     `xml:"gs:installation"`
	Connection   string   `xml:"gs:connectionEnum"`
	Address      Address  `xml:"gs:address"`
}

// Address represents a physical address
type Address struct {
	XMLName     xml.Name `xml:"gs:address"`
	Street      string   `xml:"gs:street"`
	HouseNumber string   `xml:"gs:houseNumber"`
	City        string   `xml:"gs:city"`
	PLZ         string   `xml:"gs:plz"`
	CountryCode string   `xml:"gs:countryCode"`
}

type SoapResponseEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Soapenv string   `xml:"xmlns:SOAP-ENV,attr"`
	Body    SoapResponseBody
}

type SoapResponseHeader struct {
	XMLName xml.Name `xml:"SOAP-ENV:Header"`
}

type SoapResponseBody struct {
	XMLName xml.Name `xml:"Body"`
	Output  Output
}

// Output represents the response containing products
type Output struct {
	XMLName  xml.Name  `xml:"Output"`
	Ns2      string    `xml:"xmlns:ns2,attr"`
	Products []Product `xml:"products"`
}

// Product represents an internet product offer
type Product struct {
	XMLName     xml.Name     `xml:"products"`
	Id          int          `xml:"productId"`
	Name        string       `xml:"providerName"`
	ProductInfo *ProductInfo `xml:"productInfo,omitempty"`
}

// ProductInfo contains details about an internet product
type ProductInfo struct {
	XMLName                        xml.Name `xml:"productInfo"`
	Speed                          int32    `xml:"speed"`
	MonthlyCostInCent              int32    `xml:"monthlyCostInCent"`
	MonthlyCostInCentFrom25thMonth int32    `xml:"monthlyCostInCentFrom25thMonth"`
	Voucher                        *Voucher `xml:"voucher,omitempty"`
	ContractDurationInMonths       int32    `xml:"contractDurationInMonths"`
	ConnectionType                 string   `xml:"connectionType"`
}

type Voucher struct {
	XMLName xml.Name `xml:"voucher"`

	PercentageVoucher `xml:",inline"`
	AbsoluteVoucher   `xml:",inline"`
}

// PercentageVoucher represents a percentage-based discount
type PercentageVoucher struct {
	XMLName           xml.Name `xml:"percentageVoucher"`
	Percentage        int32    `xml:"percentage"`
	MaxDiscountInCent int32    `xml:"maxDiscountInCent"`
}

// AbsoluteVoucher represents a fixed-amount discount
type AbsoluteVoucher struct {
	XMLName             xml.Name `xml:"absoluteVoucher"`
	DiscountInCent      int32    `xml:"discountInCent"`
	MinOrderValueInCent int32    `xml:"minOrderValueInCent"`
}
