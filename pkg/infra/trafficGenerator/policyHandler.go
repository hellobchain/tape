package trafficGenerator

import (
	"context"

	"github.com/open-policy-agent/opa/rego"
	"github.com/wsw365904/tape/pkg/infra/basic"
)

func CheckPolicy(input *basic.Elements, rule string) (bool, error) {
	if input.Processed {
		return false, nil
	}
	rego := rego.New(
		rego.Query("data.tape.allow"),
		rego.Module("", rule),
		rego.Input(input.Orgs),
	)
	rs, err := rego.Eval(context.Background())
	if err != nil {
		return false, err
	}
	input.Processed = rs.Allowed()
	return rs.Allowed(), nil
}
