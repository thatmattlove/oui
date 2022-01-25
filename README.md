<div align="center">
<h3>
    <code>oui</code>
</h3>
<br/>
MAC Address CLI Toolkit
<br/>
<br/>
    <a href="https://github.com/thatmattlove/oui/actions?query=workflow%3Atest">
        <img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/thatmattlove/oui/test?style=for-the-badge">
    </a>
    <a href="https://github.com/thatmattlove/oui/releases">
        <img alt="GitHub release (latest SemVer)" src="https://img.shields.io/github/v/release/thatmattlove/oui?label=version&style=for-the-badge">
    </a>
</div>

## Installation

### macOS

```console
$ brew tap thatmattlove/oui
$ brew install oui
```

### Linux

#### Debian/Ubuntu (APT)

```console
$ echo "deb [trusted=yes] https://repo.fury.io/thatmattlove/ /" > /etc/apt/sources.list.d/thatmattlove.fury.list
$ sudo apt update
$ sudo apt install oui
```

#### RHEL/CentOS (YUM)

```console
$ echo -e "[fury-thatmattlove]\nname=thatmattlove\nbaseurl=https://repo.fury.io/thatmattlove/\nenabled=1\ngpgcheck=0" > /etc/yum.repos.d/thatmattlove.fury.repo
$ sudo yum update
$ sudo yum install oui
```

### Windows

TODO

## Usage

```console
$ oui --help
NAME:
   oui - MAC Address CLI Toolkit

USAGE:
   oui [global options] command [command options] [arguments...]

VERSION:
   0.1.5


COMMANDS:
   update, u, up    Refresh the MAC address database
   convert, c, con  Convert a MAC Address to other formats
   help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug        Enable debugging (default: false)
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

### OUI Lookup

```console
$ oui F4:BD:9E:01:23:45

 F4:BD:9E:01:23:45 Results

╭──────────────────────┬────────────────────┬─────────────────────────────────────╮
│ Prefix               │ Organization       │ Range                               │
├──────────────────────┼────────────────────┼─────────────────────────────────────┤
│ f4:bd:9e:00:00:00/24 │ Cisco Systems, Inc │ f4:bd:9e:00:00:00-f4:bd:9e:ff:ff:ff │
╰──────────────────────┴────────────────────┴─────────────────────────────────────╯

```

### Conversion

```console
$ oui convert F4:BD:9E:01:23:45

 F4:BD:9E:01:23:45

╭─────────────┬───────────────────────╮
│ Hexadecimal │ f4:bd:9e:01:23:45     │
│ Dotted      │ f4bd.9e01.2345        │
│ Dashed      │ f4-bd-9e-01-23-45     │
│ Integer     │ 269095236870981       │
│ Bytes       │ {244,189,158,1,35,69} │
╰─────────────┴───────────────────────╯

```

### Updating the Database

```
$ oui update

Updating MAC Address Database
██████████████████████████████████████████████████ 100.0% Completed
Updated MAC Address database (v0) with 44,839 records in 14 seconds

```

![GitHub](https://img.shields.io/github/license/thatmattlove/oui?style=for-the-badge)