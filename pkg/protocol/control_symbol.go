package protocol

type controlSymbol byte

const (
	packageStart controlSymbol = 0xD1
	packageEnd   controlSymbol = 0xDF
	escapeSymbol controlSymbol = 0xDE
)

func findNextControlSymbol(data sliceWithBuffer, skipFirst bool) (controlSymbol, int) {
	if data.buffer.hasData {
		if skipFirst {
			skipFirst = false
		} else {
			if data.buffer.data == byte(packageStart) || data.buffer.data == byte(packageEnd) || data.buffer.data == byte(escapeSymbol) {
				return controlSymbol(data.buffer.data), 0
			}
		}
	}
	for i, b := range data.data {
		if skipFirst && i == 0 {
			continue
		}
		if b == byte(packageStart) || b == byte(packageEnd) || b == byte(escapeSymbol) {
			return controlSymbol(b), data.FromDataIndex(i)
		}
	}

	return 0, -1
}
