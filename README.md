# copy-pasta
To use, do the following setup on the two machiens you want to `copy-pasta`

```
go get github.com/jutkko/copy-pasta
```

Configure the `secret` file according to the `secret.example` file, and to copy, on one machine you do:

```
. secret # to source the environment variables
echo "Pasta-copy" | copy-pasta
```

On the other machine you do:

```
copy-pasta
```

Boom! you should see

```
Pasta-copy
```

in your terminal.

# To test
TODO: lots to do!

# To contribute
TODO: lots to do!
