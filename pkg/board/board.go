package board

import (
	"math/bits"

	"github.com/LucasSim0n/isabel/pkg/cmn"
)

type Board struct {
	Pieces    [12]cmn.Bitboard
	Occupancy [3]cmn.Bitboard

	SideToMove     cmn.Color
	CastlingRights uint8
	EnPassant      int
	HalfmoveClock  int
	FullmoveNumber int

	KingSq [2]uint8

	Hash uint64
}

func (b *Board) updateOccupancy() {
	b.Occupancy[cmn.White] = 0
	b.Occupancy[cmn.Black] = 0

	for i := range 6 {
		b.Occupancy[cmn.White] |= b.Pieces[i]
		b.Occupancy[cmn.Black] |= b.Pieces[i+6]
	}

	b.Occupancy[2] = b.Occupancy[cmn.White] | b.Occupancy[cmn.Black]
}

func (b *Board) IsSquareAttacked(sq int, by cmn.Color) bool {

	index := 0
	if by == cmn.Black {
		index = 6
	}

	if by == cmn.White {
		attackers := ((cmn.Bitboard(1) << sq) >> 7 & cmn.NotAFile) |
			((cmn.Bitboard(1) << sq) >> 9 & cmn.NotHFile)

		if attackers&b.Pieces[cmn.Pawn] != 0 {
			return true
		}
	} else {
		attackers := ((cmn.Bitboard(1) << sq) << 7 & cmn.NotHFile) | ((cmn.Bitboard(1) << sq) << 9 & cmn.NotAFile)

		if attackers&b.Pieces[int(cmn.Pawn)+index] != 0 {
			return true
		}
	}

	knights := b.Pieces[int(cmn.Knight)+index]

	if cmn.KnightAttacks[sq]&knights != 0 {
		return true
	}

	king := b.KingSq[by]
	if cmn.KingAttacks[king]&(1<<sq) != 0 {
		return true
	}

	bishops := b.Pieces[int(cmn.Bishop)+index] | b.Pieces[int(cmn.Queen)+index]
	if b.getBishopAttacks(sq)&bishops != 0 {
		return true
	}

	rooks := b.Pieces[int(cmn.Rook)+index] | b.Pieces[int(cmn.Queen)+index]
	if b.getRookAttacks(sq)&rooks != 0 {
		return true
	}

	return false
}

func (b *Board) GetPieceAt(sq int) cmn.PieceType {
	for i := range 12 {
		if (b.Pieces[i]>>sq)&1 != 0 {
			return cmn.PieceType(i % 6)
		}
	}
	panic("No piece")
}

func (b *Board) ComputeHash() uint64 {
	var h uint64

	for p := range 12 {
		bb := b.Pieces[p]

		for bb != 0 {
			sq := bits.TrailingZeros64(uint64(bb))

			h ^= cmn.PieceKeys[p][sq]

			bb &= bb - 1
		}
	}

	if b.SideToMove == cmn.Black {
		h ^= cmn.SideKey
	}

	h ^= cmn.CastleKeys[b.CastlingRights]

	if b.EnPassant != -1 {
		file := b.EnPassant % 8
		h ^= cmn.EpKeys[file]
	}

	return h
}

func (b *Board) Clone() *Board {
	c := *b
	return &c
}
