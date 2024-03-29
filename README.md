<div align="center">
   <h3>
      <code>oui</code>
   </h3>
   <br/>
   <b>MAC Address CLI Toolkit</b>
   <br/>
   <br/>
    <a href="https://github.com/thatmattlove/oui/actions/workflows/test.yml">
        <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/thatmattlove/oui/test.yml?style=for-the-badge">
    </a>
    <a href="https://github.com/thatmattlove/oui/releases">
        <img alt="GitHub release (latest SemVer)" src="https://img.shields.io/github/v/release/thatmattlove/oui?label=version&style=for-the-badge">
    </a>
   <br/>
   <a href="https://oui.is" target="_blank">
      <img alt="Check Out the Web Version" src="https://img.shields.io/badge/Check_Out_the_Web_Version-10b981?style=for-the-badge">
   </a>
</div>


## Installation

### macOS

#### Homebrew

```bash
brew tap thatmattlove/oui
brew install oui
```

#### MacPorts

```bash
sudo port install oui
```

### Linux

#### Debian/Ubuntu (APT)

```bash
echo "deb [trusted=yes] https://repo.fury.io/thatmattlove/ /" > /etc/apt/sources.list.d/thatmattlove.fury.list
sudo apt update
sudo apt install oui
```

#### RHEL/CentOS (YUM)

```bash
echo -e "[fury-thatmattlove]\nname=thatmattlove\nbaseurl=https://repo.fury.io/thatmattlove/\nenabled=1\ngpgcheck=0" > /etc/yum.repos.d/thatmattlove.fury.repo
sudo yum update
sudo yum install oui
```

### Windows

Coming Soon

## Usage

```console
$ oui --help
NAME:
   oui - MAC Address CLI Toolkit

USAGE:
   oui [global options] command [command options] [arguments...]

VERSION:
   2.0.4


COMMANDS:
   update, u, up      Refresh the MAC address database
   convert, c, con    Convert a MAC Address to other formats
   entires, e, count  Show the number of MAC addresses in the database
   help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug        Enable debugging (default: false)
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

### OUI Lookup

```console
$ oui F4:BD:9E:01:23:45

 F4:BD:9E:01:23:45 Results

╭──────────────────────┬────────────────────┬─────────────────────────────────────┬──────────╮
│ Prefix               │ Organization       │ Range                               │ Registry │
├──────────────────────┼────────────────────┼─────────────────────────────────────┼──────────┤
│ f4:bd:9e:00:00:00/24 │ Cisco Systems, Inc │ f4:bd:9e:00:00:00-f4:bd:9e:ff:ff:ff │ MA-L     │
╰──────────────────────┴────────────────────┴─────────────────────────────────────┴──────────╯

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
██████████▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ 20.0% Populating database...finished parsing vendors from MA-L registry
██████████████████▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ 37.0% Populating database...finished parsing vendors from CID registry
███████████████████████████▌░░░░░░░░░░░░░░░░░░░░░░ 54.0% Populating database...finished parsing vendors from IAB registry
███████████████████████████████████▌░░░░░░░░░░░░░░ 71.0% Populating database...finished parsing vendors from MA-M registry
████████████████████████████████████████████▌░░░░░ 88.0% Populating database...finished parsing vendors from MA-S registry
██████████████████████████████████████████████████ 100.0% Completed
Updated MAC Address database (2.0.4) with 49,949 records in 5 seconds
```

![GitHub](https://img.shields.io/github/license/thatmattlove/oui?style=for-the-badge&color=000000)