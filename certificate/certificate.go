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
	"bytes"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"sort"

	"github.com/arduino/go-paths-helper"
)

// PemToCrt ESP32 x509 certificate bundle generation utility
//
// Converts PEM and DER certificates to a custom bundle format which stores just the
// subject name and public key to reduce space
//
// The bundle will have the format: number of certificates; crt 1 subject name length; crt 1 public key length;
// crt 1 subject name; crt 1 public key; crt 2...
//
// Originally taken from:
// https://github.com/espressif/esp-idf/blob/d2471b11e78fb0af612dfa045255ac7fe497bea8/components/mbedtls/esp_crt_bundle/gen_crt_bundle.py
func PemToCrt(p *paths.Path) ([]byte, error) {
	if p == nil || !p.Exist() || p.IsDir() {
		return nil, fmt.Errorf("invalid pem path")
	}

	f, err := p.ReadFile()
	if err != nil {
		return nil, err
	}

	// pem.Decode only works for pem textual
	der := []byte{}
	for {
		block, next := pem.Decode(f)
		if block == nil && len(next) > 0 {
			return nil, errors.New("invalid pem content")
		}
		if block != nil {
			der = append(der, block.Bytes...)
			f = next
		}
		if len(next) == 0 {
			break
		}
	}

	cer, err := x509.ParseCertificates(der)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(cer, func(i, j int) bool {
		return bytes.Compare(cer[i].RawSubject, cer[j].RawSubject) == -1
	})

	bundle := []byte{}
	bundle = binary.BigEndian.AppendUint16(bundle, uint16(len(cer)))
	for _, crt := range cer {
		nameLen := len(crt.RawSubject)
		keyLen := len(crt.RawSubjectPublicKeyInfo)

		bundle = binary.BigEndian.AppendUint16(bundle, uint16(nameLen))
		bundle = binary.BigEndian.AppendUint16(bundle, uint16(keyLen))
		bundle = append(bundle, crt.RawSubject...)
		bundle = append(bundle, crt.RawSubjectPublicKeyInfo...)
	}

	return bundle, nil
}
