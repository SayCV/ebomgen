package webecd

import (
	//"fmt"
	//"regexp"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fedesog/webdriver"
	"github.com/saycv/ebomgen/pkg/types"
	"github.com/saycv/ebomgen/pkg/utils"

	"github.com/PuerkitoBio/goquery"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type NgDgkData struct {
	Props struct {
		Browser struct {
			Name    string `json:"name"`
			Version string `json:"version"`
			Major   string `json:"major"`
		} `json:"browser"`
		CorrelationID string `json:"correlationID"`
		Currency      string `json:"currency"`
		Device        struct {
		} `json:"device"`
		InterfaceView string `json:"interfaceView"`
		Lng           string `json:"lng"`
		Session       struct {
			CustomerID    int    `json:"CustomerId"`
			CustomerClass int    `json:"CustomerClass"`
			Currency      string `json:"Currency"`
			OrderModel    int    `json:"OrderModel"`
		} `json:"session"`
		ShowInterfaceViewToggle bool `json:"showInterfaceViewToggle"`
		PageProps               struct {
			Envelope struct {
				Data struct {
					PriceQuantity struct {
						QtyAvailable string `json:"qtyAvailable"`
						HasLeadTime  bool   `json:"hasLeadTime"`
						Pricing      []struct {
							DigikeyProductNumber string `json:"digikeyProductNumber"`
							ID                   string `json:"id"`
							Packaging            string `json:"packaging"`
							PricingTiers         []struct {
								BreakQty      string `json:"breakQty"`
								UnitPrice     string `json:"unitPrice"`
								ExtendedPrice string `json:"extendedPrice"`
							} `json:"pricingTiers"`
						} `json:"pricing"`
					} `json:"priceQuantity"`
					ProductOverview struct {
						DigikeyProductNumbers struct {
							Type  string `json:"type"`
							Value []struct {
								Label string `json:"label"`
								Value string `json:"value"`
							} `json:"value"`
						} `json:"digikeyProductNumbers"`
						Suppliers []struct {
							ID   int    `json:"id"`
							Name string `json:"name"`
							URL  string `json:"url"`
						} `json:"suppliers"`
						DatasheetURL              string `json:"datasheetUrl"`
						Description               string `json:"description"`
						DetailedDescription       string `json:"detailedDescription"`
						IsNormallyStocking        bool   `json:"isNormallyStocking"`
						Manufacturer              string `json:"manufacturer"`
						ManufacturerURL           string `json:"manufacturerUrl"`
						ManufacturerProductNumber string `json:"manufacturerProductNumber"`
						RolledUpProductID         string `json:"rolledUpProductId"`
						RolledUpProductNumber     string `json:"rolledUpProductNumber"`
						StandardLeadTime          string `json:"standardLeadTime"`
					} `json:"productOverview"`
					AdditionalResources struct {
						Title       string   `json:"title"`
						DataHeaders []string `json:"dataHeaders"`
						DataRows    []struct {
							DataCells []struct {
								Data struct {
									Value struct {
										Value string `json:"value"`
									} `json:"value"`
									Type string `json:"type"`
								} `json:"data"`
							} `json:"dataCells"`
						} `json:"dataRows"`
					} `json:"additionalResources"`
					Breadcrumb []struct {
						Label string `json:"label"`
						URL   string `json:"url"`
					} `json:"breadcrumb"`
					Environmental struct {
						Title       string   `json:"title"`
						DataHeaders []string `json:"dataHeaders"`
						DataRows    []struct {
							DataCells []struct {
								Data struct {
									Value struct {
										Value string `json:"value"`
									} `json:"value"`
									Type string `json:"type"`
								} `json:"data"`
							} `json:"dataCells"`
						} `json:"dataRows"`
					} `json:"environmental"`
					CarouselMedia []struct {
						AnalyticsTag string `json:"analyticsTag"`
						Title        string `json:"title"`
						HrefURL      string `json:"hrefUrl"`
						DisplayURL   string `json:"displayUrl"`
						Type         string `json:"type"`
					} `json:"carouselMedia"`
					OtherDocsAndMedia struct {
						Title       string   `json:"title"`
						DataHeaders []string `json:"dataHeaders"`
						DataRows    []struct {
							DataCells []struct {
								Data struct {
									Value struct {
										Value string `json:"value"`
									} `json:"value"`
									Type string `json:"type"`
								} `json:"data"`
							} `json:"dataCells"`
						} `json:"dataRows"`
					} `json:"otherDocsAndMedia"`
					IsBackOrderAllowed bool `json:"isBackOrderAllowed"`
					MinimumMultiplier  int  `json:"minimumMultiplier"`
					QuantityTable      []struct {
						BreakQty             int     `json:"breakQty"`
						DigikeyProductNumber string  `json:"digikeyProductNumber"`
						Fee                  int     `json:"fee"`
						ID                   string  `json:"id"`
						IsDiscount           bool    `json:"isDiscount"`
						MinMultiplier        int     `json:"minMultiplier"`
						Packaging            string  `json:"packaging"`
						PackTypeCode         int     `json:"packTypeCode"`
						UnitPrice            float64 `json:"unitPrice"`
					} `json:"quantityTable"`
					ProductAttributes struct {
						Attributes []struct {
							ID     string `json:"id"`
							Label  string `json:"label"`
							Values []struct {
								ID    string `json:"id"`
								Value string `json:"value"`
							} `json:"values"`
							IsFilterable bool   `json:"isFilterable"`
							Pt           string `json:"pt,omitempty"`
						} `json:"attributes"`
						Categories []struct {
							ID    string `json:"id"`
							Label string `json:"label"`
							URL   string `json:"url"`
						} `json:"categories"`
					} `json:"productAttributes"`
					Substitutes struct {
						Title       string   `json:"title"`
						DataHeaders []string `json:"dataHeaders"`
						DataRows    []struct {
							DataCells []struct {
								Data struct {
									Value struct {
										Label    string `json:"label"`
										URL      string `json:"url"`
										External bool   `json:"external"`
									} `json:"value"`
									Type string `json:"type"`
								} `json:"data"`
							} `json:"dataCells"`
						} `json:"dataRows"`
					} `json:"substitutes"`
					Associations struct {
						CardAssociations []struct {
							Title        string `json:"title"`
							CardCount    int    `json:"cardCount"`
							ProductCards []struct {
								Description string `json:"description"`
								DetailURL   string `json:"detailUrl"`
								ID          string `json:"id"`
								ImageURL    string `json:"imageUrl"`
								Mfr         string `json:"mfr"`
								MfrID       string `json:"mfrId"`
								MfrProduct  string `json:"mfrProduct"`
								UnitPrice   string `json:"unitPrice"`
							} `json:"productCards"`
							Type string `json:"type"`
						} `json:"cardAssociations"`
					} `json:"associations"`
					Messages []struct {
						Message string `json:"message"`
						Type    string `json:"type"`
						Code    string `json:"code"`
					} `json:"messages"`
					IsMarketplaceOnly bool `json:"isMarketplaceOnly"`
					IsMultiVendor     bool `json:"isMultiVendor"`
				} `json:"data"`
				Type string `json:"type"`
			} `json:"envelope"`
		} `json:"pageProps"`
		I18NCurrentInstance interface{} `json:"i18nCurrentInstance"`
		InitialI18NStore    struct {
			EnUs struct {
				DetailPage struct {
					CmsBomErrorDialog              string `json:"cms-bom-error-dialog"`
					CmsPreviousPage                string `json:"cms-previous-page"`
					CmsDetailedDescription         string `json:"cms-detailed-description"`
					CmsSeeAll                      string `json:"cms-see-all"`
					CmsNewBom                      string `json:"cms-new-bom"`
					CmsAddToCart                   string `json:"cms-add-to-cart"`
					CmsViewBom                     string `json:"cms-view-bom"`
					CmsTotal                       string `json:"cms-total"`
					CmsQuantity                    string `json:"cms-quantity"`
					CmsShowVat                     string `json:"cms-show-vat"`
					CmsBomNameLabel                string `json:"cms-bom-name-label"`
					CmsBomErrorProduct             string `json:"cms-bom-error-product"`
					CmsShowLess                    string `json:"cms-show-less"`
					CmsProductAddedToBom           string `json:"cms-product-added-to-bom"`
					CmsLeadtimeDisclaimer          string `json:"cms-leadtime-disclaimer"`
					CmsSupplier                    string `json:"cms-supplier"`
					CmsYourPrice                   string `json:"cms-your-price"`
					CmsMediaDownloads              string `json:"cms-media-downloads"`
					CmsFavoritesSuccess            string `json:"cms-favorites-success"`
					CmsLink                        string `json:"cms-link"`
					CmsLeadtimeInvalid             string `json:"cms-leadtime-invalid"`
					CmsStandardPrice               string `json:"cms-standard-price"`
					CmsFee                         string `json:"cms-fee"`
					CmsAttribute                   string `json:"cms-attribute"`
					CmsQtyExceededNoBackorderError string `json:"cms-qty-exceeded-no-backorder-error"`
					CmsMinQtyMultipleError         string `json:"cms-min-qty-multiple-error"`
					CmsLeadtimeTitle               string `json:"cms-leadtime-title"`
					CmsCopy                        string `json:"cms-copy"`
					CmsExistingBom                 string `json:"cms-existing-bom"`
					CmsManufacturerProductNumber   string `json:"cms-manufacturer-product-number"`
					CmsDescription                 string `json:"cms-description"`
					CmsNewParametricSearch         string `json:"cms-new-parametric-search"`
					CmsFrequentlyBoughtTogether    string `json:"cms-frequently-bought-together"`
					CmsFavorite                    string `json:"cms-favorite"`
					CmsNextPage                    string `json:"cms-next-page"`
					CmsMyPrice                     string `json:"cms-my-price"`
					CmsAddToFavorites              string `json:"cms-add-to-favorites"`
					CmsAddToBom                    string `json:"cms-add-to-bom"`
					CmsMoreCount                   string `json:"cms-more-count"`
					CmsViewFavorites               string `json:"cms-view-favorites"`
					CmsStandardLeadTime            string `json:"cms-standard-lead-time"`
					CmsQty                         string `json:"cms-qty"`
					CmsExtPrice                    string `json:"cms-ext-price"`
					CmsManufacturer                string `json:"cms-manufacturer"`
					CmsOk                          string `json:"cms-ok"`
					CmsType                        string `json:"cms-type"`
					CmsShowAll                     string `json:"cms-show-all"`
					CmsDatasheet                   string `json:"cms-datasheet"`
					CmsPriceBreak                  string `json:"cms-price-break"`
					CmsPackaging                   string `json:"cms-packaging"`
					CmsExtendedPrice               string `json:"cms-extended-price"`
					CmsDigikeyMarketplace          string `json:"cms-digikey-marketplace"`
					CmsLeadtimeUpdate              string `json:"cms-leadtime-update"`
					CmsSelect                      string `json:"cms-select"`
					CmsBomNameError                string `json:"cms-bom-name-error"`
					CmsFavoritesError              string `json:"cms-favorites-error"`
					CmsMatingProducts              string `json:"cms-mating-products"`
					CmsCategory                    string `json:"cms-category"`
					CmsFavorites                   string `json:"cms-favorites"`
					CmsCustomerReference           string `json:"cms-customer-reference"`
					CmsInvalidQtyError             string `json:"cms-invalid-qty-error"`
					CmsLeadtimeShipEstimate        string `json:"cms-leadtime-ship-estimate"`
					CmsBomName                     string `json:"cms-bom-name"`
					CmsPriceProcurement            string `json:"cms-price-procurement"`
					CmsReportError                 string `json:"cms-report-error"`
					CmsCopied                      string `json:"cms-copied"`
					CmsProductAttributes           string `json:"cms-product-attributes"`
					CmsDigikeyProductNumber        string `json:"cms-digikey-product-number"`
					CmsInStock                     string `json:"cms-in-stock"`
				} `json:"detail-page"`
				Common struct {
					CmsTopResults                string `json:"cms-top-results"`
					CmsSwitchToModernBody        string `json:"cms-switch-to-modern-body"`
					CmsJumpTo                    string `json:"cms-jump-to"`
					CmsFirstPage                 string `json:"cms-first-page"`
					CmsContinueShopping          string `json:"cms-continue-shopping"`
					CmsCancel                    string `json:"cms-cancel"`
					CmsCompareProductsCount      string `json:"cms-compare-products-count"`
					CmsFilter                    string `json:"cms-filter"`
					CmsPreviousPage              string `json:"cms-previous-page"`
					CmsLastPage                  string `json:"cms-last-page"`
					CmsClassicSearch             string `json:"cms-classic-search"`
					CmsError                     string `json:"cms-error"`
					CmsAppliedFilters            string `json:"cms-applied-filters"`
					CmsSwitchToModernLinkEnd     string `json:"cms-switch-to-modern-link-end"`
					CmsSwitchToClassicHeading    string `json:"cms-switch-to-classic-heading"`
					CmsSwitchScrolling           string `json:"cms-switch-scrolling"`
					CmsRegisteredUsersOnly       string `json:"cms-registered-users-only"`
					CmsApply                     string `json:"cms-apply"`
					CmsNoResultsDetails          string `json:"cms-no-results-details"`
					CmsShareErrorDialog          string `json:"cms-share-error-dialog"`
					CmsProductPerPage            string `json:"cms-product-per-page"`
					CmsNextPage                  string `json:"cms-next-page"`
					CmsNoResultsHelp             string `json:"cms-no-results-help"`
					CmsShareOn                   string `json:"cms-share-on"`
					CmsSeeLess                   string `json:"cms-see-less"`
					CmsRefineSearch              string `json:"cms-refine-search"`
					CmsItems                     string `json:"cms-items"`
					CmsMarketplaceMp             string `json:"cms-marketplace-mp"`
					CmsShareSocial               string `json:"cms-share-social"`
					CmsEnterQuantity             string `json:"cms-enter-quantity"`
					CmsNoSelectionsMade          string `json:"cms-no-selections-made"`
					CmsUnitPrice                 string `json:"cms-unit-price"`
					CmsTopCategories             string `json:"cms-top-categories"`
					CmsClear                     string `json:"cms-clear"`
					CmsModernSearch              string `json:"cms-modern-search"`
					CmsMarketplaceProductMessage string `json:"cms-marketplace-product-message"`
					CmsShowing                   string `json:"cms-showing"`
					CmsDetails                   string `json:"cms-details"`
					CmsFilterOptions             string `json:"cms-filter-options"`
					CmsLessFilters               string `json:"cms-less-filters"`
					CmsMissingProductPhoto       string `json:"cms-missing-product-photo"`
					CmsSwitchStacked             string `json:"cms-switch-stacked"`
					CmsShowLess                  string `json:"cms-show-less"`
					CmsViewPricesAtFilter        string `json:"cms-view-prices-at-filter"`
					CmsElectronicComponents      string `json:"cms-electronic-components"`
					CmsSwitchToClassicBody2      string `json:"cms-switch-to-classic-body-2"`
					CmsNoResults                 string `json:"cms-no-results"`
					CmsDatasheet                 string `json:"cms-datasheet"`
					CmsSearchWithinResults       string `json:"cms-search-within-results"`
					CmsNotAvailable              string `json:"cms-not-available"`
					CmsRemaining                 string `json:"cms-remaining"`
					CmsSearchEntry               string `json:"cms-search-entry"`
					CmsShareDialog               string `json:"cms-share-dialog"`
					CmsSwitchToModernHeading     string `json:"cms-switch-to-modern-heading"`
					CmsOk                        string `json:"cms-ok"`
					CmsSeeMore                   string `json:"cms-see-more"`
					CmsExactMatch                string `json:"cms-exact-match"`
					CmsSwitchToClassicBody1      string `json:"cms-switch-to-classic-body-1"`
					CmsSwitchToModernLinkStart   string `json:"cms-switch-to-modern-link-start"`
					CmsClose                     string `json:"cms-close"`
					CmsOf                        string `json:"cms-of"`
					CmsSwitchToClassicBody3      string `json:"cms-switch-to-classic-body-3"`
					CmsShare                     string `json:"cms-share"`
					CmsSearchFilter              string `json:"cms-search-filter"`
					CmsSettings                  string `json:"cms-settings"`
					CmsView                      string `json:"cms-view"`
					CmsMoreFilters               string `json:"cms-more-filters"`
					CmsResults                   string `json:"cms-results"`
					CmsSelectAnOption            string `json:"cms-select-an-option"`
					CmsNewProducts               string `json:"cms-new-products"`
					CmsShowMore                  string `json:"cms-show-more"`
					CmsResultsPerPage            string `json:"cms-results-per-page"`
					CmsSearchWithin              string `json:"cms-search-within"`
				} `json:"common"`
			} `json:"en-us"`
		} `json:"initialI18nStore"`
		InitialLanguage    string   `json:"initialLanguage"`
		NamespacesRequired []string `json:"namespacesRequired"`
	} `json:"props"`
	Page  string `json:"page"`
	Query struct {
		Num0          string `json:"0"`
		NoSiteWrapper string `json:"noSiteWrapper"`
		Any           string `json:"any"`
		ID            string `json:"id"`
	} `json:"query"`
	BuildID       string `json:"buildId"`
	AssetPrefix   string `json:"assetPrefix"`
	RuntimeConfig struct {
		APPINSIGHTSINSTRUMENTATIONKEY string `json:"APPINSIGHTS_INSTRUMENTATIONKEY"`
		BASEURL                       string `json:"BASE_URL"`
		CUSTOMENV                     string `json:"CUSTOM_ENV"`
		FEATUREFLAGDISABLELEADTIME    bool   `json:"FEATURE_FLAG_DISABLE_LEADTIME"`
		FEATUREFLAGEVERGAGE           bool   `json:"FEATURE_FLAG_EVERGAGE"`
		FEATUREFLAGFILTERSV2          bool   `json:"FEATURE_FLAG_FILTERS_V2"`
		FEATUREFLAGMOSAICCART         bool   `json:"FEATURE_FLAG_MOSAIC_CART"`
		FEATUREFLAGPRICINGCALL        bool   `json:"FEATURE_FLAG_PRICING_CALL"`
		FEATUREFLAGSETTINGSMENU       bool   `json:"FEATURE_FLAG_SETTINGS_MENU"`
	} `json:"runtimeConfig"`
	IsFallback   bool     `json:"isFallback"`
	DynamicIds   []string `json:"dynamicIds"`
	CustomServer bool     `json:"customServer"`
	Gip          bool     `json:"gip"`
	AppGip       bool     `json:"appGip"`
}

const (
	digikeyHome        = "https://www.digikey.com"
	digikeyParamSearch = "https://www.digikey.com/Search"
)

type DigikeyClient struct {
	RemoteHost   string
	client       *http.Client
	chromeDriver *webdriver.ChromeDriver
	session      *webdriver.Session
	infoCache    map[string]interface{}
}

func NewDigikeyClient() *DigikeyClient {
	hc := &DigikeyClient{
		RemoteHost: digikeyHome}
	hc.client = &http.Client{}
	hc.chromeDriver, hc.session = utils.InitChromeBrowser()
	hc.infoCache = make(map[string]interface{})
	return hc
}

func (hc *DigikeyClient) queryCallDetail(suburl string, partSpecs types.EBOMWebPart) (types.EBOMWebPart, error) {
	//var partSpecs types.EBOMWebPart
	paramString := suburl
	method := ""

	paramStringUnescaped, _ := url.QueryUnescape(paramString)
	log.Infof("Fetching: " + hc.RemoteHost + "" + method + paramStringUnescaped)
	ua := utils.GetUaHeaders()
	log.Infof("Headers: " + ua)
	//resp, err := hc.client.Get(hc.RemoteHost + "" + method + paramString)
	request, err := http.NewRequest("GET", hc.RemoteHost+""+method+paramString, nil)
	request.Header.Add("Cookie", "name=anny")
	request.Header.Add("User-Agent", ua)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return partSpecs, err
	}
	resp, err := hc.client.Do(request)
	if err != nil {
		return partSpecs, err
	}
	if resp.StatusCode != 200 {
		return partSpecs, errors.Errorf(digikeyHome + " queryCall error: " + resp.Status)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return partSpecs, err
	}

	utfBody := strings.NewReader(string(body))

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		log.Fatal(err)
	}

	// Find the key items
	content := doc.Find("#pdpMainContentDiv")

	found1st := false
	content.Find(".div-table-body .div-table-row").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		band := s.Find(".col-xs-4").Find("label").Text()
		band = strings.Replace(band, "\n", "", -1)
		band = strings.TrimSpace(band)
		title := s.Find(".col-xs-5").Text()
		title = strings.Replace(title, "\n", "", -1)
		title = strings.TrimSpace(title)
		//log.Printf("Found %d: %s - %s", i, band, title)

		if strings.HasPrefix(band, "RoHS") {
			partSpecs.RoHS = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Part Life Cycle") {
			partSpecs.Lifecycle = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Package / Case") {
			partSpecs.PackageCase = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Moisture Sensitivity") {
			partSpecs.MoistureSensitive = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Peak Reflow") {
			partSpecs.ReflowTemperaturePeak = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Minimum Operating Temperature") {
			partSpecs.OperatingTemperatureMin = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Maximum Operating Temperature") {
			partSpecs.OperatingTemperatureMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Voltage - Min") {
			partSpecs.SupplyVoltageMin = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Voltage-Max") {
			partSpecs.SupplyVoltageMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Voltage-Nom") {
			partSpecs.SupplyVoltageNom = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Current-Min") {
			partSpecs.SupplyCurrentMin = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Current-Max") {
			partSpecs.SupplyCurrentMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Operating Supply Current") {
			partSpecs.SupplyCurrentNom = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Power Dissipation-Min") {
			partSpecs.PowerDissipationMin = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Power Dissipation-Max") {
			partSpecs.PowerDissipationMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Pd - Power Dissipation") {
			partSpecs.PowerDissipationNom = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Length") {
			partSpecs.UnitLength = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Width") {
			partSpecs.UnitWidth = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Height") {
			partSpecs.UnitHeight = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Unit Weight") {
			partSpecs.UnitWeight = types.PartParameter{title, types.ParamFromDigikey}
		}
		return !found1st
	})
	return partSpecs, nil
}

