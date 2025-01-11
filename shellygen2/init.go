package shellygen2

import "github.com/quetzyg/IoTap/device"

// init registers this package's device implementations.
// It will be automatically called when the package is imported.
func init() {
	device.RegisterProber(Driver, func() device.Prober {
		return &Prober{}
	})

	device.RegisterConfig(Driver, func() device.Config {
		return &Config{}
	})
}
