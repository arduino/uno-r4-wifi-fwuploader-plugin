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
