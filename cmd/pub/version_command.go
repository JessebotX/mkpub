package main

import (
	"fmt"
)

const (
	Version           = "0.14.0"
	Homepage          = "https://github.com/JessebotX/pub"
	LicenseIdentifier = "BSL-1.0"
	LicenseLink       = "https://github.com/JessebotX/pub/blob/master/LICENSE.txt"
)

type VersionCommand struct{}

func (VersionCommand) Run() error {
	fmt.Printf(`pub v%s

Homepage :: <%s>
License  :: %s (<%s>)
\n`, Version, Homepage, LicenseIdentifier, LicenseLink)
	return nil
}
