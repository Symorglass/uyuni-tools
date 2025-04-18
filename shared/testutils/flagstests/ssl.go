// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flagstests

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// SSLGenerationFlagsTestArgs is the slice of command parameters to use with AssertSSLGenerationFlags.
var SSLGenerationFlagsTestArgs = []string{
	"--ssl-cname", "cname1",
	"--ssl-cname", "cname2",
	"--ssl-country", "OS",
	"--ssl-state", "sslstate",
	"--ssl-city", "sslcity",
	"--ssl-org", "sslorg",
	"--ssl-ou", "sslou",
}

// AssertSSLGenerationFlag checks that all the SSL certificate generation flags are parsed correctly.
func AssertSSLGenerationFlag(t *testing.T, flags *types.SSLCertGenerationFlags) {
	testutils.AssertEquals(t, "Error parsing --ssl-cname", []string{"cname1", "cname2"}, flags.Cnames)
	testutils.AssertEquals(t, "Error parsing --ssl-country", "OS", flags.Country)
	testutils.AssertEquals(t, "Error parsing --ssl-state", "sslstate", flags.State)
	testutils.AssertEquals(t, "Error parsing --ssl-city", "sslcity", flags.City)
	testutils.AssertEquals(t, "Error parsing --ssl-org", "sslorg", flags.Org)
	testutils.AssertEquals(t, "Error parsing --ssl-ou", "sslou", flags.OU)
}
