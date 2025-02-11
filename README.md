# IoTap: Effortless IoT Management
[![Go Report Card](https://goreportcard.com/badge/github.com/quetzyg/IoTap)](https://goreportcard.com/report/github.com/quetzyg/IoTap)
[![CI](https://github.com/quetzyg/IoTap/actions/workflows/ci.yml/badge.svg)](https://github.com/quetzyg/IoTap/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/quetzyg/IoTap/graph/badge.svg?token=ABiTPw2nI5)](https://codecov.io/gh/quetzyg/IoTap)

## Overview

IoTap is a command-line interface tool for tapping into IoT devices, designed to simplify their management and configuration on a large scale.

Born out of necessity, this tool fills a critical gap, offering a solution where none previously existed—particularly for managing extensive IoT deployments.

It is especially effective for performing batch operations and ensuring consistent, reproducible, and idempotent device management.

## Features

- Perform network-wide device scanning and discovery.
- Run commands across multiple devices simultaneously.
- Export detailed device information as CSV or JSON to a file or on-screen.
- Apply configurations to multiple devices.
- Activate/deactivate device authentication mechanisms.
- Identify devices running outdated software versions.
- Update firmware on outdated devices.
- Perform remote device restarts.
- Easy script deployment across compatible devices.

## Prerequisites
- Go (version 1.24 or later)
- Network access to supported IoT devices

## Installation

You can install the tool using one of the following methods.

### Binary Release
Download precompiled binaries from the [releases](https://github.com/quetzyg/IoTap/releases) page.

Binaries are available for:
- Linux (AMD64)
- Windows (AMD64)
- macOS (Apple Silicon)

After extracting the contents from the release, you can verify the binary checksum, like so:
```bash
sha256sum -c iotap.sha256 # Linux and macOS release
sha256sum -c iotap.exe.sha256 # Windows release
```

### Using `go install`
Ensure you have Go installed and your `$GOPATH/bin` is added to your system's `$PATH`. Then, execute:

```bash
go install github.com/quetzyg/IoTap/cmd/iotap@latest
```

You should now have the `iotap` binary in your `$PATH`.

### From Source
Clone the repository and compile the code manually to build the tool from the latest source.

```bash
git clone https://github.com/quetzyg/IoTap.git
cd IoTap
go build cmd/iotap/main.go -o iotap
sudo mv iotap /usr/local/bin/
```

## Usage

### Basic Syntax
```bash
iotap <IP|CIDR> <command> [flags]
```

#### Explanation:
- `<IP|CIDR>`: A single IP address (e.g. 192.168.1.1) or a network range in CIDR notation (e.g. 192.168.1.0/24).
- `<command>`: The command to be executed for each resolved device IP.
- `[flags]`: Optional parameters to customise the command execution.

### Available Commands

<details>
<summary><strong>dump</strong>: Output device scan results to <strong>STDOUT</strong> or to a file</summary>

```bash
# Dump device results to screen in CSV format (default)
iotap 192.168.1.0/24 dump

# Dump device results in CSV format to a file
iotap 192.168.1.0/24 dump -o devices.csv

# Dump device results in JSON format to a file
iotap 192.168.1.0/24 dump -o devices.json -f json

# Dump Shelly Gen2 only device results to the screen, sorted by IP
iotap 192.168.1.0/24 dump -d shellygen2 -s ip
```

Dump command help:
```bash
iotap 192.168.1.0/24 dump -h
```

Output:
```bash
Usage of dump:
 ./iotap <IP|CIDR> dump [flags]

Flags:
  -d value
        Device driver (default all)
  -f value
        Dump format (default csv)
  -o string
        Scan results output file
  -s value
        Sort devices by field (default name)
  -t duration
        Device probe timeout (default 2s)
```
</details>

<details>
<summary><strong>config</strong>: Apply a configuration to multiple devices</summary>

```bash
# Apply the configuration from `config.json` to all Shelly Gen1 devices
iotap 192.168.1.0/24 config -d shellygen1 -c config.json
```

Configuration command help:
```bash
iotap 192.168.1.0/24 config -h
```

Output:
```bash
Usage of config:
 ./iotap <IP|CIDR> config [flags]

Flags:
  -c string
        Device configuration file
  -d value
        Device driver (default all)
  -t duration
        Device probe timeout (default 2s)
```
</details>

<details>
<summary><strong>secure</strong>: Enable/disable device authentication</summary>

```bash
# Disable the authentication on all Shelly Gen1 devices
iotap 192.168.1.0/24 secure -d shellygen1 --off
```

```bash
# Enable the authentication on all devices
iotap 192.168.1.0/24 secure -c authentication.json
```

Secure command help:
```bash
iotap 192.168.1.0/24 config -h
```

Output:
```bash
Usage of secure:
 ./iotap <IP|CIDR> secure [flags]

Flags:
  -c string
        Auth configuration file (incompatible with --off)
  -d value
        Device driver (default all)
  -off
        Turn device authentication off (incompatible with -c)
  -t duration
        Device probe timeout (default 2s)
```
</details>

<details>
<summary><strong>version</strong>: Scan and check device versions</summary>
Identify device versions across the network, listing any that are out of date.

```bash
# Check versions for all devices
iotap 192.168.1.0/24 version

# Check versions of specific devices (Shelly Gen2)
iotap 192.168.1.0/24 version -d shellygen2
```

Version command help:
```bash
iotap 192.168.1.0/24 version -h
```

Output:
```bash
Usage of version:
 ./iotap <IP|CIDR> version [flags]

Flags:
  -d value
        Device driver (default all)
  -t duration
        Device probe timeout (default 2s)
```
</details>

<details>
<summary><strong>update</strong>: Update outdated devices</summary>
Update devices to the latest available vendor firmware.

```bash
# Update the firmware for all devices
iotap 192.168.1.0/24 update

# Update the firmware for specific devices (Shelly Gen1)
iotap 192.168.1.0/24 update -d shellygen1
```

Update command help:
```bash
iotap 192.168.1.0/24 update -h
```

Output:
```bash
Usage of update:
 ./iotap <IP|CIDR> update [flags]

Flags:
  -d value
        Device driver (default all)
  -t duration
        Device probe timeout (default 2s)
```
</details>

<details>
<summary><strong>deploy</strong>: Deploy scripts to devices</summary>

```bash
# Perform a deployment to Shelly Gen2 devices
iotap 192.168.1.0/24 deploy -d shellygen2 -c deployment.json
```

Deploy command help:
```bash
iotap 192.168.1.0/24 deploy -h
```

Output:
```bash
Usage of deploy:
 ./iotap <IP|CIDR> deploy [flags]

Flags:
  -c string
        Deployment configuration file
  -d value
        Device driver (default all)
  -t duration
        Device probe timeout (default 2s)
```
</details>

> [!WARNING]
> As part of the deployment task, the `deploy` command will **remove** any previously existing scripts from the device.

<details>
<summary><strong>reboot</strong>: Restart devices</summary>

```bash
# Reboot all devices
iotap 192.168.1.0/24 reboot

# Reboot devices by a specific driver
iotap 192.168.1.0/24 reboot -d shellygen1
```

Reboot command help:
```bash
iotap 192.168.1.0/24 reboot -h
```

Output:
```bash
Usage of reboot:
 ./iotap <IP|CIDR> reboot [flags]

Flags:
  -d value
        Device driver (default all)
  -t duration
        Device probe timeout (default 2s)
```
</details>

## IoTap Configuration

The only configuration that may be required is a set of credentials.

These must match the ones used when authentication was first enabled on the target devices.

The configuration can be loaded in two ways:

1. **Environment Variables:**
   - Credentials can be set using the following environment variables:
     ```bash
     export IOTAP_USERNAME=admin
     export IOTAP_PASSWORD=secret
     ```

2. **Predefined JSON File:**
   - The tool can load configuration from a JSON file located at `~/.config/iotap.json`.
   - Example file content:
     ```json
     {
         "credentials": {
             "username": "admin",
             "password": "secret"
         }
     }
     ```

## Command Configuration Files

Certain IoTap commands require a configuration file. These must be in the JSON format and are categorised as follows:

1. **Device Configuration:** Used with the `config` command to define parameters for configuring one or more IoT devices.

2. **Authentication Configuration:** Used with the `secure` command to specify the credentials that devices will use once authentication is enforced.

3. **Deployment Configuration:** Used with the `deploy` command to provide paths to scripts for deployment on supported devices.

Each configuration file allows defining a Policy, to enable the inclusion or exclusion of devices based on certain criteria (see below).

### Policies

Policies are a mechanism to selectively apply configurations to specific devices or groups of devices. By using these, users can:

- Include only (**whitelist**) devices that match:
  * Exact MAC addresses
  * RegEx patterns for device models and names

- Exclude (**blacklist**) devices that match:
  * Exact MAC addresses
  * RegEx patterns for device models and names

### Device Configuration

The device configuration file defines the parameters to be set on one or more IoT devices. This file enables users to customise the behavior and settings of the devices under management.

Currently, there are two device configuration types:
- Shelly Gen1
- Shelly Gen2

Each generation follows a distinct structure tailored to how the respective devices process and expect the data.

<details>
<summary><strong>Example (Shelly Gen1)</strong></summary>

In this scenario, devices with MAC addresses `AA:BB:CC:DD:EE:FF` and `11:22:33:44:55:66` will be skipped, while the remaining discovered devices will be configured.

```json
{
  "meta": {
    "device":"Shelly 1"
  },
  "policy": {
    "mode": "blacklist",
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
<summary><strong>Example (Shelly Gen2)</strong></summary>

In this scenario, only devices with models `SNSW-001X16EU` and `SNSW-001X8EU` will be configured.

```json
{
  "meta": {
    "device": "Shelly Plus 1 & Shelly Plus 1 Mini"
  },
  "policy": {
    "mode": "whitelist",
    "models": [
      "SNSW-001X\\d+EU",
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

### Authentication Configuration

The authentication configuration file specifies the credentials that one or more devices should use to enforce security.

<details>
<summary><strong>Example</strong></summary>

In this scenario, authentication will be set to any device that supports such feature.

```json
{
  "meta": {
    "device": "All",
  },
  "credentials": {
    "username": "admin",
    "password": "secret"
  }
}
```
</details>

### Deployment Configuration

The deployment configuration file allows users to define paths to scripts that can be deployed to devices supporting script execution.

<details>
<summary><strong>Example</strong></summary>

In this scenario, scripts will only be deployed to devices where the model name is `SPSW-001XE16EU`.

```json
{
  "meta": {
    "device": "Shelly Pro 1",
  },
  "policy": {
    "mode": "whitelist",
    "models": [
      "SPSW-001XE16EU"
    ]
  },
  "scripts": [
    "announce.js",
    "detached_input_on.js"
  ]
}
```
</details>

> [!IMPORTANT]
> Ensure the scripts you deploy have valid [Shelly Script Language](https://shelly-api-docs.shelly.cloud/gen2/Scripts/ShellyScriptLanguageFeatures) code.

## Device Support
The following table outlines the devices that have been successfully tested:

| Vendor | Generation | Model | Links |
|--------|------------|-------|-------|
| Shelly | Gen1 | Shelly 1 | [KB](https://kb.shelly.cloud/knowledge-base/shelly-1), [API](https://shelly-api-docs.shelly.cloud/gen1/#shelly1-shelly1pm)|
| Shelly | Gen2 | Shelly Plus 1| [KB](https://kb.shelly.cloud/knowledge-base/shelly-plus-1), [API](https://shelly-api-docs.shelly.cloud/gen2/Devices/Gen2/ShellyPlus1)|
| Shelly | Gen2 | Shelly Plus 1 (Mini)| [KB](https://kb.shelly.cloud/knowledge-base/shelly-plus-1-mini), [API](https://shelly-api-docs.shelly.cloud/gen2/Devices/Gen2/ShellyPlus1)|
| Shelly | Gen2 | Shelly Pro 1| [KB](https://kb.shelly.cloud/knowledge-base/shelly-pro-1-v1), [API](https://shelly-api-docs.shelly.cloud/gen2/Devices/Gen2/ShellyPro1)|

### Broader Device Support

While the above devices have been successfully tested, IoTap is designed with flexibility and broad compatibility in mind, meaning that it should already support a greater number of Shelly Gen1 and Gen2 devices.

### Hardware and Vendor Collaboration

**Device Compatibility Expansion:**
- The maintainer is actively seeking to expand the device compatibility matrix
- Hardware vendors and manufacturers are invited to contribute by:
  - Providing test devices
  - Sponsoring programming costs
  - Collaborating on compatibility testing

**Current Product Support:**
- Currently, the tool only supports Shelly products
- We welcome sponsorship and collaboration from other interested IoT and smart home device manufacturers

If you are a vendor, manufacturer, or potential sponsor interested in expanding our device support, please get in touch.

## Support & Sponsorship
If you find this project helpful, consider supporting its development:

🏆 [Sponsor on GitHub](https://github.com/sponsors/quetzyg)

## Contributions and Feature Roadmap
This project is open to community input and feature suggestions. Enhancements that provide broad utility and solve common challenges for multiple users are prioritary.
While creative ideas are welcome, implementation will be based on potential impact and alignment with the project's core objectives.

If you have a feature request, please submit it as an [**Idea**](https://github.com/quetzyg/IoTap/discussions/categories/ideas) in the Discussions—the more compelling the use case and potential benefit, the higher the likelihood of adoption.

**Sponsored development**: if you or your organisation want to accelerate a feature that has wide applicability, sponsorship can help move it to the top of the priority list.

## Disclaimer
Use responsibly! Always ensure you have proper authorisation before tapping into network devices.

## Security Policy
Please review the [security policy](https://github.com/quetzyg/IoTap/security/policy) to ensure proper procedures for disclosing security vulnerabilities.

## Credits
- [Quetzy Garcia](https://github.com/quetzyg)

## License
**IoTap** is open source software licensed under the [Apache License 2.0](LICENSE.md).