func (hc *DigikeyClient) queryCall(mpn string) (types.EBOMWebPart, error) {
	var partSpecs types.EBOMWebPart
	var detaillink string
	paramString := mpn
	method := "en?keywords="

	paramStringUnescaped, _ := url.QueryUnescape(paramString)
	log.Infof("Fetching: " + hc.RemoteHost + "/products/" + method + paramStringUnescaped)
	ua := utils.GetUaHeaders()
	log.Infof("Headers: " + ua)
	request, err := http.NewRequest("GET", hc.RemoteHost+"/products/"+method+paramString, nil)
	request.Header.Add("Cookie", "name=anny")
	request.Header.Add("User-Agent", ua)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return partSpecs, err
	}
	resp, err := hc.client.Do(request)
	if err != nil {
		return partSpecs, err
	}
	if resp.StatusCode != 200 {
		return partSpecs, errors.Errorf(digikeyHome + " queryCall error: " + resp.Status)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return partSpecs, err
	}
	//log.Infof(string(body))
	utfBody := strings.NewReader(string(body))

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		log.Fatal(err)
	}

	// Find the key items
	content := doc.Find("#productTable")

	found1st := false
	content.Find("tbody tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		band := s.Find("a").Text()
		title := s.Find("i").Text()
		log.Printf("Found %d: %s - %s", i, band, title)

		s.Children().Each(func(j int, s *goquery.Selection) {
			//_val := strings.Replace(s.Text(), "\n", "", -1)
			//_val = utils.DeleteExtraSpace(_val)
			//log.Printf("Found %d-%d: %s", i, j, _val)
			switch j {
			case 2:
				detaillink, _ = s.Find("a").Attr("href")
				_val := s.Find("a").Text()
				_val = strings.Replace(_val, "\n", "", -1)
				_val = utils.DeleteExtraSpace(_val)
				_val = strings.TrimSpace(_val)
				log.Printf(_val)
				log.Printf(detaillink)
				partSpecs.MPN = types.PartParameter{_val, types.ParamFromDigikey}
			default:
			}
			found1st = true
		})
		return !found1st
	})

	if !found1st {
		return partSpecs, errors.Errorf(digikeyHome + " not found vaildate data: " + resp.Status)
	}
	//partSpecs, err = hc.queryCallDetail(detaillink, partSpecs)

	return partSpecs, nil
}

