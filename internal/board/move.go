package board

type Move struct {
	From int
	To   int

	Piece     PieceType
	Promotion PieceType
	Capture   PieceType

	Flags uint8
}

type Undo struct {
	CastlingRights uint8
	EnPassant      int
	Halfmove       int
	FullmoveNumber int
	Hash           uint64
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
		From:  from,
		To:    to,
		Piece: piece,
	}
}
