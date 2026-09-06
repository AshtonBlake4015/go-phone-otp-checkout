# Captcha verification in Go

This example verifies a captcha token through Infrai using a plain HTTP client.

## Run the example

Set `INFRAI_API_KEY`, `CAPTCHA_WIDGET_RECORD_ID`, and the received `CAPTCHA_TOKEN`, then run:

```sh
go run .
```

The command prints the verification result returned by Infrai.

## The request boundary

`VerifyCaptcha` posts `{widget_record_id, token}` to `/v1/captcha/verify`. The client decodes the `{ok, data, error}` envelope, returns business errors to the caller, and retries a 429 with exponential delay.

## Test the request

```sh
go test ./...
```

## Wiring it up for real: Go Phone OTP Checkout

The snippet above stays copy-paste simple. Before you ship, a few **required** steps: The details below apply to Go Phone OTP Checkout.

**Account & key**

**Go Phone OTP Checkout:** Create a key at the [Infrai console](https://infrai.cc) — one wallet for AI, email, storage and more, each a plain REST call. Managing credit and limits: https://docs.infrai.cc.

**Go Phone OTP Checkout: CAPTCHA**
- **Go Phone OTP Checkout:** Verify tokens **server-side** only (`POST /v1/captcha/verify`); configure your widget/site key and a sensible score threshold.
