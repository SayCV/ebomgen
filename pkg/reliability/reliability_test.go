// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reliability

import ( //"flag"
	//"fmt"
	//"strconv"
	//"strings"
	//"regexp"
	//"strconv"
	"testing"
)

func TestCase1(t *testing.T) {

	t.Logf("%d", RpmPartsCount)
	t.Logf("%d", RpmPartStress)
	t.Logf("%d", RpmPseudoStress)
	t.Logf("%d", RpmUndef)

	t.Log(len(ClsEnv))
	t.Log(ClsEnv)

	t.Log(len(FactorEnvImported))
	t.Log(FactorEnvImported)

	t.Log(len(FactorQualityImported))
	t.Log(FactorQualityImported)

	t.Log(len(FactorTemperatureImported))
	t.Log(FactorTemperatureImported)
}
