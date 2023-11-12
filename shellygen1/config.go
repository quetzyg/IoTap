package shellygen1

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	iotune "github.com/Stowify/IoTune"
	"github.com/Stowify/IoTune/device"
)

const (
	settingsPath      = "settings"
	settingsRelayPath = "settings/relay/%d"
)

// settings for Shelly Gen1 devices.
// Read more at: https://shelly-api-docs.shelly.cloud/gen1/#settings
type settings struct {
	Reset                     *bool     `json:"reset,omitempty"`                       // Will perform a factory reset of the device
	APRoamingEnabled          *bool     `json:"ap_roaming_enabled,omitempty"`          // Enable/disable AP roaming
	APRoamingThreshold        *int8     `json:"ap_roaming_threshold,omitempty"`        // Set AP roaming threshold, dBm
	MQTTEnable                *bool     `json:"mqtt_enable,omitempty"`                 // Enable connecting to a MQTT broker
	MQTTServer                *string   `json:"mqtt_server,omitempty"`                 // MQTT broker IP address and port, ex. 10.0.0.1:1883
	MQTTCleanSession          *bool     `json:"mqtt_clean_session,omitempty"`          // MQTT clean session flag
	MQTTRetain                *bool     `json:"mqtt_retain,omitempty"`                 // MQTT retain flag
	MQTTUsername              *string   `json:"mqtt_user,omitempty"`                   // MQTT username, leave empty to disable authentication
	MQTTPassword              *string   `json:"mqtt_pass,omitempty"`                   // MQTT password
	MQTTID                    *string   `json:"mqtt_id,omitempty"`                     // MQTT ID -- by default this has the form <shellymodel>-<deviceid> e.g. shelly1-B929CC.
	MQTTReconnectTimeoutMax   *int      `json:"mqtt_reconnect_timeout_max,omitempty"`  // Maximum interval for reconnect attempts
	MQTTReconnectTimeoutMin   *int      `json:"mqtt_reconnect_timeout_min,omitempty"`  // Minimum interval for reconnect attempts
	MQTTKeepAlive             *int      `json:"mqtt_keep_alive,omitempty"`             // MQTT keep alive period in seconds
	MQTTUpdatePeriod          *int      `json:"mqtt_update_period,omitempty"`          // Periodic update in seconds, 0 to disable
	MQTTMaxQoS                *int      `json:"mqtt_max_qos,omitempty"`                // Max value of QOS for MQTT packets
	CoIotEnable               *bool     `json:"coiot_enable,omitempty"`                // Enable/disable CoIoT
	CoIotUpdatePeriod         *int      `json:"coiot_update_period,omitempty"`         // Update period of CoIoT messages 15..65535s
	CoIotPeer                 *string   `json:"coiot_peer,omitempty"`                  // Set to either mcast to switch CoIoT to multicast or to ip[:port] to switch CoIoT to unicast (port is optional, default is 5683)
	SNTPServer                *string   `json:"sntp_server,omitempty"`                 // Time-server host to be used instead of the default time.google.com. An empty value disables timekeeping and requires reboot to apply.
	Name                      *string   `json:"name,omitempty"`                        // User-configurable device name
	Discoverable              *bool     `json:"discoverable,omitempty"`                // Set whether device should be discoverable (i.e. visible)
	Timezone                  *string   `json:"timezone,omitempty"`                    // Timezone identifier
	Latitude                  *int      `json:"lat,omitempty"`                         // Degrees latitude in decimal format, South is negative
	Longitude                 *int      `json:"lng,omitempty"`                         // Degrees longitude in decimal format, -180°..180°
	TimeZoneAutodetect        *bool     `json:"tzautodetect,omitempty"`                // Set this to false if you want to set custom geolocation (lat and lng) or custom timezone.
	TimeZoneUTCOffset         *int      `json:"tz_utc_offset,omitempty"`               // UTC offset
	TimeZoneDST               *bool     `json:"tz_dst,omitempty"`                      // Daylight saving time 0/1
	TimeZoneDSTAuto           *bool     `json:"tz_dst_auto,omitempty"`                 // Auto update daylight saving time 0/1
	LEDStatusDisable          *bool     `json:"led_status_disable,omitempty"`          // For Dimmer 1/2, DW2, i3, RGBW2, Plug, Plug S, EM, 3EM, 1L, 1PM, 2.5 and Button1 Enable/Disable LED indication for network status
	DebugEnable               *bool     `json:"debug_enable,omitempty"`                // Enable/disable debug file logger
	AllowCrossOrigin          *bool     `json:"allow_cross_origin,omitempty"`          // Allow/forbid HTTP Cross-Origin Resource Sharing
	WifiRecoveryRebootEnabled *bool     `json:"wifirecovery_reboot_enabled,omitempty"` // Enable/disable WiFi-Recovery reboot. Only applicable for Shelly 1/1PM, Shelly 2, Shelly 2.5, Shelly Plug/PlugS, Shelly EM, Shelly 3EM
	Relay                     *[]*relay `json:"relay,omitempty"`
}

