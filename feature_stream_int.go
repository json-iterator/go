package jsoniter

var digits []uint8
var digitTens []uint8
var digitOnes []uint8

func init() {
	digits = []uint8{
		'0', '1', '2', '3', '4', '5',
		'6', '7', '8', '9', 'a', 'b',
		'c', 'd', 'e', 'f', 'g', 'h',
		'i', 'j', 'k', 'l', 'm', 'n',
		'o', 'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y', 'z',
	}
	digitTens = []uint8{
		'0', '0', '0', '0', '0', '0', '0', '0', '0', '0',
		'1', '1', '1', '1', '1', '1', '1', '1', '1', '1',
		'2', '2', '2', '2', '2', '2', '2', '2', '2', '2',
		'3', '3', '3', '3', '3', '3', '3', '3', '3', '3',
		'4', '4', '4', '4', '4', '4', '4', '4', '4', '4',
		'5', '5', '5', '5', '5', '5', '5', '5', '5', '5',
		'6', '6', '6', '6', '6', '6', '6', '6', '6', '6',
		'7', '7', '7', '7', '7', '7', '7', '7', '7', '7',
		'8', '8', '8', '8', '8', '8', '8', '8', '8', '8',
		'9', '9', '9', '9', '9', '9', '9', '9', '9', '9',
	}
	digitOnes = []uint8{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	}
}
func (stream *Stream) WriteUint8(val uint8) {
	if stream.Available() < 3 {
		stream.Flush()
	}
	charPos := stream.n
	if val <= 9 {
		charPos += 1;
	} else {
		if val <= 99 {
			charPos += 2;
		} else {
			charPos += 3;
		}
	}
	stream.n = charPos
	var q uint8
	for {
		q = val / 10
		r := val - ((q << 3) + (q << 1))  // r = i-(q*10) ...
		charPos--
		stream.buf[charPos] = digits[r]
		val = q;
		if val == 0 {
			break
		}
	}
}

func (stream *Stream) WriteInt8(val int8) {
	if stream.Available() < 4 {
		stream.Flush()
	}
	charPos := stream.n
	if (val < 0) {
		charPos += 1
		val = -val
		stream.buf[stream.n] = '-'
	}
	if val <= 9 {
		charPos += 1;
	} else {
		if val <= 99 {
			charPos += 2;
		} else {
			charPos += 3;
		}
	}
	stream.n = charPos
	var q int8
	for {
		q = val / 10
		r := val - ((q << 3) + (q << 1))  // r = i-(q*10) ...
		charPos--
		stream.buf[charPos] = digits[r]
		val = q;
		if val == 0 {
			break
		}
	}
}

func (stream *Stream) WriteUint16(val uint16) {
	if stream.Available() < 5 {
		stream.Flush()
	}
	charPos := stream.n
	if val <= 99 {
		if val <= 9 {
			charPos += 1;
		} else {
			charPos += 2;
		}
	} else {
		if val <= 999 {
			charPos += 3;
		} else {
			if val <= 9999 {
				charPos += 4;
			} else {
				charPos += 5;
			}
		}
	}
	stream.n = charPos
	var q uint16
	for {
		q = val / 10
		r := val - ((q << 3) + (q << 1))  // r = i-(q*10) ...
		charPos--
		stream.buf[charPos] = digits[r]
		val = q;
		if val == 0 {
			break
		}
	}
}

func (stream *Stream) WriteInt16(val int16) {
	if stream.Available() < 6 {
		stream.Flush()
	}
	charPos := stream.n
	if (val < 0) {
		charPos += 1
		val = -val
		stream.buf[stream.n] = '-'
	}
	if val <= 99 {
		if val <= 9 {
			charPos += 1;
		} else {
			charPos += 2;
		}
	} else {
		if val <= 999 {
			charPos += 3;
		} else {
			if val <= 9999 {
				charPos += 4;
			} else {
				charPos += 5;
			}
		}
	}
	stream.n = charPos
	var q int16
	for {
		q = val / 10
		r := val - ((q << 3) + (q << 1))  // r = i-(q*10) ...
		charPos--
		stream.buf[charPos] = digits[r]
		val = q;
		if val == 0 {
			break
		}
	}
}

