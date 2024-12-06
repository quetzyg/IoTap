# IoTap: IoT Device Management CLI

## Overview

A command line interface (CLI) tool designed to simplify the management and configuration of IoT devices across your network. Currently supporting Shelly relays (Generation 1 and Generation 2), IoTap provides a streamlined approach to performing bulk operations on multiple IoT devices.

## Features

- Network-wide device scanning
- Multi-device management commands
- Support for Shelly relays (Gen1 and Gen2)
- Flexible filtering by device driver
- Multiple output formats
- Easy configuration and script deployment

## Prerequisites
- Go (version 1.23 or later)
- Network access to supported IoT devices

## Installation

### From Source
```bash
git clone https://github.com/Stowify/IoTap.git
cd IoTap
go build
sudo mv IoTap /usr/local/bin/iotap
```

### Binary Release
Download the latest release from the [Releases](https://github.com/Stowify/IoTap/releases) page.

## Usage

### Basic Syntax
```bash
iotap <CIDR> <command> [flags]
```

### Commands

#### Dump
Scan and display device results
```bash
# Dump device results to screen in tabular form
iotap 192.168.1.0/24 dump

# Dump to CSV file
iotap 192.168.1.0/24 dump -f devices.csv

# Dump to JSON file
iotap 192.168.1.0/24 dump -f devices.json -format json

# Filter by driver and sort by IP
iotap 192.168.1.0/24 dump -driver shelly_gen2 -sort ip
```

#### Config
Apply configuration to multiple devices
```bash
# Apply configuration to all Shelly Gen1 devices
iotap 192.168.1.0/24 config -driver shelly_gen1 -f config.json
```

#### Version
Check device firmware versions
```bash
# Check versions for all devices
iotap 192.168.1.0/24 version

# Check versions for specific driver
iotap 192.168.1.0/24 version -driver shelly_gen2
```

#### Update
Update device firmware
```bash
# Update all devices
iotap 192.168.1.0/24 update

# Update specific driver devices
iotap 192.168.1.0/24 update -driver shelly_gen1
```

#### Deploy
Deploy scripts to devices
```bash
# Deploy a script to all Shelly Gen2 devices
iotap 192.168.1.0/24 deploy -driver shelly_gen2 -f script.js

# Deploy multiple scripts to all Shelly Gen2 devices
iotap 192.168.1.0/24 deploy -driver shelly_gen2 -f script1.js -f script2.js -f script3.js
```

#### Reboot
Reboot devices
```bash
# Reboot all devices
iotap 192.168.1.0/24 reboot

# Reboot specific driver devices
iotap 192.168.1.0/24 reboot -driver shelly_gen1
```

## Flags

### Dump Command
- `-driver`: Filter by device driver (default: all)
- `-f`: Output scan results to a file
- `-format`: Output format (default: csv)
- `-sort`: Sort devices by field (default: name)

### Other Commands
All commands support the `-driver` flag to filter by device driver.

## Configuration Files
Currently, there are two supported JSON configuration formats. Shelly Gen1 and Gen2. Each generation has a different structure, based on how the corresponding device expects the data to be passed.

### Configuration Strategy
IoTap provides flexible device targeting through a configuration strategy mechanism. In your configuration file, you can define a `strategy` section that allows precise control over which devices receive the configuration.

#### Strategy Modes

1. **Whitelist Mode**
   - When `mode` is set to `whitelist`, only devices with MAC addresses listed in the `devices` array will receive the configuration.
   - All other discovered devices will be skipped.

2. **Blacklist Mode**
   - When `mode` is set to `blacklist`, devices with MAC addresses listed in the `devices` array will be excluded from receiving the configuration.
   - All other discovered devices will receive the configuration.


##### Shelly Gen1 Config Example

In this scenario, only devices with MAC addresses `AA:BB:CC:DD:EE:FF` and `11:22:33:44:55:66` will receive the configuration.

Example:
```json
{
  "meta": {
    "device":"Shelly 1"
  },
  "strategy": {
    "mode": "whitelist",
    "devices": [
      "AA:BB:CC:DD:EE:FF",
      "11:22:33:44:55:66"
    ]
  },
  "settings": {
    "ap_roaming_enabled": true,
    "ap_roaming_threshold": -70,
    "mqtt_enable": true,
    "mqtt_server": "192.168.1.254:1883",
    "mqtt_clean_session": true,
    "mqtt_retain": false,
    "mqtt_user": "mosquitto",
    "mqtt_pass": "P@ssw0rd",
    "mqtt_reconnect_timeout_max": 60,
    "mqtt_reconnect_timeout_min": 2,
    "mqtt_keep_alive": 60,
    "mqtt_update_period": 0,
    "mqtt_max_qos": 2,
    "coiot_enable": false,
    "sntp_server": "time.cloudflare.com",
    "discoverable": false,
    "tzautodetect": true,
    "led_status_disable": false,
    "debug_enable": false,
    "allow_cross_origin": false,
    "wifirecovery_reboot_enabled": true
  },
  "settings_relay": [
    {
      "name": null,
      "appliance_type": "lock",
      "default_state": "off",
      "btn_type": "detached",
      "btn_reverse": true,
      "auto_on": 0,
      "auto_off": 3,
      "schedule": false,
      "schedule_rules": [
        "0800-012345-on"
      ]
    }
  ],
  "settings_sta": {
    "enabled": true,
    "ssid": "WIFI",
    "key": "P@ssw0rd",
    "ipv4_method": "dhcp"
  }
}
```

### Shelly Gen2 Config Example

In this scenario, devices with MAC addresses `AA:BB:CC:DD:EE:FF` and `11:22:33:44:55:66` will be skipped, and all other discovered devices will receive the configuration.

Example:
```json
{
  "meta": {
    "device": "Shelly Plus 1"
  },
  "strategy": {
    "mode": "blacklist",
    "devices": [
      "AA:BB:CC:DD:EE:FF",
      "11:22:33:44:55:66"
    ]
  },
  "sys": {
    "config": {
      "device": {
        "eco_mode": false,
        "discoverable": false
      },
      "sntp": {
        "server": "time.cloudflare.com"
      },
      "debug": {
        "mqtt": {
          "enable": false
        },
        "websocket": {
          "enable": false
        },
        "udp": {
          "addr": null
        }
      }
    }
  },
  "input": [
    {
      "id": 0,
      "config": {
        "name": null,
        "type": "switch",
        "invert": true
      }
    }
  ],
  "switch": [
    {
      "id": 0,
      "config": {
        "name": null,
        "in_mode": "detached",
        "initial_state": "off",
        "auto_on": false,
        "auto_off": true,
        "auto_off_delay": 3
      }
    }
  ],
  "wifi": {
    "config": {
      "ap": {
        "enable": false,
        "range_extender": {
          "enable": false
        }
      },
      "sta": {
        "enable": true,
        "ssid": "WIFI",
        "pass": "P@ssw0rd",
        "is_open": false,
        "ipv4mode": "dhcp"
      },
      "sta1": {
        "enable": false
      }
    }
  },
  "ble": {
    "config": {
      "enable": false
    }
  },
  "cloud": {
    "config": {
      "enable": false
    }
  },
  "mqtt": {
    "config": {
      "enable": true,
      "server": "192.168.1.254:1883",
      "user": "mosquitto",
      "pass": "P@ssw0rd",
      "ssl_ca": null,
      "enable_rpc": true,
      "rpc_ntf": false,
      "status_ntf": true,
      "enable_control": true
    }
  }
}
```

### Script File Format
[Shelly Script Language](https://shelly-api-docs.shelly.cloud/gen2/Scripts/ShellyScriptLanguageFeatures) compatible with Shelly Gen2 devices.

## Successfully tested devices:
- **Shelly**
    - **Gen1**
        - **Shelly 1** ([Product](https://www.shelly.com/en-pt/products/product-overview/shelly-1), [KB](https://kb.shelly.cloud/knowledge-base/shelly-1), [API](https://shelly-api-docs.shelly.cloud/gen1/#shelly1-shelly1pm))
    - **Gen2**
        - **Shelly Plus 1** ([Product](https://www.shelly.com/en-pt/products/product-overview/shelly-plus-1), [KB](https://kb.shelly.cloud/knowledge-base/shelly-plus-1), [API](https://shelly-api-docs.shelly.cloud/gen2/Devices/Gen2/ShellyPlus1))
        - **Shelly Plus 1 (Mini)** ([Product](https://www.shelly.com/en-pt/products/product-overview/shelly-plus-1-mini), [KB](https://kb.shelly.cloud/knowledge-base/shelly-plus-1-mini), [API](https://shelly-api-docs.shelly.cloud/gen2/Devices/Gen2/ShellyPlus1))
        - **Shelly Pro 1** ([Product](https://www.shelly.com/en-pt/products/product-overview/shelly-pro-1), [KB](https://kb.shelly.cloud/knowledge-base/shelly-pro-1-v1), [API](https://shelly-api-docs.shelly.cloud/gen2/Devices/Gen2/ShellyPro1))

## License
[Apache License 2.0](LICENSE.md)

## Disclaimer
Use responsibly. Always ensure you have proper authorization before tapping into network devices.
