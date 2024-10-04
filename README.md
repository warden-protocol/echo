# Echo

![echo](https://github.com/warden-protocol/echo/blob/main/img/echo.png?raw=true)

Echo is meant to monitor blockchain node readiness in Kubernetes clusters.
It works best to be used in a sidecar container.

## Config

| ENV              | Type   | Default                                      | Description                                                |
| ---------------- | ------ | -------------------------------------------- | ---------------------------------------------------------- |
| PORT             | string | 10010                                        | Port of the service                                        |
| ENDPOINTS        | string | http://localhost:26657,http://localhost:1317 | Endpoints to monitor                                       |
| PEERS            | string |                                              | Peers to compare node height                               |
| BEHIND_THRESHOLD | int    | 10                                           | Behind threshold, allowing configuration of block mismatch |
