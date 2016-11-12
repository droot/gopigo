package main

import (
  "github.com/kidoman/embd"
)

type mockI2CBus struct{}

// ReadByte reads a byte from the given address.
func (m *mockI2CBus) ReadByte(addr byte) (value byte, err error) {
  return 0, nil
}

// ReadBytes reads a slice of bytes from the given address.
func (m *mockI2CBus) ReadBytes(addr byte, num int) (value []byte, err error) {
  return nil, nil
}

// WriteByte writes a byte to the given address.
func (m *mockI2CBus) WriteByte(addr, value byte) error {
  return nil
}

// WriteBytes writes a slice bytes to the given address.
func (m *mockI2CBus) WriteBytes(addr byte, value []byte) error {
  return nil
}

// ReadFromReg reads n (len(value)) bytes from the given address and register.
func (m *mockI2CBus) ReadFromReg(addr, reg byte, value []byte) error {
  return nil
}

// ReadByteFromReg reads a byte from the given address and register.
func (m *mockI2CBus) ReadByteFromReg(addr, reg byte) (value byte, err error) {
  return 0, nil
}

// ReadU16FromReg reads a unsigned 16 bit integer from the given address and register.
func (m *mockI2CBus) ReadWordFromReg(addr, reg byte) (value uint16, err error) {
  return 0, nil
}

// WriteToReg writes len(value) bytes to the given address and register.
func (m *mockI2CBus) WriteToReg(addr, reg byte, value []byte) error {
  return nil
}

// WriteByteToReg writes a byte to the given address and register.
func (m *mockI2CBus) WriteByteToReg(addr, reg, value byte) error {
  return nil
}

// WriteU16ToReg
func (m *mockI2CBus) WriteWordToReg(addr, reg byte, value uint16) error {
  return nil
}

// Close releases the resources associated with the bus.
func (m *mockI2CBus) Close() error {
  return nil
}

func I2CBus() embd.I2CBus {
  return &mockI2CBus{}
}

func cleanupI2Bus() {
}
