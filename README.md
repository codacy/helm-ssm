# Helm SSM Plugin

[![CircleCI](https://circleci.com/gh/Altitude-sports/helm-ssm.svg?style=svg)](https://circleci.com/gh/Altitude-sports/helm-ssm)

This is a **helm3** plugin to help developers inject values coming from AWS SSM
parameters, on the `values.yaml` file. It also leverages the wonderful [sprig](http://masterminds.github.io/sprig/)
package, thus making all its functions available when parsing.

Since **helm2 is deprecated** the current version of the plugin only supports helm3. The last version
to support helm2 is [v2.2.1](https://github.com/codacy/helm-ssm/releases/tag/2.2.1). There will be
no further patches or updates to this legacy version.

## Usage

Loads a template file,  and writes the output.

Simply add placeholders like `{{ssm "path" "option1=value1" }}` in your
file, where you want it to be replaced by the plugin.

Currently the plugin supports the following options:

- `region=eu-west-1` - to resolve that parameter in a specific region
- `default=some-value` - to give a default **string** value when the ssm parameter is optional. The plugin will throw an error when values are not defined and do not have a default.
- `prefix=/something` - you can use this to specify a given prefix for a parameter without affecting the path. It will be concatenated with the path before resolving.

### Values file

```yaml
service:
ingress:
  enabled: false
  hosts:
    - service.{{ssm "/exists/subdomain" }}
    - service1.{{ssm "/empty/subdomain" "default=codacy.org" }}
    - service2.{{ssm "/exists/subdomain" "default=codacy.org" "region=eu-west-1" }}
    - service3.{{ssm "/subdomain" "default=codacy.org" "region=eu-west-1" "prefix=/empty" }}
    - service4.{{ssm "/securestring" }}
    - service5.{{env "ENV_VAR" }}

```

when you do not want a key to be defined, you can use a condition and an empty default value:

```yaml
service:
ingress:
  enabled: false
  hosts:
    {{- with $subdomain := (ssm "/exists/subdomain" "default=") }}{{ if $subdomain }}
    - service.{{$subdomain}}
    {{- end }}{{- end }}

```

### Command

```sh
$ helm ssm [flags]
```

### Flags

```sh
  -d, --dry-run                 does not replace the file content
  -h, --help                    help for ssm
  -p, --profile string          aws profile to fetch the ssm parameters
  -o, --target-dir string       dir to output content
  -f, --values valueFilesList   specify values in a YAML file (can specify multiple) (default [])
  -v, --verbose                 show the computed YAML values file/s
```

## Example

[![asciicast](https://asciinema.org/a/c2zut95zzbiKyk5gJov67bxsP.svg)](https://asciinema.org/a/c2zut95zzbiKyk5gJov67bxsP?t=1)

## Install

Choose the latest version from the releases and install the
appropriate version for your OS as indicated below.

```sh
$ helm plugin add https://github.com/codacy/helm-ssm
```

### Developer (From Source) Install

If you would like to handle the build yourself, instead of fetching a binary,
this is how we recommend doing it.

- Make sure you have [Go](http://golang.org) installed.

- Clone this project

- In the project directory run
```sh
$ make install
```
## License

helm-ssm is available under the MIT license. See the LICENSE file for more info.
