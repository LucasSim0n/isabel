package board

type Move struct {
	From uint8
	To   uint8

	Piece     PieceType
	Promotion PieceType
	Capture   PieceType

	Flags uint8
}

const (
	FlagCapture uint8 = 1 << iota
	FlagDoublePawnPush
	FlagEnPassant
	FlagCastling
	FlagPromotion
)

func NewMove(from, to int, piece PieceType) Move {
	return Move{
		From:  uint8(from),
		To:    uint8(to),
		Piece: piece,
	}
}
