package shellygen1

import "github.com/quetzyg/IoTap/device"

// The settings type increases flexibility when dealing with
// configuration payloads expected by the Shelly Gen1 device API.
type settings map[string]any

// Config implementation for the Shelly Gen1 driver.
type Config struct {
	Policy                 *device.Policy `json:"policy,omitempty"`
	Settings               *settings      `json:"settings,omitempty"`
	SettingsAP             *settings      `json:"settings_ap,omitempty"`
	SettingsSTA            *settings      `json:"settings_sta,omitempty"`
	SettingsSTA1           *settings      `json:"settings_sta1,omitempty"`
	SettingsLogin          *settings      `json:"settings_login,omitempty"`
	SettingsCloud          *settings      `json:"settings_cloud,omitempty"`
	SettingsActions        *[]*settings   `json:"settings_actions,omitempty"`
	SettingsRelay          *[]*settings   `json:"settings_relay,omitempty"`
	SettingsPower          *[]*settings   `json:"settings_power,omitempty"`
	SettingsExtTemperature *[]*settings   `json:"settings_ext_temperature,omitempty"`
	SettingsExtHumidity    *[]*settings   `json:"settings_ext_humidity,omitempty"`
	SettingsExtSwitch      *[]*settings   `json:"settings_ext_switch,omitempty"`
}

// Driver name of this Config implementation.
func (c *Config) Driver() string {
	return Driver
}

// Empty checks if the struct holding the configuration has a zero value.
func (c *Config) Empty() bool {
	return *c == Config{}
}
