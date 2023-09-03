package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
	"github.com/sudonym-btc/zap/service/config"
)

type WalletConnect struct {
	relayHost    string
	walletPubKey string
	secret       string
	clientPubkey string
	relay        nostr.Relay
}

type PayResponse struct {
	Result_type *string `json:"result_type"`
	Err         *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Result *struct {
		Preimage string `json:"preimage"`
	} `json:"result"`
}

func Connect() (*WalletConnect, error) {
	conf, err := config.LoadConfig()
	if err == nil && conf.WalletConnect != "" {
		return Parse_and_connect(conf.WalletConnect)
	}
	return nil, fmt.Errorf("no wallet connect in config", err)
}

func Parse_and_connect(nwc string) (*WalletConnect, error) {
	ctx := context.Background()
	parsedUrl, _ := url.Parse(nwc)
	var wc *WalletConnect = &WalletConnect{}
	wc.walletPubKey = parsedUrl.Host
	wc.relayHost = parsedUrl.Query().Get("relay")
	wc.secret = parsedUrl.Query().Get("secret")
	pub, _ := nostr.GetPublicKey(wc.secret)
	wc.clientPubkey = pub

	relay, err := nostr.RelayConnect(ctx, wc.relayHost)
	if err != nil {
		return nil, err
	}
	wc.relay = *relay
	return wc, nil
}

func Pay_invoice(wc *WalletConnect, invoice string) (*PayResponse, error) {
	ctx := context.Background()

	ss, _ := nip04.ComputeSharedSecret(wc.walletPubKey, wc.secret)
	content, _ := nip04.Encrypt(`{
			"method": "pay_invoice",
			"params": {
				"invoice": "`+invoice+`"
			}
		}`, ss)

	// Create payment request event
	tag := []string{"p", wc.walletPubKey}
	tags := []nostr.Tag{tag}
	ev := nostr.Event{
		PubKey:    wc.clientPubkey,
		CreatedAt: nostr.Now(),
		Kind:      nostr.KindNWCWalletRequest,
		Tags:      tags,
		Content:   content,
	}
	ev.Sign(wc.secret)

	// Filters for response
	var filters nostr.Filters
	t := make(map[string][]string)
	t["p"] = []string{wc.clientPubkey}
	t["e"] = []string{ev.ID}
	filters = []nostr.Filter{{
		Tags:  t,
		Kinds: []int{nostr.KindNWCWalletInfo, nostr.KindNWCWalletResponse, nostr.KindNWCWalletRequest},
		Limit: 1,
	}}
	sub, _ := wc.relay.Subscribe(ctx, filters)

	wc.relay.Publish(ctx, ev)

	for ev := range sub.Events {
		// handle returned event.
		// channel will stay open until the ctx is cancelled (in this case, context timeout)
		content, err := nip04.Decrypt(ev.Content, ss)
		if err != nil {
			return nil, err
		}
		payResponse := &PayResponse{}
		err2 := json.Unmarshal([]byte(content), payResponse)
		if err2 != nil {
			return nil, err2
		}
		if payResponse.Err != nil {
			return nil, fmt.Errorf(payResponse.Err.Message)
		}
		return payResponse, nil
	}

	return nil, fmt.Errorf("no response")

}
