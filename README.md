# IoTap: IoT Device Management CLI

## Overview

A CLI (command line interface) tool designed to simplify the management and configuration of IoT devices across your network. Currently supporting Shelly (Gen1 and Gen2) relay devices, IoTap provides a streamlined approach to performing bulk operations on multiple IoT devices.

## Features

- Network-wide device scanning
- Multi-device management commands
- Support for Shelly (Gen1 and Gen2) relay devices
- Flexible device filtering and sorting
- Multiple output formats
- Easy configuration and script deployment

## Prerequisites
- Go (version 1.23 or later)
- Network access to supported IoT devices

## Installation

### From Source
```bash
git clone https://github.com/quetzyg/IoTap.git
cd IoTap
go build
sudo mv IoTap /usr/local/bin/iotap
```

### Binary Release
Download the latest release from the [Releases](https://github.com/quetzyg/IoTap/releases) page.

## Usage

### Basic Syntax
```bash
iotap <CIDR> <command> [flags]
```

### Commands

<details>
<summary><strong>dump</strong>: Output the device scan results to <strong>STDOUT</strong> or to a file.</summary>

```bash
# Dump device results to screen in tabular form
iotap 192.168.1.0/24 dump

# Dump device results to CSV file
iotap 192.168.1.0/24 dump -f devices.csv

# Dump device results to JSON file
iotap 192.168.1.0/24 dump -f devices.json -format json

# Dump device results filtered by driver and sorted by IP
iotap 192.168.1.0/24 dump -driver shelly_gen2 -sort ip
```

Dump command help
```bash
iotap 192.168.1.0/24 dump -h
```

Output:
```bash
Usage of dump:
 ./iotap <CIDR> dump [flags]

Flags:
  -driver value
        Filter by device driver (default all)
  -f string
        Output the scan results to a file
  -format value
        Dump output format (default csv)
  -sort value
        Sort devices by field (default name)

```
</details>

<details>
<summary><strong>config</strong>: Apply a configuration to multiple devices of a specific driver</summary>

```bash
# Apply the configuration from `config.json` to all Shelly Gen1 devices
iotap 192.168.1.0/24 config -driver shelly_gen1 -f config.json
```

Configuration command help
```bash
iotap 192.168.1.0/24 config -h
```

Output:
```bash
Usage of config:
 ./iotap <CIDR> config [flags]

Flags:
  -driver value
        Filter by device driver (default all)
  -f string
        Device configuration file path
```
</details>

<details>
<summary><strong>version</strong>: Check the device firmware versions</summary>
Out of date devices will be listed.

```bash
# Check versions for all devices
iotap 192.168.1.0/24 version

# Check versions for specific driver (Shelly Gen2)
iotap 192.168.1.0/24 version -driver shelly_gen2
```

Version command help
```bash
iotap 192.168.1.0/24 version -h
```

Output:
```bash
Usage of version:
 ./iotap <CIDR> version [flags]

Flags:
  -driver value
        Filter by device driver (default all)
```
</details>

<details>
<summary><strong>update</strong>: Update devices to the latest available vendor firmware</summary>

```bash
# Update the firmware for all devices
iotap 192.168.1.0/24 update

# Update the firmware for specific devices (Shelly Gen1)
iotap 192.168.1.0/24 update -driver shelly_gen1
```

Update command help
```bash
iotap 192.168.1.0/24 update -h
```

Output:
```bash
Usage of update:
 ./iotap <CIDR> update [flags]

Flags:
  -driver value
        Filter by device driver (default all)
```
</details>

<details>
<summary><strong>deploy</strong>: Deploy one or more scripts to devices.</summary>

Note that, at the moment, only Shelly Gen2 devices support this command.

```bash
# Deploy a script to all Shelly Gen2 devices
iotap 192.168.1.0/24 deploy -driver shelly_gen2 -f script.js

# Deploy multiple scripts to all Shelly Gen2 devices
iotap 192.168.1.0/24 deploy -driver shelly_gen2 -f script1.js -f script2.js -f script3.js
```

Deploy command help
```bash
iotap 192.168.1.0/24 deploy -h
```

Output:
```bash
Usage of deploy:
 ./iotap <CIDR> deploy [flags]

Flags:
  -driver value
        Filter by device driver (default all)
  -f value
        Deploy script file path (allows multiple calls)
```
</details>

<details>
<summary><strong>reboot</strong>: Restart devices</summary>

```bash
# Reboot all devices
iotap 192.168.1.0/24 reboot

# Reboot specific devices by driver
iotap 192.168.1.0/24 reboot -driver shelly_gen1
```

Reboot command help
```bash
iotap 192.168.1.0/24 reboot -h
```

Output:
```bash
Usage of reboot:
 ./iotap <CIDR> reboot [flags]

Flags:
  -driver value
        Filter by device driver (default all)
```
</details>

## Configuration Files
Currently, there are two supported JSON configuration formats. Shelly Gen1 and Gen2. Each generation has a different structure, based on how the corresponding device expects the data to be passed.

### Configuration Policy
IoTap provides flexible device targeting through a configuration policy mechanism. In a configuration file, you can define a `policy` section that allows precise control over which devices receive the configuration.

#### Strategy Modes
1. **Whitelist Mode**
   - When `mode` is set to `whitelist`, only devices with MAC addresses listed in the `devices` array will receive the configuration.
   - All other discovered devices will be skipped.

2. **Blacklist Mode**
   - When `mode` is set to `blacklist`, devices with MAC addresses listed in the `devices` array will be excluded from receiving the configuration.
   - All other discovered devices will receive the configuration.

<details>
<summary><strong>Shelly Gen1</strong>: Configuration Example</summary>

In this scenario, only devices with MAC addresses `AA:BB:CC:DD:EE:FF` and `11:22:33:44:55:66` will receive the configuration.

**Example:**
```json
{
  "meta": {
    "device":"Shelly 1"
  },
  "policy": {
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
</details>

<details>
<summary><strong>Shelly Gen2</strong>: Configuration Example</summary>

In this scenario, devices with MAC addresses `AA:BB:CC:DD:EE:FF` and `11:22:33:44:55:66` will be skipped, and all other discovered devices will receive the configuration.

**Example:**
```json
{
  "meta": {
    "device": "Shelly Plus 1"
  },
  "policy": {
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
</details>

### Script File Format
Use [Shelly Script Language](https://shelly-api-docs.shelly.cloud/gen2/Scripts/ShellyScriptLanguageFeatures) compatible code when deploying to Shelly Gen2 devices.

## Successfully tested devices:
- **Shelly**
    - **Gen1**
        - **Shelly 1** ([Product](https://www.shelly.com/en-pt/products/product-overview/shelly-1), [KB](https://kb.shelly.cloud/knowledge-base/shelly-1), [API](https://shelly-api-docs.shelly.cloud/gen1/#shelly1-shelly1pm))
    - **Gen2**
        - **Shelly Plus 1** ([Product](https://www.shelly.com/en-pt/products/product-overview/shelly-plus-1), [KB](https://kb.shelly.cloud/knowledge-base/shelly-plus-1), [API](https://shelly-api-docs.shelly.cloud/gen2/Devices/Gen2/ShellyPlus1))
        - **Shelly Plus 1 (Mini)** ([Product](https://www.shelly.com/en-pt/products/product-overview/shelly-plus-1-mini), [KB](https://kb.shelly.cloud/knowledge-base/shelly-plus-1-mini), [API](https://shelly-api-docs.shelly.cloud/gen2/Devices/Gen2/ShellyPlus1))
        - **Shelly Pro 1** ([Product](https://www.shelly.com/en-pt/products/product-overview/shelly-pro-1), [KB](https://kb.shelly.cloud/knowledge-base/shelly-pro-1-v1), [API](https://shelly-api-docs.shelly.cloud/gen2/Devices/Gen2/ShellyPro1))

## Support & Sponsorship
If you find this project helpful, consider supporting its development:

🏆 [Sponsor on GitHub](https://github.com/sponsors/quetzyg)

## Contributions and Feature Roadmap
This project is open to community input and feature suggestions. Enhancements that provide broad utility and solve common challenges for multiple users are prioritary.
While creative ideas are welcome, implementation will be based on potential impact and alignment with the project's core objectives.

If you have a feature request, please submit it as an [**Idea**](https://github.com/quetzyg/IoTap/discussions/categories/ideas) in the Discussions—the more compelling the use case and potential benefit, the higher the likelihood of adoption.

**Sponsored development**: if you or your organization want to accelerate a feature that has wide applicability, sponsorship can help move it to the top of the priority list.

## Disclaimer
Use responsibly! Always ensure you have proper authorisation before tapping into network devices.

## Security
If you found a security related issue, please email **security (at) altek (dot) org**.

## Credits
- [Quetzy Garcia](https://github.com/quetzyg)

## License
**IoTap** is open source software licensed under the [Apache License 2.0](LICENSE.md).
