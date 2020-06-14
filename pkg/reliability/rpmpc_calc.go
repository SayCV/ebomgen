package reliability

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (b *EBOMFrPart) FrCalcCap() (string, error) {
	strlambda_b, _ := b.GetFailureRateBaseImported()
	strpi_e, _ := b.GetFactorEnvImported()
	strpi_q, _ := b.GetFactorQualityImported()
	strpi_t, _ := b.GetFactorTemperatureImported()
	strpi_s, _ := b.GetFactorStressImported()
	strpi_ch, _ := b.GetFactorChImported()

	strlambda_b = strings.Replace(strlambda_b, " ", "", -1)
	strpi_e = strings.Replace(strpi_e, " ", "", -1)
	strpi_q = strings.Replace(strpi_q, " ", "", -1)
	strpi_t = strings.Replace(strpi_t, " ", "", -1)
	strpi_s = strings.Replace(strpi_s, " ", "", -1)
	strpi_ch = strings.Replace(strpi_ch, " ", "", -1)

	lambda_b, _ := strconv.ParseFloat(strlambda_b, 64)
	pi_e, _ := strconv.ParseFloat(strpi_e, 64)
	pi_q, _ := strconv.ParseFloat(strpi_q, 64)
	pi_t, _ := strconv.ParseFloat(strpi_t, 64)
	pi_s, _ := strconv.ParseFloat(strpi_s, 64)
	pi_ch, _ := strconv.ParseFloat(strpi_ch, 64)

	reqValue := lambda_b * pi_e * pi_q * pi_t * pi_s * pi_ch
	strreqValue := strconv.FormatFloat(reqValue, 'f', -1, 64)
	return strreqValue, nil
}

func (b *EBOMFrPart) FrCalcRes() (string, error) {
	strlambda_b, _ := b.GetFailureRateBaseImported()
	strpi_e, _ := b.GetFactorEnvImported()
	strpi_q, _ := b.GetFactorQualityImported()
	strpi_t, _ := b.GetFactorTemperatureImported()
	strpi_s, _ := b.GetFactorStressImported()
	strpi_ch := "1.0"

	strlambda_b = strings.Replace(strlambda_b, " ", "", -1)
	strpi_e = strings.Replace(strpi_e, " ", "", -1)
	strpi_q = strings.Replace(strpi_q, " ", "", -1)
	strpi_t = strings.Replace(strpi_t, " ", "", -1)
	strpi_s = strings.Replace(strpi_s, " ", "", -1)
	strpi_ch = strings.Replace(strpi_ch, " ", "", -1)

	lambda_b, _ := strconv.ParseFloat(strlambda_b, 64)
	pi_e, _ := strconv.ParseFloat(strpi_e, 64)
	pi_q, _ := strconv.ParseFloat(strpi_q, 64)
	pi_t, _ := strconv.ParseFloat(strpi_t, 64)
	pi_s, _ := strconv.ParseFloat(strpi_s, 64)
	pi_ch, _ := strconv.ParseFloat(strpi_ch, 64)

	reqValue := lambda_b * pi_e * pi_q * pi_t * pi_s * pi_ch
	strreqValue := strconv.FormatFloat(reqValue, 'f', -1, 64)
	return strreqValue, nil
}

func (b *EBOMFrPart) FrCalcInd() (string, error) {
	strlambda_b, _ := b.GetFailureRateBaseImported()
	strpi_e, _ := b.GetFactorEnvImported()
	strpi_q, _ := b.GetFactorQualityImported()
	strpi_t, _ := b.GetFactorTemperatureImported()
	strpi_s := "1.0"
	strpi_ch := "1.0"

	strlambda_b = strings.Replace(strlambda_b, " ", "", -1)
	strpi_e = strings.Replace(strpi_e, " ", "", -1)
	strpi_q = strings.Replace(strpi_q, " ", "", -1)
	strpi_t = strings.Replace(strpi_t, " ", "", -1)
	strpi_s = strings.Replace(strpi_s, " ", "", -1)
	strpi_ch = strings.Replace(strpi_ch, " ", "", -1)

	lambda_b, _ := strconv.ParseFloat(strlambda_b, 64)
	pi_e, _ := strconv.ParseFloat(strpi_e, 64)
	pi_q, _ := strconv.ParseFloat(strpi_q, 64)
	pi_t, _ := strconv.ParseFloat(strpi_t, 64)
	pi_s, _ := strconv.ParseFloat(strpi_s, 64)
	pi_ch, _ := strconv.ParseFloat(strpi_ch, 64)

	reqValue := lambda_b * pi_e * pi_q * pi_t * pi_s * pi_ch
	strreqValue := strconv.FormatFloat(reqValue, 'f', -1, 64)
	return strreqValue, nil
}

