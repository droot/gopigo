package main

import (
  "github.com/kidoman/embd"
)

func I2CBus() embd.I2CBus {
  if err := embd.InitI2C(); err != nil {
    panic(err)
  }
  return embd.NewI2CBus(1)
}

func cleanupI2Bus() {
  embd.CloseI2C()
}
