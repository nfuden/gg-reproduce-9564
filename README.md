# Gloo-9317 Reproducer

## Notes
This has been inspired by Duncan's reproducers and strives to show a version of what we can do via mage.
In addition to running these commands please feel free to give feedback on the structure found within.

## Installation

Export your Gloo Edge License Key to an environment variable:
```
export GLOO_EDGE_LICENSE_KEY={your license key}
```

Make sure you have Mage https://magefile.org/

```
brew install mage
```

or 

```
cd ~/yourworkdir
git clone https://github.com/magefile/mage
cd mage
go run bootstrap.go
```

## Setup the environment


Install the environment and apply resources
```
mage fullinstall
mage demo:run2000
```

now install the new version 

```
glooctl uninstall --all
mage installgloo true
mage demo:run2000
```