func (b *EBOMFrPart) FrCalcDiodeBjt() (string, error) {
	strlambda_b, _ := b.GetFailureRateBaseImported()
	strpi_e, _ := b.GetFactorEnvImported()
	strpi_q, _ := b.GetFactorQualityImported()
	strpi_t, _ := b.GetFactorTemperatureImported()
	strpi_s, _ := b.GetFactorStressImported()
	strpi_ch := "1.0"

	strlambda_b = strings.Replace(strlambda_b, " ", "", -1)
	strpi_e = strings.Replace(strpi_e, " ", "", -1)
	strpi_q = strings.Replace(strpi_q, " ", "", -1)
	strpi_t = strings.Replace(strpi_t, " ", "", -1)
	strpi_s = strings.Replace(strpi_s, " ", "", -1)
	strpi_ch = strings.Replace(strpi_ch, " ", "", -1)

	lambda_b, _ := strconv.ParseFloat(strlambda_b, 64)
	pi_e, _ := strconv.ParseFloat(strpi_e, 64)
	pi_q, _ := strconv.ParseFloat(strpi_q, 64)
	pi_t, _ := strconv.ParseFloat(strpi_t, 64)
	pi_s, _ := strconv.ParseFloat(strpi_s, 64)
	pi_ch, _ := strconv.ParseFloat(strpi_ch, 64)

	log.Info(lambda_b, pi_e, pi_q, pi_t, pi_s, pi_ch)

	reqValue := lambda_b * pi_e * pi_q * pi_t * pi_s * pi_ch
	strreqValue := strconv.FormatFloat(reqValue, 'f', -1, 64)
	return strreqValue, nil
}

func (b *EBOMFrPart) FrCalcOptoElectronicDevices() (string, error) {
	strlambda_b, _ := b.GetFailureRateBaseImported()
	strpi_e, _ := b.GetFactorEnvImported()
	strpi_q, _ := b.GetFactorQualityImported()
	strpi_t, _ := b.GetFactorTemperatureImported()
	strpi_s := "1.0"
	strpi_ch := "1.0"

	strlambda_b = strings.Replace(strlambda_b, " ", "", -1)
	strpi_e = strings.Replace(strpi_e, " ", "", -1)
	strpi_q = strings.Replace(strpi_q, " ", "", -1)
	strpi_t = strings.Replace(strpi_t, " ", "", -1)
	strpi_s = strings.Replace(strpi_s, " ", "", -1)
	strpi_ch = strings.Replace(strpi_ch, " ", "", -1)

	lambda_b, _ := strconv.ParseFloat(strlambda_b, 64)
	pi_e, _ := strconv.ParseFloat(strpi_e, 64)
	pi_q, _ := strconv.ParseFloat(strpi_q, 64)
	pi_t, _ := strconv.ParseFloat(strpi_t, 64)
	pi_s, _ := strconv.ParseFloat(strpi_s, 64)
	pi_ch, _ := strconv.ParseFloat(strpi_ch, 64)

	reqValue := lambda_b * pi_e * pi_q * pi_t * pi_s * pi_ch
	strreqValue := strconv.FormatFloat(reqValue, 'f', -1, 64)
	return strreqValue, nil
}

func (b *EBOMFrPart) FrCalcXtal() (string, error) {
	strlambda_b, _ := b.GetFailureRateBaseImported()
	strpi_e, _ := b.GetFactorEnvImported()
	strpi_q, _ := b.GetFactorQualityImported()
	strpi_t := "1.0"
	strpi_s := "1.0"
	strpi_ch := "1.0"

	strlambda_b = strings.Replace(strlambda_b, " ", "", -1)
	strpi_e = strings.Replace(strpi_e, " ", "", -1)
	strpi_q = strings.Replace(strpi_q, " ", "", -1)
	strpi_t = strings.Replace(strpi_t, " ", "", -1)
	strpi_s = strings.Replace(strpi_s, " ", "", -1)
	strpi_ch = strings.Replace(strpi_ch, " ", "", -1)

	lambda_b, _ := strconv.ParseFloat(strlambda_b, 64)
	pi_e, _ := strconv.ParseFloat(strpi_e, 64)
	pi_q, _ := strconv.ParseFloat(strpi_q, 64)
	pi_t, _ := strconv.ParseFloat(strpi_t, 64)
	pi_s, _ := strconv.ParseFloat(strpi_s, 64)
	pi_ch, _ := strconv.ParseFloat(strpi_ch, 64)

	reqValue := lambda_b * pi_e * pi_q * pi_t * pi_s * pi_ch
	strreqValue := strconv.FormatFloat(reqValue, 'f', -1, 64)
	return strreqValue, nil
}

