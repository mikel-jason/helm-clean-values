<div align="center">
  <img src="./docs/assets/logo.svg" width="200px" height="200px" />

  <h1 style="border-bottom:none!important">helm clean-values</h1>

  <p>To all <i>"the YAML is self-explanatory"</i> people</p>

  <p>
    <a href="https://github.com/sarcaustech/helm-clean-values/actions/workflows/test-latest.yaml">
      <img src="https://github.com/sarcaustech/helm-clean-values/actions/workflows/test-latest.yaml/badge.svg" />
    </a>
    <img src="https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-cyan" />
    <a href="https://github.com/sarcaustech/helm-clean-values/releases/latest">
      <img src="https://img.shields.io/github/v/release/sarcaustech/helm-clean-values?color=lightgrey&include_prereleases&logo=github">
    <a href="https://choosealicense.com/licenses/gpl-3.0/#">
        <img src="https://img.shields.io/github/license/sarcaustech/helm-clean-values?color=blue&logo=github" />
    </a>
  </p>

  <hr />
</div>

## Install & update via helm plugin system

```sh
helm plugin install https://github.com/sarcaustech/helm-clean-values
```

```sh
helm plugin update clean-values
```

## Usage

Initialize chart repository first if you want to work with a remote repository
```sh
helm repo add bitnami https://charts.bitnami.com/bitnami
```

Check by comparing user values with the chart's default values
```sh
helm clean-values simple --chart bitnami/nginx -f ./testdata/bitnami.nginx.yaml
```

Check by mutating every user values primitive and test for changes in the templating result
```sh
helm clean-values mutate --chart bitnami/nginx --stdin < <(cat ./testdata/bitnami.nginx.yaml)
```
