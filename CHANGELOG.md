# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.2] - 2025-02-05

### Added
- Added token info and chart

## [0.2.1] - 2024-12-24

### Added
- Added vote status validation

## [0.2.0] - 2024-12-04

### Added
- Implement getting top delegates/delegators by address from storage
- Implement getting delegates/delegators list by params
- Implement getting count of delegations by address from storage

## [0.1.15] - 2024-10-08

### Added
- Extend DAO object with active proposals ids

## [0.1.14] - 2024-10-02

### Added
- Added vote now

## [0.1.13] - 2024-09-22

### Added
- Added delegates total

## [0.1.12] - 2024-09-13

### Added
- Added github actions for building docker container

## [0.1.11] - 2024-09-09

### Fixed
- Update core storage protocol

## [0.1.10] - 2024-09-09

### Added
- All delegations list by address

### Added
- Added proxy api for delegates

## [0.1.9] - 2024-08-13

### Changed
- Parameter's name for search

## [0.1.8] - 2024-07-22

### Added
- Search for votes 

## [0.1.7] - 2024-07-05

### Changed
- Extend vote response with proposal identifier

## [0.1.6] - 2024-04-19

### Added
- DAO recommendations endpoint
- DAO popularity index

## [0.1.5] - 2024-04-10

### Added
- Daos user participates in

## [0.1.4] - 2024-03-22

### Added
- Stats endpoint

## [0.1.3] - 2024-03-15

### Added
- Total vp for votes

## [0.1.2] - 2024-03-13

### Added
- Get Ens Names

## [0.1.1] - 2024-03-06

### Fixed
- Fixed Dockerfile

## [0.1.0] - 2024-03-02

### Changed
- Updated all internal dependencies
- Changed the path name of the go module

### Added
- Added LICENSE information
- Added info for contributing
- Added github issues templates
- Added linter and unit-tests workflows for github actions
- Added badges with link to the license and passed workflows

## [0.0.24] - 2024-02-06

### Added
- Active votes, verified fields to dao

### [0.0.23] - 2024-02-05

### Added
- Order by voter parameter

### [0.0.22] - 2024-01-30

### Added
- User votes

### [0.0.21] - 2023-12-14

### Added
- Author ens name field for votes

### [0.0.20] - 2023-12-06

### Added
- Author ens name field for proposals

### [0.0.19] - 2023-12-04

### Added
- Added voting methods

### [0.0.18] - 2023-10-09

### Added
- Voters count field for dao info

## [0.0.17] - 2023-10-06

### Fixed
- Proposal ends soon event

## [0.0.16] - 2023-09-18

### Added
- Feed endpoint to get feed by filters

## [0.0.15] - 2023-09-12

### Changed
- Top proposals uses new field in the request

## [0.0.14] - 2023-09-12

### Changed
- Mark votes choice field as json.RawMessage due to multiple values

## [0.0.13] - 2023-08-23

### Added
- Proposal timeline field

## [0.0.12] - 2023-07-18

### Changed
- Extend vote model

## [0.0.11] - 2023-07-14

### Fixed
- Fixed action field style in the result json in timeline

## [0.0.10] - 2023-07-14

### Fixed
- Fixed missed feed item timeline logic

## [0.0.9] - 2023-07-14

### Added
- Dao activity since field

## [0.0.8] - 2023-07-14

### Fixed
- Updated core-api protocol to v0.0.10
- Fixed type mismatched fields

## [0.0.7] - 2023-07-12

### Fixed
- Fixed getting DAO list from core storage service
- Fixed missed params in DAO and Strategy models

## [0.0.6] - 2023-07-11

### Added
- Proposal top endpoint

## [0.0.5] - 2023-07-07

### Added
- Filtering proposals by title

## [0.0.4] - 2023-07-06

### Added
- Flat feed

## [0.0.3] - 2023-06-29

### Added
- Filter dao by ids

## [0.0.2] - 2023-06-14

### Removed
- Logrus library

## [0.0.1] - 2023-06-14

### Added
- Added skeleton app
- Added daos handlers
- Added proposals handlers
- Added subscriber handlers
