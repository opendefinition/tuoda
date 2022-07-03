# Tuoda

## About 

Tuoda is 
* a threat hunting tool for importing various logs into ArangoDB
* written in GO
* using ArangoDB as database backend
* developed on a Debian box
* Named after the Finnish word for "bring, import, get, carry, inject, win". You probably see where this is going ... 

## Prerequisites

* OS: Developed and tested on Debian, other Linux distributions may work.
* DB: ArangoDB
* Go: Newest version available

## Usage

TBA

## Configuration

This application relies on two configuration file sets:

1. Application Configuration - The main configuration of the application itself. Connect to database etc.
2. Parser Configuration files - JSON files descripting how to parse formats.

### Application Configuration

Note: Work in progress and is subject to change

This application is configuration lives inside the following path: *~/tuoda/config.json*

Here's an example of config.json to connect to your ArangoDB instance:

```json
{
	"ArangoDB": {
		"Address": "http://localhost:8529",
		"Username": "",
		"Password": "",
		"Database": ""
	}
}

```

### Parser configuration

Note: Work in progress and is subject to change

Tuoda doesn't have product specific parsers. Instead we rely on the underlying formats. For now we only support CSV.
Here's an example of a CSV Zeek parser configuration you must refer to in CLI arguments:

```json
{
	"delimiter": "\t",              # CSV delimiter
	"parser_type": "csv",           # Which parser to use
	"column_headers": {             # Section for handling CSV column headers    
		"line_pos": 7,          # Line column headers is on - counting lines starting with 1. Set to value 0 if no headers present in CSV
		"column_names": [],     # Provide your own headers here. "line_pos" setting above will be ignored
		"SkipCols": [1]         # Skip header columns in position (list of indexes starting with 1)
	},
	"data": {                       # How data lines should be handled
		"starts_at_line": 9     # On which line data starts. Note: 1 denotes the very first line
	}
}
```
