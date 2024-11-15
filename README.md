```
          *          .---. ,---.  ,---.  .-. .-.  ,--,  ,---.         *
         /.\        ( .-._)| .-.\ | .-.\ | | | |.' .')  | .-'        /.\
        /..'\      (_) \   | |-' )| `-'/ | | | ||  |(_) | `-.       /..'\
        /'.'\      _  \ \  | |--' |   (  | | | |\  \    | .-'       /'.'\
       /.''.'\    ( `-'  ) | |    | |\ \ | `-')| \  `-. |  `--.    /.''.'\
       /.'.'.\     `----'  /(     |_| \)\`---(_)  \____\/( __.'    /.'.'.\
"'""""/'.''.'.\""'"'""""""(__)""""""""(__)"""""""""""""(__)""'""""/'.''.'.\""'"'"
      ^^^[_]^^^                                                   ^^^[_]^^^
```

[![Slack][slack-badge]][slack-channel] ( We'll be in `#spruce`)

## Introducing Spruce

`spruce` is a general purpose YAML & JSON merging tool.

It is designed to be an intuitive utility for merging YAML/JSON templates together
to generate complicated YAML/JSON config files in a repeatable fashion. It can be used
to stitch together some generic/top level definitions for the config and pull in overrides
for site-specific configurations to [DRY][dry-definition] your configs up as much as possible.

## How do I get started?

`spruce` is available via Homebrew, just `brew tap starkandwayne/cf; brew install spruce`

Alternatively, you can download a [prebuilt binaries for 64-bit Linux, or Mac OS X][releases]

## How do I compile from source?

1. [Install Go][install-go], e.g. on Ubuntu `sudo snap install --classic go`
1. Fetch sources via `go get github.com/bedag/spruce`
1. Change current directory to the source root `cd ~/go/src/github.com/bedag/spruce/`
1. Compile and execute tests `make all`

## A Quick Example

```sh
# Let's build the first yaml file we will merge
$ cat <<EOF first.yml
some_data: this will be overwritten later
a_random_map:
  key1: some data
heres_an_array:
- first element
EOF

# and now build the second yaml file to merge on top of it
$ cat <<EOF second.yml
some_data: 42
a_random_map:
  key2: adding more data
heres_an_array:
- (( prepend ))
- zeroth element
more_data: 84

# what happens when we spruce merge?
$ spruce merge first.yml second.yml
a_random_map:
  key1: some data
  key2: adding more data
heres_an_array:
- zeroth element
- first element
more_data: 84
some_data: 42
```

The data in `second.yml` is overlayed on top of the data in `first.yml`. Check out the
[merge semantics][merge-semantics] and [array merging][array-merge] for more info on how that was done. Or,
check out [this example on play.spruce.cf][play.spruce-example]

## Documentation

- [What are all the spruce operators, and how do they work?][operator-docs]
- [What are the merge semantics of spruce?][merge-semantics]
- [How can I manipulate arrays with spruce?][array-merge]
- [Can I specify defaults for an operation, or use environment variables?][env-var-defaults]
- [Can I use spruce with go-patch files?][go-patch-support]
- [Can I use spruce with CredHub?][credhub-support]
- [Can I use spruce with Vault?][vault-support]
- [How can I generate spruce templates with spruce itself?][defer]
- [How can I use spruce with BOSH's Cloud Config?][cloud-config-support]

## What else can Spruce do for you?

`spruce` doesn't just stop at merging datastructures together. It also has the following
helpful subcommands:

`spruce diff` - Allows you to get a useful diff of two YAML files, to see where they differ
semantically. This is more than a simple diff tool, as it examines the functional differences,
rather than just textual (e.g. key-ordering differences would be ignored)

`spruce json` - Allows you to convert a YAML document into JSON, for consumption by something
that requires a JSON input. `spruce merge` will handle both YAML + JSON documents, but produce
only YAML output.

`spruce vaultinfo` - Takes a list of files that would be merged together, and analyzes what paths
in Vault would be looked up. Useful for determining explicitly what access an automated process
might need to Vault to obtain the right credentials, and nothing more. Also useful if you need
to audit what credentials your configs are retrieving for a system..

## License

Licensed under [the MIT License][license]


[slack-channel]:        https://cloudfoundry.slack.com/messages/spruce/
[slack-badge]:          http://slack.cloudfoundry.org/badge.svg
[dry-definition]:       https://en.wikipedia.org/wiki/Don%27t_repeat_yourself
[releases]:             https://github.com/bedag/spruce/releases/
[operator-docs]:        https://github.com/bedag/spruce/blob/master/doc/operators.md
[merge-semantics]:      https://github.com/bedag/spruce/blob/master/doc/merging.md
[array-merge]:          https://github.com/bedag/spruce/blob/master/doc/array-merging.md
[env-var-defaults]:     https://github.com/bedag/spruce/blob/master/doc/environment-variables-and-defaults.md
[go-patch-support]:     https://github.com/bedag/spruce/blob/master/doc/merging-go-patch-files.md
[credhub-support]:      https://github.com/bedag/spruce/blob/master/doc/integrating-with-credhub.md
[vault-support]:        https://github.com/bedag/spruce/blob/master/doc/pulling-creds-from-vault.md
[defer]:                https://github.com/bedag/spruce/blob/master/doc/generating-spruce-with-spruce.md
[cloud-config-support]: https://github.com/bedag/spruce/blob/master/doc/integrating-with-cloud-config.md
[license]:              https://github.com/bedag/spruce/blob/master/LICENSE
[install-go]:           https://golang.org/doc/install
