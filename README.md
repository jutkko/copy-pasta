# copy-pasta
To use, do the following setup on the two machines you want to `copy-pasta`

```
go get github.com/jutkko/copy-pasta
```

Configure the `secret` file according to the `secret.example` file.

```
. secret # to source the environment variables
```

## Single lined stuff
 To copy, on one machine you do

```
echo "Pasta-copy" | copy-pasta
```

On the other machine you do

```
copy-pasta
```

Boom! you should see

```
Pasta-copy
```

in your terminal.

## Multiline
```
cat myFileOnMachine0 | copy-pasta
```

On the other machine you do

```
copy-pasta > myFileOnMachine1
```

Boom! You should see a copy of `myFileOnMachine0` on your machine 1.

# To test
You will need to have a working go environment, and go to the repo

```
cd copy-pasta
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
TODO: lots to do!
