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

func (m Move) String() string {
	s := squareToNotation(m.From) + squareToNotation(m.To)

	switch m.Piece {
	case Queen:
		s = "q " + s
	case Rook:
		s = "r " + s
	case Bishop:
		s = "b " + s
	case Knight:
		s = "n " + s
	case King:
		s = "k " + s
	}

	if m.Flags&FlagPromotion != 0 {
		switch m.Promotion {
		case Queen:
			s += "=q"
		case Rook:
			s += "=r"
		case Bishop:
			s += "=b"
		case Knight:
			s += "=n"
		}
	}
	return s
}