func (b *EBOMFrPart) FrCalcRelay() (string, error) {
	strlambda_b, _ := b.GetFailureRateBaseImported()
	strpi_e, _ := b.GetFactorEnvImported()
	strpi_q, _ := b.GetFactorQualityImported()
	strpi_t, _ := b.GetFactorTemperatureImported()
	strpi_s, _ := b.GetFactorStressImported()
	strpi_ch := "1.0"

	strlambda_b = strings.Replace(strlambda_b, " ", "", -1)
	strpi_e = strings.Replace(strpi_e, " ", "", -1)
	strpi_q = strings.Replace(strpi_q, " ", "", -1)
	strpi_t = strings.Replace(strpi_t, " ", "", -1)
	strpi_s = strings.Replace(strpi_s, " ", "", -1)
	strpi_ch = strings.Replace(strpi_ch, " ", "", -1)

	lambda_b, _ := strconv.ParseFloat(strlambda_b, 64)
	pi_e, _ := strconv.ParseFloat(strpi_e, 64)
	pi_q, _ := strconv.ParseFloat(strpi_q, 64)
	pi_t, _ := strconv.ParseFloat(strpi_t, 64)
	pi_s, _ := strconv.ParseFloat(strpi_s, 64)
	pi_ch, _ := strconv.ParseFloat(strpi_ch, 64)

	reqValue := lambda_b * pi_e * pi_q * pi_t * pi_s * pi_ch
	strreqValue := strconv.FormatFloat(reqValue, 'f', -1, 64)
	return strreqValue, nil
}

func (b *EBOMFrPart) FrCalcSwitch() (string, error) {
	strlambda_b1, _ := b.GetFailureRateBaseImported()
	strlambda_b2, _ := b.GetFailureRateActiveContactImported(2)
	strpi_e, _ := b.GetFactorEnvImported()
	strpi_q, _ := b.GetFactorQualityImported()
	strpi_t, _ := b.GetFactorTemperatureImported()
	strpi_s, _ := b.GetFactorStressImported()
	strpi_ch := "1.0"

	strlambda_b1 = strings.Replace(strlambda_b1, " ", "", -1)
	strlambda_b2 = strings.Replace(strlambda_b2, " ", "", -1)
	strpi_e = strings.Replace(strpi_e, " ", "", -1)
	strpi_q = strings.Replace(strpi_q, " ", "", -1)
	strpi_t = strings.Replace(strpi_t, " ", "", -1)
	strpi_s = strings.Replace(strpi_s, " ", "", -1)
	strpi_ch = strings.Replace(strpi_ch, " ", "", -1)

	lambda_b1, _ := strconv.ParseFloat(strlambda_b1, 64)
	lambda_b2, _ := strconv.ParseFloat(strlambda_b2, 64)
	pi_e, _ := strconv.ParseFloat(strpi_e, 64)
	pi_q, _ := strconv.ParseFloat(strpi_q, 64)
	pi_t, _ := strconv.ParseFloat(strpi_t, 64)
	pi_s, _ := strconv.ParseFloat(strpi_s, 64)
	pi_ch, _ := strconv.ParseFloat(strpi_ch, 64)

	reqValue := (lambda_b1 + lambda_b2) * pi_e * pi_q * pi_t * pi_s * pi_ch
	strreqValue := strconv.FormatFloat(reqValue, 'f', -1, 64)
	return strreqValue, nil
}

