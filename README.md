# Speechmatics Terraform Provider

## Testing Locally

Ensure you have your `GOBIN` env variable set: `go env GOBIN`

You will need to add the following to `~/.terraformrc` so terraform looks locally for custom providers:

```hcl
provider_installation {

  dev_overrides {
      "hashicorp.com/edu/speechmatics" = "YOUR_GOBIN_PATH"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

After making changes, run `go install .` which will build your new provider to your GOBIN path