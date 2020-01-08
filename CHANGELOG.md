# Changelog

## 2.0.x

**NOTE:** Some initial versions of the 2.0.x cycle were wrongly published and for that reason they should start on 2.0.2,
you can check [GitHub releases](https://github.com/codacy/helm-ssm/releases) for the last release with artifacts attached.

- Removed `required` field
- Added `default` field.
  Now you should explicitly say you want the `default=` (empty string) instead of that being implicit.
