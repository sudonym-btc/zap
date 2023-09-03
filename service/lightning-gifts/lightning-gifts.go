package lightninggifts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gookit/slog"
)

type LightningGift struct {
	Amount           int              `json:"amount"`
	LightningInvoice lightningInvoice `json:"lightningInvoice"`
	OrderId          string           `json:"orderId"`
}
type lightningInvoice struct {
	Payreq string `json:"payreq"`
}

func CreateGift(amount int, sender string, comment string) (*LightningGift, error) {
	slog.Debug("Creating gift", amount, sender, comment)
	postBody := []byte(`{
		"amount":        ` + fmt.Sprint(amount) + `,
		"senderName":    "` + sender + `"
	}`)
	// ,
	// 	"senderMessage": "` + comment + `"
	// slog.Info("Creating gift", string(postBody))
	bodyReader := bytes.NewReader(postBody)

	resp, err := http.Post("https://api.lightning.gifts/create", "application/json", bodyReader)
	if err != nil {
		slog.Warn("Failed creating gift", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		slog.Warn("Failed creating gift", err)
		return nil, err
	}
	gift := LightningGift{}
	err = json.Unmarshal(body, &gift)
	if err != nil {
		slog.Warn("Failed creating gift", body, err)
		return nil, err
	}
	if gift.OrderId == "" {
		slog.Warn("Failed creating gift: Order ID not set", body, err)
		return nil, fmt.Errorf("Order ID not set")
	}
	slog.Info("Created gift", gift)
	return &gift, nil
}

func RedeemUrl(orderId string) string {
	return "https://lightning.gifts/redeem/" + orderId
}
