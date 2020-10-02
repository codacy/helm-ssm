# Changelog

## 3.x

Helm 2 support is dropped and dependencies on its api were removed. This means helm extensions to go templates
are not supported by this plugin, since the helm engine funcMap is now private. The functions from the
[sprig](http://masterminds.github.io/sprig/) are now used instead.

## 2.0.x

**NOTE:** Some initial versions of the 2.0.x cycle were wrongly published and for that reason they should start on 2.0.2,
you can check [GitHub releases](https://github.com/codacy/helm-ssm/releases) for the last release with artifacts attached.

- Removed `required` field
- Added `default` field.
  Now you should explicitly say you want the `default=` (empty string) instead of that being implicit.
- The input file flag `-t` was changed to `-o`