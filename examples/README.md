## requirements.
- [tile-config-generator](https://github.com/pivotalservices/tile-config-generator/releases).
- go 1.12 =>
- a running PCF Ops Manager

## generate/update tile configs.
get a legacy pivnet token [here](https://network.pivotal.io/users/dashboard/edit-profile).

```
tile-config-generator generate \
    --token=YOUR_LEGACY_TOKEN \
    --product-slug=elastic-runtime \
    --product-version=2.5.1 \
    --product-glob='srt*.pivotal' \
    --include-errands \
    --do-not-include-product-version \
    --base-directory=templates/assets/tiles/srt
```

## update pattern.yml.
We already provided a skeleton pattern.yml which work with pcf 2.5.1
review `templates/assets/pattern.yml`

## generate template.
the yml templates in the assets directory should be embedded using go generate:

`go generate templates/templates.go`

### running localy.
```
export PIVNET_TOKEN=#CHANGEME
export OPSMAN_TARGET=#CHANGEME
export OPSMAN_PASSWORD=#CHANGEME
go run -tags dev main.go
```
