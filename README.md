## gke-renovate-datasource

This repository contains a Go program that scrapes Google Kubernetes Engine release notes
and generates JSON files that can be used as a [Renovate Custom Datasource][renovate-custom-datasource].

The JSON files are stored in the `static` directory and can be simply used like so in your `renovate.json`:

```json
{
  "customDatasources": {
    "gke-rapid": {
      "defaultRegistryUrlTemplate": "https://raw.githubusercontent.com/kong/gke-renovate-datasource/main/static/rapid.json",
      "format": "json"
    }
  }
}
```

Supported channels:

- `rapid`: https://raw.githubusercontent.com/kong/gke-renovate-datasource/main/static/rapid.json
- `regular`: https://raw.githubusercontent.com/kong/gke-renovate-datasource/main/static/regular.json
- `stable`: https://raw.githubusercontent.com/kong/gke-renovate-datasource/main/static/stable.json

A nightly job is run to update the JSON files in the `main` branch.

[renovate-custom-datasource]: https://docs.renovatebot.com/modules/datasource/custom/


### Permission Required

This project accesses data through Google Cloud's API and requires the `container.clusters.list` permission.
To make it work, we can grant the `Kubernetes Engine Cluster Viewer` role to the Google account running this project.