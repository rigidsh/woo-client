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
TBD