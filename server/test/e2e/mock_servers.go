package e2e

import (
	"net/http"
	"net/http/httptest"
)

// createByteMeTestServer creates a mock ByteMe provider server
func createByteMeTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/csv")
		csvData := `productId,providerName,speed,monthlyCostInCent,afterTwoYearsMonthlyCost,durationInMonths,connectionType,installationService,tv,limitFrom,maxAge,voucherType,voucherValue
1,ByteMe Basic,100,2999,1999,24,DSL,true,Premium,0,65,,0
2,ByteMe Premium,200,3999,2999,12,FIBER,true,Ultra,0,75,percentage,10
3,ByteMe Economy,50,1999,1499,24,DSL,false,Basic,50,65,absolute,500`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(csvData))
	}))
}

// createWebWunderTestServer creates a mock WebWunder SOAP provider server
func createWebWunderTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !hasSuffix(r.URL.Path, "/endpunkte/soap/ws") {
			http.NotFound(w, r)
			return
		}
		if r.Method != "POST" && r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/xml")

		// Return test SOAP response for WebWunder
		soapResponse := `<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
	<Body>
		<Output xmlns:ns2="http://webwunder.gendev7.check24.fun/offerservice">
			<products>
				<productId>1</productId>
				<providerName>WebWunder DSL Basic</providerName>
				<productInfo>
					<speed>100</speed>
					<monthlyCostInCent>2999</monthlyCostInCent>
					<monthlyCostInCentFrom25thMonth>3499</monthlyCostInCentFrom25thMonth>
					<contractDurationInMonths>24</contractDurationInMonths>
					<connectionType>DSL</connectionType>
				</productInfo>
			</products>
			<products>
				<productId>2</productId>
				<providerName>WebWunder Fiber Premium</providerName>
				<productInfo>
					<speed>250</speed>
					<monthlyCostInCent>4999</monthlyCostInCent>
					<monthlyCostInCentFrom25thMonth>5499</monthlyCostInCentFrom25thMonth>
					<contractDurationInMonths>12</contractDurationInMonths>
					<connectionType>FIBER</connectionType>
				</productInfo>
				<vouchers>
					<absoluteDiscount>
						<amount>1000</amount>
						<limit>2000</limit>
					</absoluteDiscount>
				</vouchers>
			</products>
		</Output>
	</Body>
</SOAP-ENV:Envelope>`

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(soapResponse))
	}))
}

// createVerbynDichTestServer creates a mock VerbynDich provider server
func createVerbynDichTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !hasSuffix(r.URL.Path, "/check24/data") {
			http.NotFound(w, r)
			return
		}
		if r.Method != "POST" && r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		jsonResponse := `{"product": "VerbynDich Premium", "description": "Für nur 29€ im Monat erhalten Sie eine DSL-Verbindung mit einer Geschwindigkeit von 100 Mbit/s. Zusätzlich sind folgende Fernsehsender enthalten TestTV.", "last": true, "valid": true}`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonResponse))
	}))
}

// createServusSpeedTestServer creates a mock ServusSpeed provider server
func createServusSpeedTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hasSuffix(r.URL.Path, "/api/external/available-products") && (r.Method == "POST" || r.Method == "GET") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"availableProducts": ["servus1", "servus2"]}`))
			return
		}
		if hasPrefix(r.URL.Path, "/api/external/product-details/") && (r.Method == "POST" || r.Method == "GET") {
			w.Header().Set("Content-Type", "application/json")
			id := r.URL.Path[len("/api/external/product-details/"):]
			var jsonResponse string
			switch id {
			case "servus1":
				jsonResponse = `{"servusSpeedProduct": {"providerName": "ServusSpeed Basic", "productInfo": {"speed": 100, "contractDurationInMonths": 24, "connectionType": "DSL"}, "pricingDetails": {"monthlyCostInCent": 2999, "installationService": false}, "discount": 0}}`
			case "servus2":
				jsonResponse = `{"servusSpeedProduct": {"providerName": "ServusSpeed Premium", "productInfo": {"speed": 500, "contractDurationInMonths": 12, "connectionType": "FIBER"}, "pricingDetails": {"monthlyCostInCent": 4999, "installationService": true}, "discount": 1000}}`
			default:
				http.NotFound(w, r)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(jsonResponse))
			return
		}
		http.NotFound(w, r)
	}))
}

// createPingPerfectTestServer creates a mock PingPerfect provider server
func createPingPerfectTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		jsonResponse := `[
			{"providerName": "PingPerfect Basic", "productInfo": {"speed": 100, "contractDurationInMonths": 24, "connectionType": "DSL"}, "pricingDetails": {"monthlyCostInCent": 2999, "installationService": "no"}},
			{"providerName": "PingPerfect Ultra", "productInfo": {"speed": 1000, "contractDurationInMonths": 12, "connectionType": "FIBER", "tv": "Premium Package"}, "pricingDetails": {"monthlyCostInCent": 6999, "installationService": "yes"}}
		]`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonResponse))
	}))
}

// Helpers for path matching
func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