func (hc *DigikeyClient) QueryWDCall(mpn string) (types.EBOMWebPart, error) {
	var partSpecs types.EBOMWebPart
	var detaillink webdriver.WebElement
	//var cookie webdriver.Cookie
	reDigit := regexp.MustCompile("\\d*\\.?\\d+")

	paramString := mpn
	method := "en?keywords="

	paramStringUnescaped, _ := url.QueryUnescape(paramString)
	log.Infof("Fetching: " + hc.RemoteHost + "/products/" + method + paramStringUnescaped)

	//chromeDriver := hc.chromeDriver
	session := hc.session
	//chromeDriver, session := utils.InitChromeBrowser()
	//err := session.SetCookie(cookie)
	//if err != nil {
	//	return partSpecs, err
	//}

	//err := session.SetTimeouts("page load", 10000)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err := session.SetTimeoutsAsyncScript(3000)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = session.SetTimeoutsImplicitWait(5000)
	//if err != nil {
	//	log.Fatal(err)
	//}

	wh := session.GetCurrentWindowHandle()
	size, err := wh.GetSize()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("width: %d, height: %d", size.Width, size.Height)
	err = wh.MaximizeWindow()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	size2, err := wh.GetSize()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("width: %d, height: %d", size2.Width, size2.Height)

	err = session.Url(hc.RemoteHost + "/products/" + method + paramString)
	if err != nil {
		return partSpecs, err
	}
	//time.Sleep(10 * time.Second)

	// How to process Choose Your Location
	// Now Does not work, by add --disable-notifications argument | NG
	// Try click twice
	el, err := session.FindElement("tag name", "body")
	if err != nil {
		return partSpecs, err
	}
	session.MoveTo(el, 20, 600)
	session.Click(webdriver.LeftButton)
	session.Click(webdriver.LeftButton)

	// data-testid = result-page | category-page | filter-page | detail-page
	// https://stackoverflow.com/questions/57101417/find-an-element-where-data-tb-test-id-attribute-is-present-instead-of-id-using-s
	// https://www.cnblogs.com/yaoze2018/p/10387461.html
	dgkPageType := "result-page"
	var wePageType webdriver.WebElement
	wePageType, err = session.FindElement(webdriver.CSS_Selector, "div[data-testid='result-page']")
	if err != nil {
		dgkPageType = "category-page"
		wePageType, err = session.FindElement(webdriver.CSS_Selector, "div[data-testid='category-page']")
		if err != nil {
			dgkPageType = "filter-page"
			wePageType, err = session.FindElement(webdriver.CSS_Selector, "section[data-testid='filter-page']")
			if err != nil {
				return partSpecs, errors.Errorf(digikeyHome + " get unknown page")
			}
		}
	}
	log.Println(dgkPageType)

	if dgkPageType == "result-page" {
		//div, err := wePageType.FindElement(webdriver.CSS_Selector, "div>div>div>div:nth-child(2)")
		//if err != nil {
		//	return partSpecs, err
		//}
		//log.Println(div.GetAttribute("innerHTML"))
		//section, err := div.FindElement(webdriver.CSS_Selector, "section:nth-child(6)") // 6th elem is <section>
		//section, err := div.FindElement(webdriver.CSS_Selector, "div:nth-child(3)") // 3th elem is <div>
		//if err != nil {
		//	return partSpecs, err
		//}
		//log.Println(section.GetAttribute("innerHTML"))

		tbody, err := wePageType.FindElement(webdriver.CSS_Selector, "tbody.MuiTableBody-root")
		if err != nil {
			return partSpecs, err
		}
		trs, err := tbody.FindElements(webdriver.TagName, "tr")
		if err != nil {
			return partSpecs, err
		}

		td, err := trs[0].FindElement(webdriver.CSS_Selector, "td")
		if err != nil {
			return partSpecs, err
		}
		hrefs, err := td.FindElements(webdriver.CSS_Selector, "a")
		if err != nil {
			return partSpecs, err
		}
		href, err := hrefs[0].GetAttribute("href")
		if err != nil {
			return partSpecs, err
		}
		log.Printf(href)

		err = hrefs[0].Click()
		if err != nil {
			return partSpecs, err
		}
		time.Sleep(2 * time.Second)
		// expect to get filter-page
	}
	if dgkPageType == "category-page" {
		section, err := wePageType.FindElement(webdriver.CSS_Selector, "section>div>ul")
		if err != nil {
			return partSpecs, err
		}
		hrefs, err := section.FindElements(webdriver.CSS_Selector, "a")
		if err != nil {
			return partSpecs, err
		}
		href, err := hrefs[0].GetAttribute("href")
		if err != nil {
			return partSpecs, err
		}
		log.Printf(href)

		err = hrefs[0].Click()
		if err != nil {
			return partSpecs, err
		}
		time.Sleep(2 * time.Second)
		// expect to get filter-page
	}

	time.Sleep(1 * time.Second)
	session.ExecuteScript("window.scrollBy(0, 400)", make([]interface{}, 0))
	wePageType, err = session.FindElement(webdriver.CSS_Selector, "section[data-testid='filter-page']")
	if err != nil {
		return partSpecs, errors.Errorf(digikeyHome + " not expect get filter-page")
	}

	tbody, err := wePageType.FindElement(webdriver.CSS_Selector, "tbody.MuiTableBody-root")
	if err != nil {
		return partSpecs, err
	}
	trs, err := tbody.FindElements(webdriver.TagName, "tr")
	if err != nil {
		return partSpecs, err
	}
	for _, tr := range trs {
		tds, err := tr.FindElements(webdriver.TagName, "td")
		if err != nil {
			return partSpecs, err
		}

		// Discontinued at Digi-Key
		_PartStatus, err := tds[10].Text()
		if err != nil {
			return partSpecs, err
		}
		if strings.HasPrefix(_PartStatus, "Discontinued") {
			continue
		}

		hrefs, err := tds[1].FindElements(webdriver.CSS_Selector, "a")
		if err != nil {
			return partSpecs, err
		}
		href, err := hrefs[1].GetAttribute("href")
		if err != nil {
			return partSpecs, err
		}
		log.Printf(href)
		detaillink = hrefs[1]
		break
	}

	err = detaillink.Click()
	if err != nil {
		return partSpecs, err
	}
	time.Sleep(2 * time.Second)
	// expect to get detail-page

	session.ExecuteScript("window.scrollBy(0, 400)", make([]interface{}, 0))
	wePageType, err = session.FindElement(webdriver.CSS_Selector, "div[data-testid='detail-page']")
	if err != nil {
		return partSpecs, errors.Errorf(digikeyHome + " not expect get detail-page")
	}

	// https://mholt.github.io/json-to-go/
	// <script id="__NEXT_DATA__">
	session.ExecuteScript("window.scrollBy(0, 400)", make([]interface{}, 0))
	session.Refresh()
	time.Sleep(2 * time.Second)
	ngDgkDataXML, err := session.FindElement(webdriver.CSS_Selector, "#__NEXT_DATA__")
	if err != nil {
		return partSpecs, err
	}
	ngDgkDataXMLText, _ := ngDgkDataXML.GetAttribute("outerHTML")
	//log.Println(ngDgkDataXMLText)
	//log.Println(ngDgkDataXML.GetAttribute("outerHTML"))
	re, _ := regexp.Compile(`<script.*?>(.*)</script>`)
	ngDgkDataXMLText = re.ReplaceAllString(ngDgkDataXMLText, "$1")
	if false {
		//err = ioutil.WriteFile("dgk_newdata.json", []byte(ngDgkDataXMLText), os.ModeAppend)
		fl, err := os.OpenFile("dgk_newdata.json", os.O_CREATE, 0644) // os.O_APPEND
		if err != nil {
			return partSpecs, err
		}
		defer fl.Close()
		n, err := fl.Write([]byte(ngDgkDataXMLText))
		if err == nil && n < len([]byte(ngDgkDataXMLText)) {
			return partSpecs, err
		}
	}

	ngDgkDataJSON := &NgDgkData{}
	err = json.Unmarshal([]byte(ngDgkDataXMLText), &ngDgkDataJSON)
	if err != nil {
		log.Println("Umarshal failed:", err)
		return partSpecs, err
	}

	PageProps := ngDgkDataJSON.Props.PageProps
	productAttributes := PageProps.Envelope.Data.ProductAttributes.Attributes

	for _, attr := range productAttributes {
		band := attr.Label
		title := attr.Values[0].Value
		if strings.HasPrefix(band, "RoHS") {
			partSpecs.RoHS = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Part Status") {
			partSpecs.Lifecycle = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Package / Case") {
			partSpecs.PackageCase = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Moisture Sensitivity") {
			partSpecs.MoistureSensitive = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Peak Reflow") {
			partSpecs.ReflowTemperaturePeak = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Operating Temperature") {
			if strings.Contains(title, "~") {
				_val := strings.Split(title, "~")
				baseval := string(reDigit.FindAll([]byte(_val[0]), -1)[0])
				partSpecs.OperatingTemperatureMin = types.PartParameter{baseval, types.ParamFromDigikey}
				baseval = string(reDigit.FindAll([]byte(_val[1]), -1)[0])
				partSpecs.OperatingTemperatureMax = types.PartParameter{baseval, types.ParamFromDigikey}
			} else {
				dval := reDigit.FindAll([]byte(title), -1)
				baseval := ""
				if len(dval) >= 1 {
					baseval = string(dval[0])
				}
				partSpecs.OperatingTemperatureMin = types.PartParameter{baseval, types.ParamFromDigikey}
			}
		} else if strings.HasPrefix(band, "Maximum Operating Temperature") {
			partSpecs.OperatingTemperatureMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Voltage - Min") {
			partSpecs.SupplyVoltageMin = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Voltage-Max") {
			partSpecs.SupplyVoltageMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Voltage - Supply") {
			if strings.Contains(title, "~") {
				_val := strings.Split(title, "~")
				baseval := string(reDigit.FindAll([]byte(_val[0]), -1)[0])
				partSpecs.SupplyVoltageMin = types.PartParameter{baseval, types.ParamFromDigikey}
				baseval = string(reDigit.FindAll([]byte(_val[1]), -1)[0])
				partSpecs.SupplyVoltageMax = types.PartParameter{baseval, types.ParamFromDigikey}
			} else {
				dval := reDigit.FindAll([]byte(title), -1)
				baseval := ""
				if len(dval) >= 1 {
					baseval = string(dval[0])
				}
				partSpecs.SupplyVoltageMin = types.PartParameter{baseval, types.ParamFromDigikey}
			}
		} else if strings.HasPrefix(band, "Supply Current-Min") {
			partSpecs.SupplyCurrentMin = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Current-Max") {
			partSpecs.SupplyCurrentMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Current - Supply") {
			partSpecs.SupplyCurrentNom = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Power Dissipation-Min") {
			partSpecs.PowerDissipationMin = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Power Dissipation-Max") {
			partSpecs.PowerDissipationMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Pd - Power Dissipation") {
			partSpecs.PowerDissipationNom = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Length") {
			partSpecs.UnitLength = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Width") {
			partSpecs.UnitWidth = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Height") {
			partSpecs.UnitHeight = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Unit Weight") {
			partSpecs.UnitWeight = types.PartParameter{title, types.ParamFromDigikey}
		}
	}

	Pricing := PageProps.Envelope.Data.PriceQuantity.Pricing
	PricingTiers := Pricing[0].PricingTiers

	valPrice := ""
	lastPrice := ""
	for _, pricing := range PricingTiers {
		qty, err := strconv.Atoi(strings.ReplaceAll(pricing.BreakQty, ",", ""))
		if err != nil {
			return partSpecs, err
		}
		if qty <= 1000 {
			valPrice = pricing.UnitPrice
		} else if valPrice == "" {
			valPrice = pricing.UnitPrice
		} else {
			break
		}
		lastPrice = pricing.UnitPrice
		//log.Println(lastPrice)
	}
	if valPrice == "" {
		valPrice = lastPrice
	}
	log.Println(valPrice)
	partSpecs.UnitPrice = types.PartParameter{valPrice, types.ParamFromDigikey}

	//session.Delete()
	//chromeDriver.Stop()

	return partSpecs, nil
}

