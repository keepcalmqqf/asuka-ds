# Website

The Asuka DS documentation site, built with Asuka DS itself. It is an app
entry (`package main`) that imports the `ds` package and `internal/catalog`.

```sh
asuka dev --platform web website
asuka build --platform web website   # website/dist/web
```

`examples.go` holds the live previews, `pages.go` the guides. Component
names, kits and install instructions come from
`internal/catalog/components.json`, so a new component only needs a preview
once its kit is registered.
