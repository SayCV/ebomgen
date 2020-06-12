package reliability

// Attribute defines a single Code
type Attribute int

// Base attributes
const (
	RpmPartsCount Attribute = iota
	RpmPartStress
	RpmPseudoStress
	RpmUndef
)

// ClsEnv - Environment Class
var (
	ClsEnv = []struct {
		mark     string
		title    string
		desc     string
		temp     string
	}{
		{
			"GB",
			"Ground, Fixed, Controlled",
			"",
			"30",
		},
		{
			"GMS",
			"Ground, Fixed, Controlled",
			"",
			"30",
		},
		{
			"GF1",
			"Ground, Fixed, Uncontrolled",
			"",
			"40",
		},
		{
			"GF2",
			"Ground, Fixed, Uncontrolled",
			"",
			"40",
		},
		{
			"GM1",
			"Ground, Mobile, Steady",
			"",
			"55",
		},
		{
			"GM2",
			"Ground, Mobile, Violent",
			"",
			"60",
		},
		{
			"MP",
			"Ground, Fixed, Controlled",
			"",
			"40",
		},
		{
			"NSB",
			"Ground, Fixed, Controlled",
			"",
			"45",
		},
		{
			"NS1",
			"Ground, Fixed, Controlled",
			"",
			"40",
		},
		{
			"NS2",
			"Ground, Fixed, Controlled",
			"",
			"45",
		},
		{
			"NU",
			"Ground, Fixed, Controlled",
			"",
			"70",
		},
		{
			"AIF",
			"Ground, Fixed, Controlled",
			"",
			"55",
		},
		{
			"AUF",
			"Ground, Fixed, Controlled",
			"",
			"70",
		},
		{
			"AIC",
			"Airborne, Commercial",
			"",
			"55",
		},
		{
			"AUC",
			"Ground, Fixed, Controlled",
			"",
			"70",
		},
		{
			"ARW",
			"Ground, Fixed, Controlled",
			"",
			"55",
		},
		{
			"SF",
			"Ground, Fixed, Controlled",
			"",
			"30",
		},
		{
			"ML",
			"Ground, Fixed, Controlled",
			"",
			"55",
		},
		{
			"MF",
			"Ground, Fixed, Controlled",
			"",
			"55",
		},
	}
)

// ClsEnv - Environment Factor
var (
	FactorEnv = [2]map[string]interface{} {
		"IC":  {[]string{"1.0", "1.2", "1.8", "4.2", "4.1", "7.0", "4.6", "4.5", "3.0", "6.0", "10", "9.5", "13", "5.8", "10.0", "11", "1.0", "18", "9"}},
		"IC1": []string{"1.0", "1.2", "1.8", "4.2", "4.1", "7.0", "4.6", "4.5", "3.0", "6.0", "10", "9.5", "13", "5.8", "10.0", "11", "1.0", "18", "9"},
	}
)
