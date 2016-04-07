Thank you for contributing goma!  Here are few guidelines.

Issues / Bugs
-------------

If you think you found a bug, open an issue and supply the
minimum configuration that triggers the bug reproducibly.

Pull requests
-------------

New plugins should be accompanied with enough tests.  
Bug fixes should update or add tests to cover that bug.

Codes must be formatted by [goimports][].  Please configure
your editor to run it automatically.

Pull requests will be tested on travis-ci with [golint][] and
[go vet][govet].  Please run them on your local machine before
submission.

By submitting the code, you agree to license your code under
the [MIT][] license.

[gofmt]: https://golang.org/cmd/gofmt/
[goimports]: https://godoc.org/golang.org/x/tools/cmd/goimports
[golint]: https://github.com/golang/lint
[govet]: https://golang.org/cmd/vet/
[MIT]: https://opensource.org/licenses/MIT
