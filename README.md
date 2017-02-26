[![Coverage Status](https://coveralls.io/repos/github/jutkko/copy-pasta/badge.svg?branch=master)](https://coveralls.io/github/jutkko/copy-pasta?branch=master)
[![Build Status](https://travis-ci.org/jutkko/copy-pasta.svg?branch=master)](https://travis-ci.org/jutkko/copy-pasta)

# How to use
## Single lined stuff
 To copy, on one machine you do

```
echo "I don't like ravioli" | copy-pasta
```

On the other machine you do

```
copy-pasta
```

Boom! you should see

```
I don't like ravioli
```

in the terminal.

## Multiline / Files
```
cat myPenne.jpg | copy-pasta
```

On the other machine you do

```
copy-pasta > myPenne.jpg
```

Boom! You should see a copy of `myPenne.jpg` on your other machine.

## Multi-user
Are you sharing a machine with others? Or you want to have multiple clipboards?
`copy-pasta` now supports [concourse](https://concourse.ci) `fly` like targets.
Remember the `--target` option in the `login` command?  After specifying
another user like

```
copy-pasta login --target your-copy-pasta
<Enter your S3 accesskey>
<Enter your S3 secretaccesskey>
```

You can do

```
copy-pasta target your-copy-pasta
```

You will be using another `copy-pasta` destination. **Note the credentials can
be the same one!**

# Installation
To use, do the following setup on the two machines you want to `copy-pasta`

```
go get github.com/jutkko/copy-pasta
```

Login on the machines you want to do `copy-pasta`

```
copy-pasta login --target my-copy-pasta
<Enter your S3 accesskey>
<Enter your S3 secretaccesskey>
```

If you are not using Amazon S3, and not in London you might want to pass the
`endpoint` and `location` of your S3 backend implementation when you target.

# Running the tests
You will need to have a working go environment, and go to the repo

```
cd $GOPATH/src/github.com/jutkko/copy-pasta
```

Install the awesome ginkgo testing framework

```
go get github.com/onsi/ginkgo
go get github.com/onsi/gomega
```

To run the tests

```
ginkgo -r
```

# To contribute
## Features, concerns or bugs
Please open an issue and talk about the feature/bug you have, I will get back
to you very soon.

## Use cases
Got an interesting use case for `copy-pasta`? Make a PR and I will include it
here! Here's some

### Bash
Non-interactive logon

```
printf "%s\n%s\n" "$ACCESSKEY" "$SECRETACCESSKEY" | copy-pasta login --target my-target
```

Paste straight to pbcopy
```
#!/bin/bash
copy-pasta-to-pbcopy() {
  copy-pasta | pbcopy
}
```

Copy straight into copy-paste

```
#!/bin/bash
pbpaste-to-copy-pasta() {
  pbpaste | copy-pasta
}
```

... And yours?

# copy-pasta?
Credits to my colleague [Vlad](https://github.com/vlad-stoian). Genius!
