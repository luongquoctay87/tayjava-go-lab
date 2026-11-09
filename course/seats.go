package course

import "errors"

// ErrInvalidEnrollment is returned when enrolled exceeds capacity.
var ErrInvalidEnrollment = errors.New("course: invalid enrollment")

// RemainingSeats returns capacity minus enrolled, or an error.
func RemainingSeats(capacity, enrolled int) (int, error) {
	if capacity < 0 || enrolled < 0 {
		return 0, errors.New("negative seats")
	}
	if enrolled > capacity {
		return 0, ErrInvalidEnrollment
	}
	return capacity - enrolled, nil
}