func (stream *Stream) WriteUint32(val uint32) {
	if stream.Available() < 10 {
		stream.Flush()
	}
	charPos := stream.n
	if val <= 99999 {
		if val <= 999 {
			if val <= 9 {
				charPos += 1;
			} else {
				if val <= 99 {
					charPos += 2;
				} else {
					charPos += 3;
				}
			}
		} else {
			if val <= 9999 {
				charPos += 4;
			} else {
				charPos += 5;
			}
		}
	} else {
		if val < 99999999 {
			if val <= 999999 {
				charPos += 6;
			} else {
				if val <= 9999999 {
					charPos += 7;
				} else {
					charPos += 8;
				}
			}
		} else {
			if val <= 999999999 {
				charPos += 9;
			} else {
				charPos += 10;
			}
		}
	}
	stream.n = charPos

	var q uint32
	for val >= 65536 {
		q = val / 100;
		// really: r = i - (q * 100);
		r := val - ((q << 6) + (q << 5) + (q << 2));
		val = q;
		charPos--
		stream.buf[charPos] = digitOnes[r];
		charPos--
		stream.buf[charPos] = digitTens[r];
	}

	for {
		q = val / 10
		r := val - ((q << 3) + (q << 1))  // r = i-(q*10) ...
		charPos--
		stream.buf[charPos] = digits[r]
		val = q;
		if val == 0 {
			break
		}
	}
}

func (stream *Stream) WriteInt32(val int32) {
	if stream.Available() < 11 {
		stream.Flush()
	}
	charPos := stream.n
	if (val < 0) {
		charPos += 1
		val = -val
		stream.buf[stream.n] = '-'
	}
	if val <= 99999 {
		if val <= 999 {
			if val <= 9 {
				charPos += 1;
			} else {
				if val <= 99 {
					charPos += 2;
				} else {
					charPos += 3;
				}
			}
		} else {
			if val <= 9999 {
				charPos += 4;
			} else {
				charPos += 5;
			}
		}
	} else {
		if val < 99999999 {
			if val <= 999999 {
				charPos += 6;
			} else {
				if val <= 9999999 {
					charPos += 7;
				} else {
					charPos += 8;
				}
			}
		} else {
			if val <= 999999999 {
				charPos += 9;
			} else {
				charPos += 10;
			}
		}
	}
	stream.n = charPos

	var q int32
	for val >= 65536 {
		q = val / 100;
		// really: r = i - (q * 100);
		r := val - ((q << 6) + (q << 5) + (q << 2));
		val = q;
		charPos--
		stream.buf[charPos] = digitOnes[r];
		charPos--
		stream.buf[charPos] = digitTens[r];
	}

	for {
		q = val / 10
		r := val - ((q << 3) + (q << 1))  // r = i-(q*10) ...
		charPos--
		stream.buf[charPos] = digits[r]
		val = q;
		if val == 0 {
			break
		}
	}
}

func (stream *Stream) WriteUint64(val uint64) {
	if stream.Available() < 10 {
		stream.Flush()
	}
	charPos := stream.n
	if val <= 99999 {
		if val <= 999 {
			if val <= 9 {
				charPos += 1;
			} else {
				if val <= 99 {
					charPos += 2;
				} else {
					charPos += 3;
				}
			}
		} else {
			if val <= 9999 {
				charPos += 4;
			} else {
				charPos += 5;
			}
		}
	} else if val < 9999999999 {
		if val < 99999999 {
			if val <= 999999 {
				charPos += 6;
			} else {
				if val <= 9999999 {
					charPos += 7;
				} else {
					charPos += 8;
				}
			}
		} else {
			if val <= 999999999 {
				charPos += 9;
			} else {
				charPos += 10;
			}
		}
	} else {
		stream.writeUint64SlowPath(val)
		return
	}
	stream.n = charPos
	var q uint64
	for val >= 65536 {
		q = val / 100;
		// really: r = i - (q * 100);
		r := val - ((q << 6) + (q << 5) + (q << 2));
		val = q;
		charPos--
		stream.buf[charPos] = digitOnes[r];
		charPos--
		stream.buf[charPos] = digitTens[r];
	}

	for {
		q = val / 10
		r := val - ((q << 3) + (q << 1))  // r = i-(q*10) ...
		charPos--
		stream.buf[charPos] = digits[r]
		val = q;
		if val == 0 {
			break
		}
	}
}

func (stream *Stream) WriteInt64(val int64) {
	if (val < 0) {
		val = -val
		stream.writeByte('-')
	}
	stream.WriteUint64(uint64(val))
}

func (stream *Stream) writeUint64SlowPath(val uint64) {
	var temp [20]byte
	charPos := 20
	var q uint64
	for val >= 65536 {
		q = val / 100;
		// really: r = i - (q * 100);
		r := val - ((q << 6) + (q << 5) + (q << 2));
		val = q;
		charPos--
		temp[charPos] = digitOnes[r];
		charPos--
		temp[charPos] = digitTens[r];
	}

	for {
		q = val / 10
		r := val - ((q << 3) + (q << 1))  // r = i-(q*10) ...
		charPos--
		temp[charPos] = digits[r]
		val = q;
		if val == 0 {
			break
		}
	}
	stream.Write(temp[charPos:])
}