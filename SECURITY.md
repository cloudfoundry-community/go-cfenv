# Security Policy

The Cloud Foundry Foundation (CFF) Security Team provides a single point of contact for the reporting of security vulnerabilities in open source Cloud Foundry codebases and coordinates the process of investigating any reports. Please see [this page](https://cve.mitre.org/about/terminology.html) for more information about what might qualify as a vulnerability.

## Reporting a Vulnerability

We strongly encourage people to report security vulnerabilities privately to our security team before disclosing them in a public forum.

The e-mail address to use to contact the CFF Security Team is **security@cloudfoundry.org**.

Please note that the e-mail address above should only be used for reporting undisclosed security vulnerabilities in open source Cloud Foundry codebases and managing the process of fixing such vulnerabilities. We cannot accept regular bug reports or other security-related queries at this address.

If you wish to send encrypted email, our public key can be obtained from a public key server such as keys.openpgp.org. The fingerprint is: `3FC8 9AF3 940B E270 CF25  E122 9965 0006 EF9D C642`

## Maintenance and review

This library has a single maintainer. GitHub does not allow a pull request to be approved by its author, so no change here can carry an approving review, and the OpenSSF Scorecard `Code-Review` check scores 0 as a consequence. That score is accepted rather than overlooked: no repository-side setting can raise it while the maintainer count is one.

What is enforced on `master` in place of a second reviewer:

- Changes land only through a pull request. Direct pushes, force pushes and branch deletion are blocked, and review conversations must be resolved before a merge.
- Five status checks must pass, on a branch up to date with `master`: `Build at go directive floor`, `Test`, `Quality gates`, `Security audit` and `Analyze Go`.
- `Security audit` runs govulncheck, gosec, gitleaks and modrot on every pull request, and `Analyze Go` runs CodeQL. OpenSSF Scorecard runs on every push to `master` and weekly.
- Every GitHub Action is pinned to a full commit SHA, with Dependabot proposing updates weekly.

This is revisited if a second maintainer joins the project: requiring an approving review becomes possible at that point, and the check can pass on its merits.
