package assets

import _ "embed"

//go:generate ./getCertificate.sh

//go:embed cdata.zip
var CertificateZipBytes []byte
