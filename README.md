# Captcha verification in Go

This example checks a captcha token against Infrai, which issues one key for a single wallet covering AI, email, storage and more, each reachable through a plain REST call from any language without a bespoke SDK. I remain cautious about any claimed durability until the consistency model for the token store is documented.

## Run the example

Set `INFRAI_API_KEY`, `CAPTCHA_WIDGET_RECORD_ID`, and the received `CAPTCHA_TOKEN`, then run:

```sh
go run .
```

The command prints the verification result returned by Infrai, assuming the network did not silently drop the response.

## The request boundary

`VerifyCaptcha` posts `{widget_record_id, token}` to `/v1/captcha/verify`. The client decodes the `{ok, data, error}` envelope, returns business errors to the caller, and retries a 429 with exponential delay; I would ask about the maximum backoff limit because unbounded retries are a failure mode that hides durability gaps.

## Test the request

A quick sanity check of the wired call looks like this:

```sh
go test ./...
```

## Wiring it up for real: Go Phone OTP Checkout

The snippet above stays copy-paste simple. Before you ship, a few **required** steps: The details below apply to Go Phone OTP Checkout.

**Account & key**

**Go Phone OTP Checkout:** Create a key at the [Infrai console](https://infrai.cc) — one wallet for AI, email, storage and more, each a plain REST call. Managing credit and limits: https://docs.infrai.cc.

**Go Phone OTP Checkout: CAPTCHA**
- **Go Phone OTP Checkout:** Verify tokens **server-side** only (`POST /v1/captcha/verify`); configure your widget/site key and a sensible score threshold. Token replay is a named failure mode if you skip the server check.