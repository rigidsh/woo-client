Implementation of BLE woo sport protocol on Go

# Protocol
## Transport protocol
### Package structure
```
<PackageStart> [data...] [checksum] <PackageEnd>
```
if in data there is a control symbol, it will be escaped with EscapeSymbol

### Control symbols
| Symbol       | Byte |
|--------------|------|
| PackageStart | 0xD1 |
| PackageEnd   | 0xDF |
| EscapeSymbol | 0xDE |

### Checksum
Sum of all data bytes plus checksum byte should be equal to 0x00 

## Logical protocol
see [woo.dump](https://github.com/rigidsh/woo-client/blob/main/woo.dump) as an example (recorder from a real device)
### Event types
| Event          | Byte |
|----------------|------|
| RecordingEvent | 0x58 |
| JumpEventType  | 0x44 |
### RecordingEvent
| Field            | Size | Description                   |
|------------------|------|-------------------------------|
| Start            | 1    | 0x01 if recording was started |
| Unknown          | 1    | always 0x00                   |
| Number of events | 2    | Events in memory              |
| Counter          | 4    | ???                           |
| Unknown          | 4    | always 0x00 0x00 0xF6 0x01    |
### JumpEvent
| Field      | Size | Description                  |
|------------|------|------------------------------|
| JumpType   | 1    | always 0x02                  |
| JumpNumber | 1    | Event number after last sync |
| Unknown    | 10   |                              |
| JumpHeight | 2    | Unknown format               |
| Unknown    | 13   |                              |
| JumpTime   | 6    | Woo date/time format. in UTC |
| Unknown    | 14   |                              |