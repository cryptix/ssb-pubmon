# ssb-pubmon

a botpage that periodically pings [secure scuttlebutt](https://www.scuttlebutt.nz/) pubs to monitor their livelyness.

![](screenshot.png)

# TODO

* make a list of pubs
* make a worker that iterates over them and attempt shs/muxrpc connections to them.

## QOR maintainance

* overwrite auth template
* instantiate own middleware
* don't use global default cookie secret
* go to address > try to change pub: dropdown has no names?

# Fork of QOR example application

see [qor-example](https://github.com/qor/qor-example).

## License

Released under the MIT License.

