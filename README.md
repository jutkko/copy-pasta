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
<Enter your S3 ACCESSKEY>
<Enter your S3 SECRETACCESSKEY>
```

You can do

```
copy-pasta target your-copy-pasta
```

You will be using another `copy-pasta` destination. **Note the credentials can
be the same one!**

## How does it work?
Are you super paranoid about security? Do you sweat if you copy your
credentials into a copy buffer and leave it there? Then you should read on.
Here is a diagram that briefly describes how `copy-pasta` works.

<img src="/figures/how-it-works.png" width="750">

The communication between the machines and the storage server is done in `SSL`,
so we can assume that it is relatively safe.

We can see that the things you copy into `copy-pasta` gets stored in plain text
on the storage server. The weakest link here will be the security of your
backend store. Take S3 as an example, if your bucket is private and you haven't
shared with anyone your `ACCESSKEY` and `SECRETACCESSKEY`, you should be pretty
safe. On the other hand, if the backend store is either public or compromised,
the content copied to `copy-pasta` is in danger.

In general it is **not** advised to copy confidential content to `copy-pasta`,
`copy-pasta` is also **not** responsible for keeping the content secure. But if
you are a security lax person like me, you probably can take the advantage of
the overwrite nature of `copy-pasta`, copy confidential content, use it and
quickly copy something else into `copy-pasta`.

# Installation
Looking good? Can't wait to hack with `copy-pasta`? There are two ways to
install `copy-pasta`. Using go, do the following setup on the two machines you
want to `copy-pasta`

```
go get github.com/jutkko/copy-pasta
```

Using `homebrew`, do

```
brew tap jutkko/homebrew-copy-pasta
brew install copy-pasta
```

Login on the machines you want to do `copy-pasta`

```
copy-pasta login --target my-copy-pasta
<Enter your S3 ACCESSKEY>
<Enter your S3 SECRETACCESSKEY>
```

If you are not using Amazon S3, or your bucket location is  not in London you
might want to pass the `endpoint` and `location` of your S3 backend
implementation when you target.

If you are interested in using another storage solution, please let me know
is the issues page and we get the conversation started.

# Uninstall
It depends on how you installed the binary. If by go, you should remove both
the `copy-pasta` repo and the binary

```
rm -rf $GOPATH/src/github.com/jutkko/copy-pasta
rm $GOPATH/bin/copy-pasta
```

If by homebrew, you can first remove the binary and then the tap

```
brew uninstall copy-pasta
brew untap homebrew-copy-pasta
```

To remove the config file leftover by `copy-pasta`, simply delete the
`.copy-pastarc` file in your `$HOME`.

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
