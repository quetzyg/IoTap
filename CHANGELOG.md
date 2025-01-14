# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project uses a date-based versioning system.

## [v25.01.14] - 2025-01-14

### Changed
- Remove the underscore from the Shelly driver names [[73683d6](https://github.com/quetzyg/IoTap/commit/73683d6ec2c3bd03f1167b839ef0c7438597cb3e)], [[d0fcfca](https://github.com/quetzyg/IoTap/commit/d0fcfca6aa5618c85ff6f6ba8582eb2e4157f81c)].

## [v25.01.13] - 2025-01-13

### Fixed
- Nil pointer dereference if a configuration file isn't present [[6b1a8e2](https://github.com/quetzyg/IoTap/commit/6b1a8e2dfea15aaddefa9f64f5582c36c71ccaf4)], [[48723d1](https://github.com/quetzyg/IoTap/commit/48723d17f82e67ff7f9ce7bc3c0d98f0db4af9d1)].

## [v25.01.11] - 2025-01-11

### Added
- Support for executing commands on a single IP address [[7d94d71](https://github.com/quetzyg/IoTap/commit/7d94d71c653484f5ca2bca491caa18764df84f66)].
- Ability to sort the output of the `dump` command by Vendor and Generation [[d06d4a8](https://github.com/quetzyg/IoTap/commit/d06d4a815007d9999be7e2e1bbf279e2adcf0d82)].

## [v25.01.10] - 2025-01-10

### Added
- `secure` command for enabling/disabling device authentication [[1d0907a](https://github.com/quetzyg/IoTap/commit/1d0907ad8dc0766ec8f03ac8f07292c5961187f8)].
- Authentication support for Shelly Gen1 and Gen2 devices [[3102594](https://github.com/quetzyg/IoTap/commit/3102594c55458fe22d2008da8f9cc8dfbe2a520d)].
- Configuration values support [[08181f0](https://github.com/quetzyg/IoTap/commit/08181f0863492ee4e76ba05b0d95850b94b76569)].

### Changed
- The `dump` command no longer displays the driver, showing the device vendor and generation instead [[ab7a6dd](https://github.com/quetzyg/IoTap/commit/ab7a6dd2a0f24b269c24f8d2e74c05c5c4ad55d1)].
- In order to find secured devices, Shelly Gen1 devices are now probed via the `/shelly` API endpoint, but names can no longer be fetched, showing as `N/A` instead [[511b8ba](https://github.com/quetzyg/IoTap/commit/511b8baa889f46097275317377cff945d77d7158)].

### Fixed
- Remove a line feed from a few error messages [[eb03c99](https://github.com/quetzyg/IoTap/commit/eb03c99e4b3999ce87ff0a26c5f97abd0a54bbdb)].

## [v24.12.20] - 2024-12-20

### Added
- Display the number of affected devices when a procedure finishes [[f808b93](https://github.com/quetzyg/IoTap/commit/f808b931ceabfd02e67d7dcbc08654b78b3026d3)].

### Changed
- Update how the `deploy` command works. Deployments are now defined in files with optional policy enforcement [[274b105](https://github.com/quetzyg/IoTap/commit/274b1058636f2e5f4079f5792fbc4f89d6fba552)].

### Fixed
- Ignore policy excluded errors [[f808b93](https://github.com/quetzyg/IoTap/commit/f808b931ceabfd02e67d7dcbc08654b78b3026d3)].

## [v24.12.17] - 2024-12-17

### Added
- Initial public release of the IoTap command-line interface tool
- Open-sourced project, now freely available for community use
