package numbers

import "strconv"

func FormatNumber(n int64) string {
	// Ubah angka ke string
	str := strconv.FormatInt(n, 10)

	// Hitung berapa digit
	length := len(str)

	if length <= 3 {
		return str
	}

	// Mulai dari belakang, tambahkan titik setiap 3 digit
	var result []byte
	for i := length - 1; i >= 0; i-- {
		if (length-i-1)%3 == 0 && i != length-1 {
			result = append([]byte{'.'}, result...)
		}
		result = append([]byte{str[i]}, result...)
	}

	return string(result)
}
