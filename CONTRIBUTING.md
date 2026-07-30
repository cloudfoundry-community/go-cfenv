# Thank you for your interest in contributing to Cloud Foundry

We try to tag issues that should be reasonable for a new contributor to take on with a [good first issue](https://github.com/search?q=org%3Acloudfoundry+org%3Acommunity+label%3A%22good+first+issue%22+state%3Aopen&type=Issues) label so you have somewhere to start. 

**How to contribute:**
- Fork the repo you'd like to make a contribution to
- Clone your fork to your local workstation
- Create a new branch for the issue
- Make the necessary changes on that branch
- Commit and push to that branch
- Make a Pull Request against the repo
- Sign the Contributor Licensing Agreement, if you haven’t already

All committers to a Cloud Foundry Foundation project must sign a Contributor License Agreement. 
To sign these documents and begin contributing, please [sign in to the EasyCLA app](https://corporate.v1.easycla.lfx.linuxfoundation.org/). 
Or, you’ll simply be prompted to do so when you put in your first Pull Request.

If you’d like to have someone review the agreement before signing, you may download them here: 
- [Individual Contributor License Agreement](https://www.cloudfoundry.org/wp-content/uploads/icla.pdf)
- [Corporate Contributor License Agreement](https://www.cloudfoundry.org/wp-content/uploads/ccla.pdf)

**Testing policy**

New functionality must arrive with tests covering it, and a bug fix must arrive
with a test that fails without the fix. That second half matters more than it
sounds: a test written after the fix proves only that the code does what it
does, never that it catches the bug.

`make check` and `make test-race` must both pass before a pull request is
merged, and CI enforces both on every pull request. Run `make cover` to see
where a change leaves coverage.

Parsing changes deserve a look at `make fuzz` as well. `VCAP_APPLICATION` and
`VCAP_SERVICES` come from the platform and from bound service brokers rather
than from the application, so this package is a trust boundary: it must return
an error or a value for any input, and never panic.

**Interim versions**

Work is reviewable between releases, not only at release time. Every change
lands on `master` as its own commit through a pull request, and interim
versions are tagged as `dev` prereleases so reviewers have something citable
to build and test against:

```bash
make bump dev                              # v1.23.0 -> v1.24.0-dev.1
git push origin refs/tags/v1.24.0-dev.1
make bump dev                              # -dev.1  -> -dev.2
make bump final                            # -dev.2  -> v1.24.0, ready to release
```

The prerelease is derived from the *next* release, not the last one. Semver
orders a prerelease before its own release, so tagging `v1.23.0-dev.1` after
`v1.23.0` would produce an interim version older than the release it follows,
which no resolver would select. Which component advances is
`BUMP_PRERELEASE_STEP` in `mk/bump.mk`; this project releases minor versions,
which is the default.

**Where can I reach out to the team?**

- _Want to report concerns/bugs?_ Create an issue on the affected repo.
- _Usage issues/help?_ Reach out to us on [Slack](https://slack.cloudfoundry.org/) or our [cf-dev mailing list](https://lists.cloudfoundry.org/g/cf-dev/topics).
- _Want to participate in deeper architectural discussions?_ Attend a relevant [working group meeting](https://github.com/cloudfoundry/community/blob/main/toc/working-groups/WORKING-GROUPS.md) or our monthly CAB call. 

You’ll find these listed on our [community calendar](https://www.cloudfoundry.org/community-calendar/).
