package aiff

import "errors"

var (
	formID = [4]byte{'F', 'O', 'R', 'M'}
	aiffID = [4]byte{'A', 'I', 'F', 'F'}
	aifcID = [4]byte{'A', 'I', 'F', 'C'}
	COMMID = [4]byte{'C', 'O', 'M', 'M'}
	SSNDID = [4]byte{'S', 'S', 'N', 'D'}

	// AIFC encodings
	encNone = [4]byte{'N', 'O', 'N', 'E'}
	// inverted byte order LE instead of BE (not really compression)
	encSowt = [4]byte{'s', 'o', 'w', 't'}
	// inverted byte order LE instead of BE (not really compression)
	encTwos = [4]byte{'t', 'w', 'o', 's'}
	encRaw  = [4]byte{'r', 'a', 'w', ' '}
	encIn24 = [4]byte{'i', 'n', '2', '4'}
	enc42n1 = [4]byte{'4', '2', 'n', '1'}
	encIn32 = [4]byte{'i', 'n', '3', '2'}
	enc23ni = [4]byte{'2', '3', 'n', 'i'}

	encFl32 = [4]byte{'f', 'l', '3', '2'}
	encFL32 = [4]byte{'F', 'L', '3', '2'}
	encFl64 = [4]byte{'f', 'l', '6', '4'}
	encFL64 = [4]byte{'F', 'L', '6', '4'}

	envUlaw = [4]byte{'u', 'l', 'a', 'w'}
	encULAW = [4]byte{'U', 'L', 'A', 'W'}
	encAlaw = [4]byte{'a', 'l', 'a', 'w'}
	encALAW = [4]byte{'A', 'L', 'A', 'W'}

	encDwvw = [4]byte{'D', 'W', 'V', 'W'}
	encGsm  = [4]byte{'G', 'S', 'M', ' '}
	encIma4 = [4]byte{'i', 'm', 'a', '4'}

	// ErrFmtNotSupported is a generic error reporting an unknown format.
	ErrFmtNotSupported = errors.New("format not supported")
	// ErrUnexpectedData is a generic error reporting that the parser encountered unexpected data.
	ErrUnexpectedData = errors.New("unexpected data content")
)
