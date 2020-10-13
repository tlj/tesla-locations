package locations

type LocationType string

const (
	Delivery            LocationType = "delivery"
	DestinationCharger  LocationType = "destination charger"
	Logistics           LocationType = "logistics"
	SalesRepresentative LocationType = "sales representative"
	Service             LocationType = "service"
	StandardCharger     LocationType = "standard charger"
	Store               LocationType = "store"
	SuperCharger        LocationType = "supercharger"
)

type Location struct {
	AddressLine1           string `json:"address_line_1"`
	AddressLine2           string `json:"address_line_2"`
	AddressNotes           string `json:"address_notes"`
	Address                string `json:"address"`
	Amenities              string `json:"amenities"`
	BaiduLat               string `json:"baidu_lat"`
	BaiduLng               string `json:"baidu_lng"`
	Chargers               string `json:"locations"`
	City                   string `json:"city"`
	CommonName             string `json:"common_name"`
	Country                string `json:"country"`
	DestinationChargerLogo string `json:"destination_charger_logo"`
	DestinationWebsite     string `json:"destination_website"`
	DirectionsLink         string `json:"directions_link"`
	Emails                 []struct {
		Label string `json:"label"`
		Email string `json:"email"`
	} `json:"emails"`
	Geocode       string         `json:"geocode"`
	Hours         string         `json:"hours"`
	IsGallery     bool           `json:"is_gallery"`
	KioskPinX     string         `json:"kiosk_pin_x"`
	KioskPinY     string         `json:"kiosk_pin_y"`
	KioskZoomPinX string         `json:"kiosk_zoom_pin_x"`
	KioskZoomPinY string         `json:"kiosk_zoom_pin_y"`
	Latitude      string         `json:"latitude"`
	LocationID    string         `json:"location_id"`
	LocationType  []LocationType `json:"location_type"`
	Longitude     string         `json:"longitude"`
	Nid           string         `json:"nid"`
	OpenSoon      string         `json:"open_soon"`
	Path          string         `json:"path"`
	PostalCode    string         `json:"postal_code"`
	ProvinceState string         `json:"province_state"`
	Region        string         `json:"region"`
	SalesPhone    []struct {
		Label     string `json:"label"`
		Number    string `json:"number"`
		LineBelow string `json:"line_below"`
	} `json:"sales_phone"`
	SalesRepresentative bool   `json:"sales_representative"`
	SubRegion           string `json:"sub_region"`
	Title               string `json:"title"`
	TrtID               string `json:"trt_id"`
}

func (c *Location) HasCharger() bool {
	return c.Has([]LocationType{SuperCharger, StandardCharger, DestinationCharger})
}

func (c *Location) Has(types []LocationType) bool {
	for _, l := range c.LocationType {
		for _, t := range types {
			if t == l {
				return true
			}
		}
	}

	return false
}
