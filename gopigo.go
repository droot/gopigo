// package gopigo implements APIs to control a GoPiGo robot.
// This package has been tested using Raspberry PI 3. If you have any trouble
// using it with older version of PIs, report it in Github issues.
package gopigo

import (
  "github.com/kidoman/embd"

  "fmt"
  _ "github.com/kidoman/embd/host/all"
  "math"
  "sync"
  "time"
)

const (
  // i2c device address for GoPiGo device.
  address = 0x08
)

const (
  wheelRadius       = 3.25
  wheelDistance     = 2 * math.Pi * wheelRadius
  pulsesPerRotation = 18
)

// Command values supported by GoPiGo.
const (
  fwdCmd            = 119
  motorFwdCmd       = 105
  motorBwdCmd       = 107
  bwdCmd            = 115
  stopCmd           = 120
  leftCmd           = 97  // Turn Left by turning off one motor
  leftRotCmd        = 98  // Rotate left by running both motors is opposite direction
  rightCmd          = 100 // Turn Right by turning off one motor
  rightRotCmd       = 110 // Rotate Right by running both motors is opposite direction
  ispdCmd           = 116 // Increase the speed by 10
  dspdCmd           = 103 // Decrease the speed by 10
  m1Cmd             = 111 // Control motor1
  m2Cmd             = 112 // Control motor2
  readMotorSpeedCmd = 114 // Get motor speed back
  voltCmd           = 118 // Read the voltage of the batteries
  usCmd             = 117 // Read the distance from the ultrasonic sensor
  ledCmd            = 108 // Turn On/Off the LED's
  servoCmd          = 101 // Rotate the servo
  encTgtCmd         = 50  // Set the encoder targeting
  fwVerCmd          = 20  // Read the firmware version
  enEncCmd          = 51  // Enable the encoders
  disEncCmd         = 52  // Disable the encoders
  readEncStatusCmd  = 53  // Read encoder status
  enServoCmd        = 61  // Enable the servo's
  disServoCmd       = 60  // Disable the servo's
  setLeftSpeedCmd   = 70  // Set the speed of the right motor
  setRightSpeedCmd  = 71  // Set the speed of the left motor
  enComTimeoutCmd   = 80  // Enable communication timeout
  disComTimeoutCmd  = 81  // Disable communication timeout
  timeoutStatusCmd  = 82  // Read the timeout status
  encReadCmd        = 53  // Read encoder values
  trimTestCmd       = 30  // Test the trim values
  trimWriteCmd      = 31  // Write the trim values
  trimReadCmd       = 32

  digitalWriteCmd = 12 // Digital write on a port
  digitalReadCmd  = 13 // Digital read on a port
  analogReadCmd   = 14 // Analog read on a port
  analogWriteCmd  = 15 // Analog read on a port
  pinModeCmd      = 16 // Set up the pin mode on a port

  irReadCmd    = 21
  irRecvPinCmd = 22
  cpuSpeedCmd  = 25
)

// GoPiGo represents a GoPiGo device.
type GoPiGo struct {
  mu  sync.Mutex
  bus embd.I2CBus
}

// New instantiates GoPiGo devices given I2C bus device.
func New(bus embd.I2CBus) *GoPiGo {
  return &GoPiGo{
    bus: bus,
  }
}

// cmdValues is a helper function to construct the command byte array.
func cmdValues(mainCmd byte, cmdParams ...byte) []byte {
  values := make([]byte, 4)
  values[0] = mainCmd
  for i, v := range cmdParams {
    values[i+1] = v
  }
  return values
}

func (p *GoPiGo) sendCmd(cmdValues []byte) error {
  err := p.bus.WriteToReg(address, 1, cmdValues)
  time.Sleep(5 * time.Millisecond)
  return err
}

// writeByte is helper function to write a byte to the GoPiGo.
func (p *GoPiGo) writeByte(b byte) error {
  return p.bus.WriteByte(address, b)
}

// readByte reads a byte at the GoPiGo address.
func (p *GoPiGo) readByte() (byte, error) {
  return p.bus.ReadByte(address)
}

// DirectMotor1 sets the direction and speed of motor1.
func (p *GoPiGo) DirectMotor1(direction, speed byte) error {
  p.mu.Lock()
  defer p.mu.Unlock()
  return p.sendCmd(cmdValues(m1Cmd, direction, speed))
}

// DirectMotor2 sets the direction and speed of motor2.
func (p *GoPiGo) DirectMotor2(direction, speed byte) error {
  p.mu.Lock()
  defer p.mu.Unlock()
  return p.sendCmd(cmdValues(m2Cmd, direction, speed))
}

// Fwd moves the GoPiGo forward. It takes distance in cms as input. If specified
// 0, it moves forward without stopping.
func (p *GoPiGo) Fwd(distance int) error {
  p.mu.Lock()
  defer p.mu.Unlock()
  if distance > 0 {
    pulses := pulsesPerRotation * (float64(distance) / wheelDistance)
    if err := p.setEncoderTarget(true, true, int(pulses)); err != nil {
      return err
    }
  }
  return p.sendCmd(cmdValues(motorFwdCmd))
}

