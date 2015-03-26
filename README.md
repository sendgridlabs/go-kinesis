# go-kinesis
[![Build Status](https://travis-ci.org/sendgridlabs/go-kinesis.png?branch=master)](https://travis-ci.org/sendgridlabs/go-kinesis)

GO-lang library for AWS Kinesis API.

## [API Documentation](http://godoc.org/github.com/sendgridlabs/go-kinesis)

## Example

Example you can find in folder `examples`.

## Command line interface

You can find a tool for interacting with kinesis from the command line in folder `kinesis-cli`.

## Testing

The tests require a local Kinesis server such as [Kinesalite](https://github.com/mhart/kinesalite)
to be running and reachable at `http://127.0.0.1:4567`.

To make the tests complete faster, you might want to have Kinesalite perform stream creation and
deletion faster than the default of 500ms, like so:

    kinesalite --createStreamMs 5 --deleteStreamMs 5 &

The `&` runs Kinesalite in the background, which is probably what you want.
