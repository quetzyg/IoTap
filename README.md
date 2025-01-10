# IoTap: Effortless IoT Device Orchestration

## Overview

IoTap is a command-line interface tool for tapping into IoT devices, streamlining the configuration and control of network-connected devices while ensuring reproducible execution of tasks. Currently supporting [Shelly](https://shelly.cloud) devices (both Gen1 and Gen2), IoTap enables you to manage multiple devices simultaneously through a unified interface. This makes it particularly powerful for batch operations and consistent device management across your network.

## Features

- Perform network-wide device scanning and discovery.
- Run commands across multiple devices simultaneously.
- Export detailed device information in CSV or JSON formats or view it directly on-screen.
- Apply configurations to multiple devices.
- Activate/deactivate device authentication mechanisms.
- Identify devices running outdated software versions.
- Update firmware on outdated devices.
- Perform remote device restarts.
- Easy script deployment across compatible devices.

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
<summary><strong>dump</strong>: Output device scan results to <strong>STDOUT</strong> or to a file</summary>

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
<summary><strong>config</strong>: Apply a configuration to multiple devices</summary>

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
<summary><strong>secure</strong>: Enable/disable device authentication</summary>

```bash
# Disable the authentication on all Shelly Gen1 devices
iotap 192.168.1.0/24 secure -driver shelly_gen1 --off
```

```bash
# Enable the authentication on all devices
iotap 192.168.1.0/24 secure -f authentication.json
```

Secure command help
```bash
iotap 192.168.1.0/24 config -h
```

Output:
```bash
Usage of secure:
 ./iotap <CIDR> secure [flags]

Flags:
  -driver value
        Filter by device driver (default all)
  -f string
        Auth configuration file path (incompatible with --off)
  -off
        Turn device authentication off (incompatible with -f)
```
</details>

<details>
<summary><strong>version</strong>: Scan and check device versions</summary>
Identify device versions across the network, listing any that are out of date.

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
<summary><strong>update</strong>: Update outdated devices</summary>
Update devices to the latest available vendor firmware.

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
<summary><strong>deploy</strong>: Deploy scripts to devices</summary>

```bash
# Perform a deployment to Shelly Gen2 devices
iotap 192.168.1.0/24 deploy -driver shelly_gen2 -f deployment.json
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
  -f string
        Device deployment file path
```

### Important Notes
- At the moment, **only** the `shelly_gen2` driver supports this command.

- As part of the deployment task, the `deploy` command will **remove** any previously existing scripts from the device.

- Use only [Shelly Script Language](https://shelly-api-docs.shelly.cloud/gen2/Scripts/ShellyScriptLanguageFeatures) code.

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

## Device Configuration Files
Currently, there are two device configuration types:
- Shelly Gen1
- Shelly Gen2

Each generation follows a distinct structure tailored to how the respective devices process and expect the data.

### Policy Support
IoTap offers flexible device targeting using a policy mechanism. Within a configuration file, you can include a `policy` section to precisely specify which devices should receive the configuration.

#### Policy Modes
1. **Whitelist Mode**
    - When `mode` is set to `whitelist`, configurations will only be applied to devices that meet either of these criteria:
        - Their model name matches an entry in the `models` array
        - Their MAC address matches an entry in the `devices` array

2. **Blacklist Mode**
    - When `mode` is set to `blacklist`, configurations will not be applied to devices if either:
        - Their model name matches an entry in the `models` array
        - Their MAC address matches an entry in the `devices` array

<details>
<summary><strong>Shelly Gen1</strong>: Configuration Example</summary>

In this scenario, devices with MAC addresses `AA:BB:CC:DD:EE:FF` and `11:22:33:44:55:66` will be skipped, while the remaining discovered devices will be configured.

**Example:**
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
<summary><strong>Shelly Gen2</strong>: Configuration Example</summary>

In this scenario, only devices with models `SNSW-001X16EU` and `SNSW-001X8EU` will be configured.

**Example:**
```json
{
  "meta": {
    "device": "Shelly Plus 1 & Shelly Plus 1 Mini"
  },
  "policy": {
    "mode": "whitelist",
    "models": [
      "SNSW-001X16EU",
      "SNSW-001X8EU"
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

## Auth Configuration file
The `secure` command requires a single, well-defined auth configuration file.

Similarly to the configuration file formats, this one also supports policy enforcement.

### Policy Support
IoTap enables flexible device targeting through a policy mechanism. In an auth configuration file, you can include a `policy` section to precisely define which devices the authentication should apply to.

#### Policy Modes
1. **Whitelist Mode**
    - When `mode` is set to `whitelist`, the authentication will only be set to devices that meet either of these criteria:
        - Their model name matches an entry in the `models` array
        - Their MAC address matches an entry in the `devices` array

2. **Blacklist Mode**
    - When `mode` is set to `blacklist`, the authentication will not be set to devices if either:
        - Their model name matches an entry in the `models` array
        - Their MAC address matches an entry in the `devices` array

<details>
<summary>Set Authentication Example</summary>

In this scenario, authentication will be set to any device found.

**Example:**
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

### Authentication Credentials
When authentication is enabled, `IoTap` requires credentials to make authenticated requests to secured devices. These credentials must match those used when initially enabling authentication on the devices.

#### Configuration Options
You can provide credentials to IoTap using either of these two methods:

##### 1. Environment Variables
Set the following environment variables:
```bash
export IOTAP_USERNAME=admin
export IOTAP_PASSWORD=secret
```

##### 2. Configuration File
Create a configuration file at `~/.config/iotap.json` with the following content:
```json
{
  "credentials": {
    "username": "admin",
    "password": "secret"
  }
}
```

## Deployment file
The `deploy` command requires a single, well-defined deployment file.

Similarly to the auth configuration file format, this one also supports policy enforcement.

### Policy Support
IoTap enables flexible device targeting through a policy mechanism. In a deployment file, you can include a `policy` section to precisely define which devices the deployment applies to.

#### Policy Modes
1. **Whitelist Mode**
    - When `mode` is set to `whitelist`, scripts will only be deployed to devices that meet either of these criteria:
        - Their model name matches an entry in the `models` array
        - Their MAC address matches an entry in the `devices` array

2. **Blacklist Mode**
    - When `mode` is set to `blacklist`, scripts will not be deployed to devices if either:
        - Their model name matches an entry in the `models` array
        - Their MAC address matches an entry in the `devices` array

<details>
<summary><strong>Shelly Gen2</strong>: Deployment Example</summary>

In this scenario, scripts will only be deployed to devices where the model name is `SPSW-001XE16EU`.

**Example:**
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

## Device Support
The following table outlines the devices that have been successfully tested:

| Vendor | Generation | Model | Links |
|--------|------------|-------|-------|
| Shelly | Gen1 | Shelly 1 | ([Product](https://www.shelly.com/en-pt/products/product-overview/shelly-1), [KB](https://kb.shelly.cloud/knowledge-base/shelly-1), [API](https://shelly-api-docs.shelly.cloud/gen1/#shelly1-shelly1pm)) |
| Shelly | Gen2 | Shelly Plus 1| ([Product](https://www.shelly.com/en-pt/products/product-overview/shelly-plus-1), [KB](https://kb.shelly.cloud/knowledge-base/shelly-plus-1), [API](https://shelly-api-docs.shelly.cloud/gen2/Devices/Gen2/ShellyPlus1))|
| Shelly | Gen2 | Shelly Plus 1 (Mini)| ([Product](https://www.shelly.com/en-pt/products/product-overview/shelly-plus-1-mini), [KB](https://kb.shelly.cloud/knowledge-base/shelly-plus-1-mini), [API](https://shelly-api-docs.shelly.cloud/gen2/Devices/Gen2/ShellyPlus1))|
| Shelly | Gen2 | Shelly Pro 1| ([Product](https://www.shelly.com/en-pt/products/product-overview/shelly-pro-1), [KB](https://kb.shelly.cloud/knowledge-base/shelly-pro-1-v1), [API](https://shelly-api-docs.shelly.cloud/gen2/Devices/Gen2/ShellyPro1))|

### Broader Device Support

While the above devices have been successfully tested, IoTap is designed with flexibility and broad compatibility in mind, meaning that it may already support a greater number of Shelly Gen1 and Gen2 devices.

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

**Sponsored development**: if you or your organization want to accelerate a feature that has wide applicability, sponsorship can help move it to the top of the priority list.

## Disclaimer
Use responsibly! Always ensure you have proper authorisation before tapping into network devices.

## Security Policy
Please review the [security policy](https://github.com/quetzyg/IoTap/security/policy) to ensure proper procedures for disclosing security vulnerabilities.

## Credits
- [Quetzy Garcia](https://github.com/quetzyg)

## License
**IoTap** is open source software licensed under the [Apache License 2.0](LICENSE.md).
