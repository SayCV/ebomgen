package types

import (
	"io"
	//"bytes"
	"fmt"
	//"net/url"
	//"sort"
	//"strconv"
	"strings"
	//"github.com/pkg/errors"
	//"github.com/sirupsen/logrus"
	//log "github.com/sirupsen/logrus"
	//yaml "gopkg.in/yaml.v2"
)

// EBOMItem part propery
type EBOMItem struct {
	Quantity   int
	References []string
	Value      string
	Library    string
	Footprint  string
	Desc       string
	Attributes map[string]string
}

// EBOMSheet include all parts
type EBOMSheet struct {
	Headers       []string
	Items         []EBOMItem
	CustomHeaders []string
}

func (b *EBOMSheet) appendField(fieldName string) {
	for _, f := range b.Headers {
		if f == fieldName {
			return
		}
	}
	b.CustomHeaders = append(b.CustomHeaders, fieldName)
	b.Headers = append(b.Headers, fieldName)
}

func (b *EBOMSheet) generateHeaders() error {
	b.Headers = []string{"Quantity", "References", "Value", "Footprint"}
	return nil
}

func (b *EBOMSheet) makeUniqueIdentifier(comp EBOMItem) string {
	ident := fmt.Sprintf("Value=%s_Footprint=%s", comp.Value, comp.Footprint)

	return ident
}

func (b *EBOMSheet) writeItem(w io.Writer, i EBOMItem) error {
	res := make([]string, 0, len(i.Attributes)+5)
	res = append(res, fmt.Sprintf("%d", i.Quantity))
	res = append(res, fmt.Sprintf(`"%s"`, strings.Join(i.References, ",")))
	res = append(res, `"`+i.Value+`"`)
	res = append(res, `"`+i.Footprint+`"`)
	//for _, f := range i.Attributes {
	//	res = append(res, `"`+f+`"`)
	//}

	_, err := fmt.Fprintln(w, strings.Join(res, ","))

	return err
}

// WriteCSV saveas csv file
func (b *EBOMSheet) WriteCSV(w io.Writer) error {

	_, err := fmt.Fprintln(w, strings.Join(b.Headers, ","))
	if err != nil {
		return err
	}

	for _, i := range b.Items {
		err = b.writeItem(w, i)
		if err != nil {
			return err
		}
	}
	return nil
}