// Bwd moves the GoPiGo backwards. It takes distance in cms as input. If distance
// is 0, it moves backward without stopping.
func (p *GoPiGo) Bwd(distance int) error {
  p.mu.Lock()
  defer p.mu.Unlock()
  if distance > 0 {
    pulses := pulsesPerRotation * (float64(distance) / wheelDistance)
    if err := p.setEncoderTarget(true, true, int(pulses)); err != nil {
      return err
    }
  }
  return p.sendCmd(cmdValues(motorBwdCmd))
}

// Left turns GoPiGo left slowly.
func (p *GoPiGo) Left() error {
  p.mu.Lock()
  defer p.mu.Unlock()
  return p.sendCmd(cmdValues(leftCmd))
}

// LeftRotate turns GoPiGo left more aggresively by rotating both the wheels in
// reverse direction.
func (p *GoPiGo) LeftRotate() error {
  p.mu.Lock()
  defer p.mu.Unlock()
  return p.sendCmd(cmdValues(leftRotCmd))
}

// TurnLeft turns GoPiGo left by specified degrees.
func (p *GoPiGo) TurnLeft(degrees float64) error {
  p.mu.Lock()
  defer p.mu.Unlock()
  if degrees > 0 {
    DPR := 360.0 / 64
    pulses := int(degrees / DPR)
    p.setEncoderTarget(false, true, pulses)
  }
  return p.sendCmd(cmdValues(leftCmd))
}

// Right turns GoPiGo to the right slowly.
func (p *GoPiGo) Right() error {
  p.mu.Lock()
  defer p.mu.Unlock()
  return p.sendCmd(cmdValues(rightCmd))
}

// RightRotate turns GoPiGo to the right aggresively.
func (p *GoPiGo) RightRotate() error {
  p.mu.Lock()
  defer p.mu.Unlock()
  return p.sendCmd(cmdValues(rightRotCmd))
}

// TurnRight turns GoPiGo to the right by specified angle in degrees.
func (p *GoPiGo) TurnRight(degrees float64) error {
  p.mu.Lock()
  defer p.mu.Unlock()
  if degrees > 0 {
    DPR := 360.0 / 64
    pulses := int(degrees / DPR)
    p.setEncoderTarget(true, false, pulses)
  }
  return p.sendCmd(cmdValues(rightCmd))
}

// Stop stops the GoPiGo devices (if its in motion).
func (p *GoPiGo) Stop() error {
  p.mu.Lock()
  defer p.mu.Unlock()
  return p.sendCmd(cmdValues(stopCmd))
}

// IncreaseSpeed bumps up the motor's speed by 10.
func (p *GoPiGo) IncreaseSpeed() error {
  p.mu.Lock()
  defer p.mu.Unlock()
  return p.sendCmd(cmdValues(ispdCmd))
}

// DecreaseSpeed slows down GoPiGo by speed of 10.
func (p *GoPiGo) DecreaseSpeed() error {
  p.mu.Lock()
  defer p.mu.Unlock()
  return p.sendCmd(cmdValues(dspdCmd))
}

// setEncoderTarget sets the encoder targeting to number of target pulses.
func (p *GoPiGo) setEncoderTarget(motor1, motor2 bool, targetPulses int) error {
  if targetPulses < 0 {
    return fmt.Errorf("invalid target Pulse value")
  }
  var m1, m2 byte

  if motor1 {
    m1 = 1
  }
  if motor2 {
    m2 = 1
  }
  m := 2*m1 + m2
  return p.sendCmd(cmdValues(encTgtCmd, m, byte(targetPulses/256), byte(targetPulses%256)))
}

// EnableEncoders enables the encoders (enabled by default).
func (p *GoPiGo) EnableEncoders() error {
  p.mu.Lock()
  defer p.mu.Unlock()
  return p.sendCmd(cmdValues(enEncCmd))
}

// DisableEncoders disables encoders.
func (p *GoPiGo) DisableEncoders() error {
  p.mu.Lock()
  defer p.mu.Unlock()
  return p.sendCmd(cmdValues(disEncCmd))
}

// func (p *GoPiGo) TrimTest(value byte) error {
//   if value > 100 {
//     value = 100
//   } else if value < -100 {
//     value = -100
//   }
//   value += 100
//   return p.sendCmd(cmdValues(trimTestCmd, value))
// }

// func (p *GoPiGo) TrimRead() byte {
//   err := p.sendCmd(cmdValues(trimReadCmd))
//   if err != nil {
//     return -1
//   }
//   b1, err := p.readByte()
//   if err != nil {
//     return -1
//   }
//   b2, err := p.readByte()
//   if err != nil {
//     return -1
//   }

//   if b1 != -1 && b2 != -1 {
//     val := b1*256 + b2
//     if val == 255 {
//       return -3
//     }
//     return val
//   }
//   return -1
// }

// func (p *GoPiGo) TrimWrite(val byte) error {
//   if val > 100 {
//     val = 100
//   } else if val < 100 {
//     val = -100

// BatteryVoltage returns the Battery Voltage reading.
func (p *GoPiGo) BatteryVoltage() (volt float32, err error) {
  p.mu.Lock()
  defer p.mu.Unlock()

  err = p.sendCmd(cmdValues(voltCmd))
  b1, err := p.readByte()
  if err != nil {
    return
  }
  b2, err := p.readByte()
  if err != nil {
    return
  }

  v := b1*255 + b2
  volt = (5 * float32(v) / 1024) / 0.4
  return
}
