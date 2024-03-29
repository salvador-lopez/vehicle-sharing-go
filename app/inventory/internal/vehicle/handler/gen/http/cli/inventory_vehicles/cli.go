// Code generated by goa v3.11.1, DO NOT EDIT.
//
// inventory/vehicles HTTP client CLI support package
//
// Command:
// $ goa gen
// vehicle-sharing-go/internal/inventory/vehicle/infrastructure/controller/design
// -o internal/inventory/vehicle/infrastructure/controller

package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"

	client2 "vehicle-sharing-go/app/inventory/internal/vehicle/handler/gen/http/car/client"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//	command (subcommand1|subcommand2|...)
func UsageCommands() string {
	return `car (create|get)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` car create --body '{
      "color": "Nihil aperiam tempore dolor aut recusandae.",
      "id": "c2dfea5c-9040-11ee-ae47-9a5c7f2ee299",
      "vin": "Optio repellat fuga quas."
   }'` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(
	scheme, host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restore bool,
) (goa.Endpoint, interface{}, error) {
	var (
		carFlags = flag.NewFlagSet("car", flag.ContinueOnError)

		carCreateFlags    = flag.NewFlagSet("create", flag.ExitOnError)
		carCreateBodyFlag = carCreateFlags.String("body", "REQUIRED", "")

		carGetFlags  = flag.NewFlagSet("get", flag.ExitOnError)
		carGetIDFlag = carGetFlags.String("id", "REQUIRED", "Car id in uuid format")
	)
	carFlags.Usage = carUsage
	carCreateFlags.Usage = carCreateUsage
	carGetFlags.Usage = carGetUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if flag.NArg() < 2 { // two non flag args are required: SERVICE and ENDPOINT (aka COMMAND)
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = flag.Arg(0)
		switch svcn {
		case "car":
			svcf = carFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(flag.Args()[1:]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = svcf.Arg(0)
		switch svcn {
		case "car":
			switch epn {
			case "create":
				epf = carCreateFlags

			case "get":
				epf = carGetFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if svcf.NArg() > 1 {
		if err := epf.Parse(svcf.Args()[1:]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     interface{}
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "car":
			c := client2.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "create":
				endpoint = c.Create()
				data, err = client2.BuildCreatePayload(*carCreateBodyFlag)
			case "get":
				endpoint = c.Get()
				data, err = client2.BuildGetPayload(*carGetIDFlag)
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// carUsage displays the usage of the car command and its subcommands.
func carUsage() {
	fmt.Fprintf(os.Stderr, `The car service performs operations on car vehicles inventory
Usage:
    %[1]s [globalflags] car COMMAND [flags]

COMMAND:
    create: Create implements create.
    get: Get implements get.

Additional help:
    %[1]s car COMMAND --help
`, os.Args[0])
}
func carCreateUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] car create -body JSON

Create implements create.
    -body JSON: 

Example:
    %[1]s car create --body '{
      "color": "Nihil aperiam tempore dolor aut recusandae.",
      "id": "c2dfea5c-9040-11ee-ae47-9a5c7f2ee299",
      "vin": "Optio repellat fuga quas."
   }'
`, os.Args[0])
}

func carGetUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] car get -id STRING

Get implements get.
    -id STRING: Car id in uuid format

Example:
    %[1]s car get --id "c2dffd62-9040-11ee-ae47-9a5c7f2ee299"
`, os.Args[0])
}
