package cmn

const (
	White Color = iota
	Black
)

const (
	Pawn PieceType = iota
	Knight
	Bishop
	Rook
	Queen
	King
)

const Infinity = 1000000
const MateScore = 900000

const (
	FlagCapture uint8 = 1 << iota
	FlagDoublePawnPush
	FlagEnPassant
	FlagCastling
	FlagPromotion
)

var PieceCharMap = map[rune]struct {
	Type  PieceType
	Color Color
}{
	'P': {Pawn, White}, 'N': {Knight, White}, 'B': {Bishop, White},
	'R': {Rook, White}, 'Q': {Queen, White}, 'K': {King, White},
	'p': {Pawn, Black}, 'n': {Knight, Black}, 'b': {Bishop, Black},
	'r': {Rook, Black}, 'q': {Queen, Black}, 'k': {King, Black},
}

var KnightAttacks [64]Bitboard
var KingAttacks [64]Bitboard

var PieceKeys [12][64]uint64
var CastleKeys [16]uint64
var EpKeys [8]uint64
var SideKey uint64

const (
	Rank1 Bitboard = 0x00000000000000FF
	Rank2 Bitboard = 0x000000000000FF00
	Rank7 Bitboard = 0x00FF000000000000
	Rank8 Bitboard = 0xFF00000000000000

	NotAFile Bitboard = 0xfefefefefefefefe
	NotHFile Bitboard = 0x7f7f7f7f7f7f7f7f
)
