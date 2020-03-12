# Helm SSM Plugin

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/d3cd080edd8644e085f2f8adfd43510c)](https://www.codacy.com?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=codacy/helm-ssm&amp;utm_campaign=Badge_Grade)
[![CircleCI](https://circleci.com/gh/codacy/helm-ssm.svg?style=svg)](https://circleci.com/gh/codacy/helm-ssm)

This is a plugin to help developers inject values coming from AWS SSM
parameters, on the `values.yaml` file.

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

```

when you do not want a key to be defined, use a with using a default empty value in the go templates:

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
  -t, --target-dir string       dir to output content
  -f, --values valueFilesList   specify values in a YAML file (can specify multiple) (default [])
  -v, --verbose                 show the computed YAML values file/s
```

## Example

[![asciicast](https://asciinema.org/a/c2zut95zzbiKyk5gJov67bxsP.svg)](https://asciinema.org/a/c2zut95zzbiKyk5gJov67bxsP?t=1)

## Install

Choose the latest version from the releases and install the
appropriate version for your OS:

### Linux

```sh
$ helm plugin install https://github.com/codacy/helm-ssm/releases/download/latest/helm-ssm-linux.tgz
```

### MacOS

```sh
$ helm plugin install https://github.com/codacy/helm-ssm/releases/download/latest/helm-ssm-macos.tgz
```

### Windows

```sh
$ helm plugin install https://github.com/codacy/helm-ssm/releases/download/latest/helm-ssm-windows.tgz
```

### Developer (From Source) Install

If you would like to handle the build yourself, instead of fetching a binary,
this is how recommend doing it.

First, set up your environment:

- You need to have [Go](http://golang.org) installed. Make sure to set `$GOPATH`

Clone this repo into your `$GOPATH`. You can use `go get -d github.com/codacy/helm-ssm`
for that.

```sh
$ cd $GOPATH/src/github.com/codacy/helm-ssm
$ make install
```

## What is Codacy

[Codacy](https://www.codacy.com/) is an Automated Code Review Tool that monitors your technical debt, helps you improve your code quality, teaches best practices to your developers, and helps you save time in Code Reviews.

### Among Codacy’s features

- Identify new Static Analysis issues
- Commit and Pull Request Analysis with GitHub, BitBucket/Stash, GitLab (and also direct git repositories)
- Auto-comments on Commits and Pull Requests
- Integrations with Slack, HipChat, Jira, YouTrack
- Track issues in Code Style, Security, Error Proneness, Performance, Unused Code and other categories

Codacy also helps keep track of Code Coverage, Code Duplication, and Code Complexity.

Codacy supports PHP, Python, Ruby, Java, JavaScript, and Scala, among others.

## Free for Open Source

Codacy is free for Open Source projects.

## License

helm-ssm is available under the MIT license. See the LICENSE file for more info.
