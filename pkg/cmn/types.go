package cmn

type Bitboard uint64

type Color uint8

type PieceType uint8

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

func NewMove(from, to int, piece PieceType) Move {
	return Move{
		From:  from,
		To:    to,
		Piece: piece,
	}
}

func (m Move) String() string {
	s := SquareToNotation(m.From) + SquareToNotation(m.To)

	if m.Flags&FlagPromotion != 0 {
		switch m.Promotion {
		case Queen:
			s += "q"
		case Rook:
			s += "r"
		case Bishop:
			s += "b"
		case Knight:
			s += "n"
		}
	}
	return s
}
