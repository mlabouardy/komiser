We are always thrilled to receive pull requests, and do our best to process them as fast as possible. Not sure if that typo is worth a pull request? Do it! We will appreciate it.

`Note: If your pull request is not accepted on the first try, don’t be discouraged! If there’s a problem with the implementation, hopefully you received feedback on what to improve.`

## Guidelines 

We recommend discussing your plans on our Slack (join our <a href="https://community.komiser.io/">community channel</a>) before starting to code - especially for more ambitious contributions. This gives other contributors a chance to point you in the right direction, give feedback on your design, and maybe point out if someone else is working on the same thing.

Any significant improvement should be documented as a github issue before anybody starts working on it. Please take a moment to check that an issue doesn’t already exist documenting your bug report or improvement proposal. If it does, it never hurts to add a quick “+1” or “I have this problem too”. This will help prioritize the most common problems and requests

## Conventions

Fork the repo and make changes on your fork in a feature branch based on the master branch:

- If it’s a bugfix branch, name it fix/XXX-something where XXX is the number of the issue
- If it’s a feature branch, create an enhancement issue to announce your intentions, and name it feature/XXX-something where XXX is the number of the issue.
- Submit unit tests for your changes. Go has a great test framework built in; use it! Take a look at existing tests for inspiration. Run the full test suite on your branch before submitting a pull request.
- Make sure you include relevant updates or additions to documentation when creating or modifying features.
- Write clean code. Universally formatted code promotes ease of writing, reading, and maintenance. Always run go fmt before committing your changes. Most editors have plugins that do this automatically.

## How Can I Contribute to Komiser Dashboard?

* Ensure that your go version supports `go.mod` (>= 1.11)

* Clone the project
* Install Go dependencies:

```
go mod download
```

* Switch to public folder and install npm dependencies and angular cli:

```
cd public
npm install -g @angular/cli
npm install
```

* Deploy to a local server:

```
ng serve
```

* Once you implemented your frontend changes, build the artifact, build it as golang assets:

```
go-bindata-assetfs -o template.go dist/ dist/assets/images/
```
