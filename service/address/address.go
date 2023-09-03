package address

import (
	"math"

	lnurl "github.com/fiatjaf/go-lnurl"
	"github.com/gookit/slog"
)

func Fetch(address string) (*lnurl.LNURLPayParams, error) {
	slog.Debug("Fetching ln address", address)

	_, params, err := lnurl.HandleLNURL(address)
	if err != nil {
		slog.Warn("Failed fetching ln address info")
		return nil, err
	}
	slog.Debug("Fetched ln address info", params)
	test := params.(lnurl.LNURLPayParams)
	test2 := &test
	return test2, nil
}
func FetchFromParams(msats int64, comment string, params lnurl.LNURLPayParams) (*string, error) {
	slog.Debug("Fetching invoice from params", msats, comment, params)
	var pd lnurl.PayerDataValues = lnurl.PayerDataValues{}

	comment = comment[0:int(math.Min(float64(len(comment)), float64(params.CommentAllowed)))]

	result, err := params.Call(msats, comment, &pd)
	if err != nil {
		slog.Error("Failed fetching invoice", err)
		return nil, err
	}
	slog.Debug("Fetched invoice", result)

	return &result.PR, err
}
func FetchInvoice(msats int64, comment string, address string) (*string, error) {
	slog.Debug("Fetching invoice", msats, comment, address)
	params, _ := Fetch(address)
	return FetchFromParams(msats, comment, *params)
}
