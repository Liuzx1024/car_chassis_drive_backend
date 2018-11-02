# Design of Serial-Bus
&emsp;&emsp;Use just one serial port,serve serval device.<br>
&emsp;&emsp;Contain one master and several slaves.
# Hardware

## Master
0. Set CE pin to LOW
1. Send data in buf to Arduino
2. Send FINISH signal to Arduino
3. Copy data to buf from Arduino
4. Waiting for FINISH signal
5. Set CE pin to HIGH

## Slave
0. Enter handler function if interrupt by CE
1. Copy data to buf from raspberryPi
2. Waiting for FINISH signal from raspberryPi
3. Send data in buf to raspberryPi
4. Send FINISH signal to raspberryPi
5. Exit interrupt function

# Software

## Data Flow
It's a string with "\n" to represent the beginning of another JSON message.When all data is sent,put a `FINISH` signal in the data flow.
```
{...}\n
{...}\n
{...}\n
...
FINISH\n
```

### `FINISH` signal
Just a string:
```
FINISH
```