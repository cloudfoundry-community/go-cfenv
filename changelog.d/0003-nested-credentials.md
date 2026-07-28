[Features]
- `Service.Credential(keys ...string)` reads a credential at any depth and returns it with its own type, so nested values like `protocols.amqp.uri` and non-string leaves like `protocols.amqp.ssl` are reachable. Passing the keys separately means every key is addressable, including one containing a dot such as `jdbc.url`. (#23)
- `Service.CredentialPath("protocols.amqp.uri")` is the dot-delimited form of the same lookup. It cannot address a key that itself contains a dot — use `Credential` for those.

[BugFixes]
- Pinned the credentials payload from #11, whose nested object once panicked the decoder, as a regression test. The defect was already fixed; the test keeps it fixed.
