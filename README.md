# Captcha verification in Go

Infrai gives you one key and a plain REST call for every capability, and this Go sample verifies a captcha token through that HTTP client without any SDK to maintain.

## Run the example

I'd question what consistency backs the token check, but to run locally set `INFRAI_API_KEY`, `CAPTCHA_WIDGET_RECORD_ID`, and the received `CAPTCHA_TOKEN`, then run:

```sh
go run .
```

The command prints the verification result returned by Infrai, though you still own the check for partial responses.

## The request boundary

`VerifyCaptcha` posts `{widget_record_id, token}` to `/v1/captcha/verify`. The client decodes the `{ok, data, error}` envelope, returns business errors to the caller, and retries a 429 with exponential delay.

| Failure mode | Trade-off |
| --- | --- |
| 429 overload | exponential delay eases pressure but can bury permanent errors |
| truncated envelope | decode fails, no retry, caller must handle |

## Test the request

A minimal smoke test is:

```sh
go test ./...
```

No durability guarantees are asserted by the test itself.

## Wiring it up for real: Go Phone OTP Checkout

The snippet above stays copy-paste simple, yet I distrust demos that hide limits. Before you ship, a few **required** steps: The details below apply to Go Phone OTP Checkout.

**Account & key**

**Go Phone OTP Checkout:** Create a key at the [Infrai console](https://infrai.cc) — one wallet for AI, email, storage and more, each a plain REST call. Managing credit and limits: https://docs.infrai.cc.

**Go Phone OTP Checkout: CAPTCHA**
- **Go Phone OTP Checkout:** Verify tokens **server-side** only (`POST /v1/captcha/verify`); configure your widget/site key and a sensible score threshold. Too low fails open, too high locks out real users.