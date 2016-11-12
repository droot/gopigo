# gopigo
The GoPiGo is a delightful and complete robot for the Raspberry Pi that turns your Pi into a fully operating robot.  GoPiGo is a mobile robotic platform for the Raspberry Pi developed by [Dexter Industries.](http://www.dexterindustries.com/GoPiGo)  

![ GoPiGo ](https://raw.githubusercontent.com/DexterInd/GoPiGo/master/GoPiGo_Chassis-300.jpg)

This repository contains Go library for interacting with the GoPiGo robot. This is a port of python library of GoPigo in Go language.

## Getting Started

Get the source with: `go get -d -u github.com/droot/gopigo`

#### Simple example

```go
package main

import (
  "github.com/droot/gopigo"
  "github.com/kidoman/embd"

  _ "github.com/kidoman/embd/host/all"
)

func main() {
  if err := embd.InitI2C(); err != nil {
    panic(err)
  }
  defer embd.CloseI2C()

  bus := embd.NewI2CBus(1)

  // create GoPiGo instance.
  gp := gopigo.New(bus)

  // move GoPiGo forward by 40 cms
  err := gp.fwd(40)

  // move GoPiGo backward
  err = gp.Bwd(0)

  // make GoPiGo stop
  err = gp.Stop()

  // make GoPiGo take Left turn
  err = gp.Left()

  // make GoPiGo take right turn
  err = gp.Right()


  // Read Battery level of GoPiGo
  volt, err := gp.BatteryVoltage()


  .....
}

```

## API Documentation
API documentation is available at https://godoc.org/github.com/droot/gopigo

## Credits
Special Thanks to Dexter Industries for providing [@GoPiGo kit](http://www.dexterindustries.com/shop/gopigo-starter-kit-2/) for the development.

## Need help ?
 * Issues: https://github.com/droot/gopigo/issues
 * twitter: [@droot](https://twitter.com/_sunil_)
