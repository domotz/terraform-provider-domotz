# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0]

### Added
- Add context propagation
- Add retry logic
- Input validation
- Add User-Agent

## [1.0.0]

### Added
- Initial release
- Support for Domotz Public API v1
- Basic provider configuration
- Core resource and data source implementations
- Resources:
  - `domotz_device` - Manage external IP devices
  - `domotz_custom_tag` - Manage custom tags
  - `domotz_device_tag_binding` - Bind tags to devices
  - `domotz_snmp_sensor` - Configure SNMP sensors
  - `domotz_tcp_sensor` - Configure TCP port sensors
- Data sources:
  - `domotz_agent` - Query agent details
  - `domotz_device` - Query device details
  - `domotz_devices` - List all devices
  - `domotz_device_variables` - Query device metrics
- API client with full CRUD operations
- Import support for all resources
- Comprehensive documentation and examples