func (b *EBOMFrPart) FrCalcConn() (string, error) {
	strlambda_b, _ := b.GetFailureRateBaseImported()
	strpi_e, _ := b.GetFactorEnvImported()
	strpi_q, _ := b.GetFactorQualityImported()
	strpi_t, _ := b.GetFactorTemperatureImported()
	strpi_s := "1.0"
	strpi_ch := "1.0"
	strpi_p, _ := b.GetFactorProcessImported()

	strlambda_b = strings.Replace(strlambda_b, " ", "", -1)
	strpi_e = strings.Replace(strpi_e, " ", "", -1)
	strpi_q = strings.Replace(strpi_q, " ", "", -1)
	strpi_t = strings.Replace(strpi_t, " ", "", -1)
	strpi_s = strings.Replace(strpi_s, " ", "", -1)
	strpi_ch = strings.Replace(strpi_ch, " ", "", -1)
	strpi_p = strings.Replace(strpi_p, " ", "", -1)

	lambda_b, _ := strconv.ParseFloat(strlambda_b, 64)
	pi_e, _ := strconv.ParseFloat(strpi_e, 64)
	pi_q, _ := strconv.ParseFloat(strpi_q, 64)
	pi_t, _ := strconv.ParseFloat(strpi_t, 64)
	pi_s, _ := strconv.ParseFloat(strpi_s, 64)
	pi_ch, _ := strconv.ParseFloat(strpi_ch, 64)
	pi_p, _ := strconv.ParseFloat(strpi_p, 64)

	reqValue := lambda_b * pi_e * pi_q * pi_t * pi_s * pi_ch * pi_p
	strreqValue := strconv.FormatFloat(reqValue, 'f', -1, 64)
	return strreqValue, nil
}

func (b *EBOMFrPart) FrCalcRotator() (string, error) {
	strlambda_b, _ := b.GetFailureRateBaseImported()
	strpi_e, _ := b.GetFactorEnvImported()
	strpi_q, _ := b.GetFactorQualityImported()
	strpi_t := "1.0"
	strpi_s := "1.0"
	strpi_ch := "1.0"
	strpi_p := "1.0"

	strlambda_b = strings.Replace(strlambda_b, " ", "", -1)
	strpi_e = strings.Replace(strpi_e, " ", "", -1)
	strpi_q = strings.Replace(strpi_q, " ", "", -1)
	strpi_t = strings.Replace(strpi_t, " ", "", -1)
	strpi_s = strings.Replace(strpi_s, " ", "", -1)
	strpi_ch = strings.Replace(strpi_ch, " ", "", -1)
	strpi_p = strings.Replace(strpi_p, " ", "", -1)

	lambda_b, _ := strconv.ParseFloat(strlambda_b, 64)
	pi_e, _ := strconv.ParseFloat(strpi_e, 64)
	pi_q, _ := strconv.ParseFloat(strpi_q, 64)
	pi_t, _ := strconv.ParseFloat(strpi_t, 64)
	pi_s, _ := strconv.ParseFloat(strpi_s, 64)
	pi_ch, _ := strconv.ParseFloat(strpi_ch, 64)
	pi_p, _ := strconv.ParseFloat(strpi_p, 64)

	reqValue := lambda_b * pi_e * pi_q * pi_t * pi_s * pi_ch * pi_p
	strreqValue := strconv.FormatFloat(reqValue, 'f', -1, 64)
	return strreqValue, nil
}

func (b *EBOMFrPart) FrCalcIc() (string, error) {
	partType := b.FrType

	strc1, err := b.GetFactorC1Imported()
	if err != nil {
		log.Errorf("Error: %v", err)
	}
	strc2, _ := b.GetFactorC2Imported()
	strpi_e, _ := b.GetFactorEnvImported()
	strpi_q, _ := b.GetFactorQualityImported()
	strpi_t, _ := b.GetFactorTemperatureImported()
	strpi_p := "1.0"
	strpi_a := "1.0"
	if partType == "GaAsMMIC" {
		strpi_p, _ = b.GetFactorProcessImported()
		strpi_a, _ = b.GetFactorApplicationImported()
	}

	strc1 = strings.Replace(strc1, " ", "", -1)
	strc2 = strings.Replace(strc2, " ", "", -1)
	strpi_e = strings.Replace(strpi_e, " ", "", -1)
	strpi_q = strings.Replace(strpi_q, " ", "", -1)
	strpi_t = strings.Replace(strpi_t, " ", "", -1)
	strpi_p = strings.Replace(strpi_p, " ", "", -1)
	strpi_a = strings.Replace(strpi_a, " ", "", -1)

	c1, _ := strconv.ParseFloat(strc1, 64)
	c2, _ := strconv.ParseFloat(strc2, 64)
	pi_e, _ := strconv.ParseFloat(strpi_e, 64)
	pi_q, _ := strconv.ParseFloat(strpi_q, 64)
	pi_t, _ := strconv.ParseFloat(strpi_t, 64)
	pi_p, _ := strconv.ParseFloat(strpi_p, 64)
	pi_a, _ := strconv.ParseFloat(strpi_a, 64)

	log.Info(c1, pi_e, pi_q, pi_t, pi_p, pi_a)

	reqValue := (c1*pi_t*pi_p*pi_a + c2*pi_e) * pi_q
	strreqValue := strconv.FormatFloat(reqValue, 'f', -1, 64)
	return strreqValue, nil
}
