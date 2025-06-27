package mpls

func isNumeric(buf [5]byte) bool {
	if len(buf) == 0 {
		return false
	}
	for _, b := range buf {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

func isAlphanumericUppercase(buf [4]byte) bool {
	for _, b := range buf {
		if !('A' <= b && b <= 'Z' || '0' <= b && b <= '9') {
			return false
		}
	}
	return true
}

//func isLowerCase(buf []byte) bool {
//	for _, b := range buf {
//		if b < 'a' || b > 'z' {
//			return false
//		}
//	}
//	return true
//}
//
//func isUpperCase(buf [4]byte) bool {
//	for _, b := range buf {
//		if b < 'A' || b > 'Z' {
//			return false
//		}
//	}
//	return true
//}
