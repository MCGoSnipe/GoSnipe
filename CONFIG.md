# Configuration

## Installing

Clone repo and run this command
```
go install gosnipe.go
```

## Running after install
`gosnipe`

## CLI Flags

```
-h, --help:
	Invokes this help page.
Required:
	-n, --name:
		Sets the name to snipe.
		Value: string
	One of either:
		-m, --microsoft:
			Uses a Microsoft account if true. Requires console input or --bearer set.
		-p, --path:
			Path to your accounts file.
			Value: string
			Default: "accounts.txt"
Optional:
	-b, --bearer:
		Provides the response for --microsoft.
		Value: String
		Default: ""
	-l, --speed-limit:
		Delay between requests (in milliseconds). Ex:
			Sent at .0001
			Delay of .0003
			Sent at .0004
		Value: int
		Default: 0
	-o, --offset:
		Sets the offset (in milliseconds) to value specified.
		Value: float
		Default: 0
	-a, --auto-offset
		Automatically sets the offset using X requests.
		Value: int
		Default: 3 (if flag is passed without value)
	-r, --requests:
		Sends this many requests.
		Value: int
		Default: 2
```
