// This file is part of uno-r4-wifi-fwuploader-plugin.
//
// Copyright (c) 2023 Arduino LLC.  All right reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package certificate

import (
	"testing"

	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/require"
)

func TestCertificates(t *testing.T) {
	t.Parallel()

	certificates := []struct {
		testName string
		pemFile  *paths.Path
		crtFile  *paths.Path
	}{
		{
			testName: "Single certificate",
			pemFile:  paths.New("testdata/single_certificate.pem"),
			crtFile:  paths.New("testdata/single_certificate.crt"),
		},
		{
			testName: "Multiple same pub keys",
			pemFile:  paths.New("testdata/google.pem"),
			crtFile:  paths.New("testdata/google.crt"),
		},
		{
			testName: "Multiple different pub keys",
			pemFile:  paths.New("testdata/portenta.pem"),
			crtFile:  paths.New("testdata/portenta.crt"),
		},
	}

	for _, cert := range certificates {
		cert := cert
		t.Run(cert.testName, func(t *testing.T) {
			t.Parallel()

			p := cert.pemFile

			f, err := PemToCrt(p)
			require.NoError(t, err)

			expected, err := cert.crtFile.ReadFile()
			require.NoError(t, err)
			require.Equal(t, expected, f)
		})
	}
}

func TestCertificatesWithInvalidInputs(t *testing.T) {
	i1 := `
-----BEGIN CERTIFICATE-----
asdf
-----END CERTIFICATE-----
invalid`

	i2 := `
-----BEGIN CERTIFICATE-----
asdf
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----`

	i3 := `
-----BEGIN CERTIFICATE-----
asdf
-----END CERTIFICATE-----
-----END CERTIFICATE-----`

	for _, i := range []string{i1, i2, i3} {
		f, err := paths.MkTempFile(paths.New(t.TempDir()), "invalid-pem")
		require.NoError(t, err)

		_, err = f.WriteString(i)
		require.NoError(t, err)
		require.NoError(t, f.Close())

		result, err := PemToCrt(paths.New(f.Name()))
		require.ErrorContains(t, err, "invalid pem content")
		require.Nil(t, result)
	}
}