// relay for Shelly Gen1 devices.
// Read more at: https://shelly-api-docs.shelly.cloud/gen1/#shelly1-1pm-settings-relay-0
type relay struct {
	Reset         *bool     `json:"reset,omitempty"`          // Submitting a non-empty value will reset settings for the output to factory defaults
	Name          *string   `json:"name,omitempty"`           // Set relay name
	ApplianceType *string   `json:"appliance_type,omitempty"` // Set custom configurable appliance type
	DefaultState  *string   `json:"default_state,omitempty"`  // Accepted values: off, on, last, switch
	ButtonType    *string   `json:"btn_type,omitempty"`       // Accepted values: momentary, toggle, edge, detached, action, momentary_on_release
	ButtonReverse *bool     `json:"btn_reverse,omitempty"`    // Inverts the logical state of the input
	AutoOn        *int      `json:"auto_on,omitempty"`        // Automatic flip back timer, seconds. Will engage after turning Shelly1/1PM OFF
	AutoOff       *int      `json:"auto_off,omitempty"`       // Automatic flip back timer, seconds. Will engage after turning Shelly1/1PM ON
	Power         *int      `json:"power,omitempty"`          // Shelly1 only Set user power constant to display in meters when relay is on, see /settings/power/0
	Schedule      *bool     `json:"schedule,omitempty"`       // Enable schedule timer
	ScheduleRules *[]string `json:"schedule_rules,omitempty"` // Rules for schedule activation, e.g. 0000-0123456-on
	MaxPower      *int      `json:"max_power,omitempty"`      // Shelly1PM only Set power threshold above which an overpower condition will be triggered
}

// Config implementation for the Shelly Gen1 driver.
type Config struct {
	Settings *settings `json:"settings"`
}

// Driver name of this Config implementation.
func (c *Config) Driver() string {
	return Driver
}

// MakeRequests returns a Shelly Gen1 specific HTTP request collection.
func (c *Config) MakeRequests(dev device.Resource) ([]*http.Request, error) {
	if dev.Driver() != c.Driver() {
		return nil, fmt.Errorf("device mismatch, expected %s, got %s", c.Driver(), dev.Driver())
	}

	var requests []*http.Request

	if c.Settings != nil {
		r, err := makeRequest(c.Settings, dev, settingsPath)
		if err != nil {
			return nil, err
		}
		requests = append(requests, r)

		if c.Settings.Relay != nil {
			for i, rel := range *c.Settings.Relay {
				r, err = makeRequest(rel, dev, fmt.Sprintf(settingsRelayPath, i))
				if err != nil {
					return nil, err
				}
				requests = append(requests, r)
			}
		}
	}

	return requests, nil
}

// Empty checks if the struct holding the configuration has a zero value.
func (c *Config) Empty() bool {
	return *c == Config{}
}

// makeRequest for a Shelly endpoint.
func makeRequest(i any, dev device.Resource, path string) (*http.Request, error) {
	values := structToValues(i)

	r, err := http.NewRequest(http.MethodPost, buildURL(dev.IP(), path), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	r.Header.Set(iotune.ContentTypeHeader, iotune.URLEncodedFormMimeType)

	return r, nil
}

// structToValues returns the url.Values representation of a struct. Unfortunately, the
// Shelly Gen1 API doesn't support JSON requests, only HTTP GET with a query-string in
// the URL or HTTP POST with an application/x-www-form-urlencoded POST payload.
// Read more at: https://shelly-api-docs.shelly.cloud/gen1/#common-http-api
func structToValues(cfg any) url.Values {
	cfgVal := reflect.Indirect(reflect.ValueOf(cfg))
	cfgTyp := cfgVal.Type()

	var values = url.Values{}

	for i := 0; i < cfgVal.NumField(); i++ {
		fieldValue := cfgVal.Field(i)

		// Ignore nil pointers
		if cfgTyp.Field(i).Type.Kind() == reflect.Ptr && !fieldValue.Elem().CanAddr() {
			continue
		}

		value := fmt.Sprint(reflect.Indirect(fieldValue).Interface())

		// Convert the Schedule Rules array to string (i.e. CSV), since that's what
		// the Shelly API expects. Otherwise, we'll get "Bad schedule rules!" errors
		// when passing URL encoded arrays.
		key := strings.TrimSuffix(cfgTyp.Field(i).Tag.Get("json"), ",omitempty")

		if key == "schedule_rules" {
			value = strings.Trim(value, "[]")
			value = strings.ReplaceAll(value, " ", ",")
		}

		values.Add(key, value)
	}

	return values
}
