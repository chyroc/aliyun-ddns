# aliyun-ddns

update aliyun dns record with your public ip(v4/v6).

## Install

```shell
go install github.com/chyroc/aliyun-ddns@latest
```

## Usage

- *Set DNS Record*

```shell
aliyun-ddns set \
  -access-key-id <access-key-id> \
  -access-key-secret <access-key-secret> \
  -domain <domain> \
  -rr <rr> \
  -ip <ip>
```

will update `rr.domain` to `ip`.

`access-key-id` and `access-key-secret` is aliyun credential, can use the environment variable `ALIYUN_ACCESS_KEY_ID` and `ALIYUN_ACCESS_KEY_SECRET` instead

- *Get DNS Record*

```shell
aliyun-ddns get \
  -access-key-id <access-key-id> \
  -access-key-secret <access-key-secret> \
  -domain <domain> \
  -rr <rr>
```

- *Auto Update*

```shell
aliyun-ddns auto-update \
  -access-key-id <access-key-id> \
  -access-key-secret <access-key-secret> \
  -domain <domain> \
  -rr <rr> \
  -ip-type ipv4/ipv6
```