func (hc *DigikeyClient) QueryWDCall2(mpn string) (types.EBOMWebPart, error) {
	var partSpecs types.EBOMWebPart
	//var detaillink webdriver.WebElement
	//var cookie webdriver.Cookie
	reDigit := regexp.MustCompile("\\d*\\.?\\d+")

	paramString := mpn
	method := "en?keywords="

	paramStringUnescaped, _ := url.QueryUnescape(paramString)
	log.Infof("Fetching: " + hc.RemoteHost + "/products/" + method + paramStringUnescaped)

	//chromeDriver := hc.chromeDriver
	session := hc.session
	//chromeDriver, session := utils.InitChromeBrowser()
	//err := session.SetCookie(cookie)
	//if err != nil {
	//	return partSpecs, err
	//}

	//err := session.SetTimeouts("page load", 10000)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err := session.SetTimeoutsAsyncScript(3000)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = session.SetTimeoutsImplicitWait(5000)
	//if err != nil {
	//	log.Fatal(err)
	//}

	wh := session.GetCurrentWindowHandle()
	size, err := wh.GetSize()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("width: %d, height: %d", size.Width, size.Height)
	err = wh.MaximizeWindow()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	size2, err := wh.GetSize()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("width: %d, height: %d", size2.Width, size2.Height)

	err = session.Url(hc.RemoteHost + "/products/" + method + paramString)
	if err != nil {
		return partSpecs, err
	}
	//time.Sleep(10 * time.Second)

	// How to process Choose Your Location
	// Now Does not work, by add --disable-notifications argument
	_, err = session.GetAlertText()
	if err == nil {
		err = session.DismissAlert()
	}

	// data-testid = result-page | category-page | filter-page | detail-page
	// https://stackoverflow.com/questions/57101417/find-an-element-where-data-tb-test-id-attribute-is-present-instead-of-id-using-s
	// https://www.cnblogs.com/yaoze2018/p/10387461.html
	dgkPageType := "result-page"
	var wePageType webdriver.WebElement
	wePageType, err = session.FindElement(webdriver.CSS_Selector, "div[data-testid='result-page']")
	if err != nil {
		dgkPageType = "category-page"
		wePageType, err = session.FindElement(webdriver.CSS_Selector, "div[data-testid='category-page']")
		if err != nil {
			dgkPageType = "filter-page"
			wePageType, err = session.FindElement(webdriver.CSS_Selector, "section[data-testid='filter-page']")
			if err != nil {
				return partSpecs, errors.Errorf(digikeyHome + " get unknown page")
			}
		}
	}
	log.Println(dgkPageType)

	if dgkPageType == "result-page" {
		// expect to get filter-page
	}
	if dgkPageType == "category-page" {
		section, err := wePageType.FindElement(webdriver.CSS_Selector, "section>div>ul")
		if err != nil {
			return partSpecs, err
		}
		hrefs, err := section.FindElements(webdriver.CSS_Selector, "a")
		if err != nil {
			return partSpecs, err
		}
		href, err := hrefs[0].GetAttribute("href")
		if err != nil {
			return partSpecs, err
		}
		log.Printf(href)

		err = hrefs[0].Click()
		if err != nil {
			return partSpecs, err
		}
		time.Sleep(2 * time.Second)
		// expect to get filter-page
	}

	time.Sleep(1 * time.Second)
	session.ExecuteScriptAsync("window.scrollBy(0, 400)", make([]interface{}, 0))
	wePageType, err = session.FindElement(webdriver.CSS_Selector, "section[data-testid='filter-page']")
	if err != nil {
		return partSpecs, errors.Errorf(digikeyHome + " not expect get filter-page")
	}

	tbody, err := wePageType.FindElement(webdriver.CSS_Selector, "tbody.MuiTableBody-root")
	if err != nil {
		return partSpecs, err
	}
	trs, err := tbody.FindElements(webdriver.TagName, "tr")
	if err != nil {
		return partSpecs, err
	}
	tds, err := trs[0].FindElements(webdriver.TagName, "td")
	if err != nil {
		return partSpecs, err
	}
	hrefs, err := tds[1].FindElements(webdriver.CSS_Selector, "a")
	if err != nil {
		return partSpecs, err
	}
	href, err := hrefs[2].GetAttribute("href")
	if err != nil {
		return partSpecs, err
	}
	log.Printf(href)

	err = hrefs[2].Click()
	if err != nil {
		return partSpecs, err
	}
	time.Sleep(2 * time.Second)
	// expect to get detail-page

	session.ExecuteScriptAsync("window.scrollBy(0, 400)", make([]interface{}, 0))
	wePageType, err = session.FindElement(webdriver.CSS_Selector, "div[data-testid='detail-page']")
	if err != nil {
		return partSpecs, errors.Errorf(digikeyHome + " not expect get detail-page")
	}

	// time.Sleep(2 * time.Second)
	// #pdp_content
	//   #product-photo
	//   product-details-procurement
	//   product-details-overview
	//     product-overview-photo-spacer
	//     product-details-documents-media product-details-section
	//     product-details-product-attributes product-details-section
	//     product-details-environmental-export product-details-section
	// <table class="MuiTable-root" id="product-attributes">
	session.ExecuteScript("window.scrollBy(0, 600)", make([]interface{}, 0))
	prodAttr, err := session.FindElement(webdriver.ID, "product-attributes")
	if err != nil {
		return partSpecs, err
	}
	prodAttrTrs, err := prodAttr.FindElements(webdriver.TagName, "tr")
	if err != nil {
		return partSpecs, err
	}
	for k, trv := range prodAttrTrs {
		//log.Println(k, trv)
		if err != nil {
			continue
		}
		if k == 0 {
			continue
		}
		_val, err := trv.FindElement(webdriver.TagName, "th")
		if err != nil {
			continue
		}
		band, err := _val.Text()
		if err != nil {
			continue
		}
		//band = strings.Replace(band, "\n", "", -1)
		//band = strings.TrimSpace(band)

		tds, err := trv.FindElements(webdriver.TagName, "td")
		if err != nil {
			continue
		}
		title, err := tds[0].Text()
		if err != nil {
			continue
		}
		if strings.HasPrefix(band, "RoHS") {
			partSpecs.RoHS = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Part Status") {
			partSpecs.Lifecycle = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Package / Case") {
			partSpecs.PackageCase = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Moisture Sensitivity") {
			partSpecs.MoistureSensitive = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Peak Reflow") {
			partSpecs.ReflowTemperaturePeak = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Operating Temperature") {
			if strings.Contains(title, "~") {
				_val := strings.Split(title, "~")
				baseval := string(reDigit.FindAll([]byte(_val[0]), -1)[0])
				partSpecs.OperatingTemperatureMin = types.PartParameter{baseval, types.ParamFromDigikey}
				baseval = string(reDigit.FindAll([]byte(_val[1]), -1)[0])
				partSpecs.OperatingTemperatureMax = types.PartParameter{baseval, types.ParamFromDigikey}
			} else {
				dval := reDigit.FindAll([]byte(title), -1)
				baseval := ""
				if len(dval) >= 1 {
					baseval = string(dval[0])
				}
				partSpecs.OperatingTemperatureMin = types.PartParameter{baseval, types.ParamFromDigikey}
			}
		} else if strings.HasPrefix(band, "Maximum Operating Temperature") {
			partSpecs.OperatingTemperatureMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Voltage - Min") {
			partSpecs.SupplyVoltageMin = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Voltage-Max") {
			partSpecs.SupplyVoltageMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Voltage - Supply") {
			if strings.Contains(title, "~") {
				_val := strings.Split(title, "~")
				baseval := string(reDigit.FindAll([]byte(_val[0]), -1)[0])
				partSpecs.SupplyVoltageMin = types.PartParameter{baseval, types.ParamFromDigikey}
				baseval = string(reDigit.FindAll([]byte(_val[1]), -1)[0])
				partSpecs.SupplyVoltageMax = types.PartParameter{baseval, types.ParamFromDigikey}
			} else {
				dval := reDigit.FindAll([]byte(title), -1)
				baseval := ""
				if len(dval) >= 1 {
					baseval = string(dval[0])
				}
				partSpecs.SupplyVoltageMin = types.PartParameter{baseval, types.ParamFromDigikey}
			}
		} else if strings.HasPrefix(band, "Supply Current-Min") {
			partSpecs.SupplyCurrentMin = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Current-Max") {
			partSpecs.SupplyCurrentMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Current - Supply") {
			partSpecs.SupplyCurrentNom = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Power Dissipation-Min") {
			partSpecs.PowerDissipationMin = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Power Dissipation-Max") {
			partSpecs.PowerDissipationMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Pd - Power Dissipation") {
			partSpecs.PowerDissipationNom = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Length") {
			partSpecs.UnitLength = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Width") {
			partSpecs.UnitWidth = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Height") {
			partSpecs.UnitHeight = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Unit Weight") {
			partSpecs.UnitWeight = types.PartParameter{title, types.ParamFromDigikey}
		}
	}

	session.Refresh()
	time.Sleep(2 * time.Second)
	// https://mholt.github.io/json-to-go/
	// <script id="__NEXT_DATA__">
	prodPrice, err := session.FindElement(webdriver.CSS_Selector, "#__NEXT_DATA__")
	//prodPrice, err := session.FindElement(webdriver.ClassName, "product-dollars")
	if err != nil {
		return partSpecs, err
	}
	//prodPriceText, _ := prodPrice.GetAttribute("innerHTML")
	prodPriceText, _ := prodPrice.GetAttribute("outerHTML")
	//log.Println(prodPriceText)
	//log.Println(prodPrice.GetAttribute("outerHTML"))
	re, _ := regexp.Compile(`<script.*?>(.*)</script>`)
	prodPriceText = re.ReplaceAllString(prodPriceText, "$1")

	prodPriceMap := &NgDgkData{}
	err = json.Unmarshal([]byte(prodPriceText), &prodPriceMap)
	if err != nil {
		log.Println("Umarshal failed:", err)
		return partSpecs, err
	}
	PageProps := prodPriceMap.Props.PageProps
	Pricing := PageProps.Envelope.Data.PriceQuantity.Pricing
	PricingTiers := Pricing[0].PricingTiers

	valPrice := ""
	lastPrice := ""
	for _, pricing := range PricingTiers {
		qty, err := strconv.Atoi(strings.ReplaceAll(pricing.BreakQty, ",", ""))
		if err != nil {
			return partSpecs, err
		}
		if qty <= 1000 {
			valPrice = pricing.UnitPrice
		} else if valPrice == "" {
			valPrice = pricing.UnitPrice
		} else {
			break
		}
		lastPrice = pricing.UnitPrice
		//log.Println(lastPrice)
	}
	if valPrice == "" {
		valPrice = lastPrice
	}
	log.Println(valPrice)
	partSpecs.UnitPrice = types.PartParameter{valPrice, types.ParamFromDigikey}

	//session.Delete()
	//chromeDriver.Stop()

	return partSpecs, nil
}

func (hc *DigikeyClient) Close() {
	hc.session.Delete()
	hc.chromeDriver.Stop()
}
