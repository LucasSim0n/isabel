package board

import (
	"math/bits"

	"github.com/LucasSim0n/isabel/pkg/cmn"
)

func (b *Board) GenerateLegalCaptures() []cmn.Move {
	legal := make([]cmn.Move, 0, 64)

	pseudo := b.generateCaptures()

	for _, move := range pseudo {
		valid := b.isLegalMove(move)
		if valid {
			legal = append(legal, move)
		}
	}

	return legal
}

func (b *Board) generateCaptures() []cmn.Move {

	captures := make([]cmn.Move, 0, 64)

	b.generateKingCaptures(&captures)
	b.generatePawnCaptures(&captures)

	b.generateCapturesByPiece(&captures, cmn.Queen)
	b.generateCapturesByPiece(&captures, cmn.Rook)
	b.generateCapturesByPiece(&captures, cmn.Bishop)
	b.generateCapturesByPiece(&captures, cmn.Knight)

	return captures
}

func (b *Board) generateKingCaptures(moves *[]cmn.Move) {
	color := b.SideToMove
	enemy := cmn.GetOpposite(color)

	from := int(b.KingSq[color])

	attacks := cmn.KingAttacks[from]
	attacks &= b.Occupancy[enemy]

	b.iterateMoves(from, attacks, cmn.King, moves)
}

func (b *Board) generatePawnCaptures(moves *[]cmn.Move) {
	color := b.SideToMove
	enemy := cmn.GetOpposite(color)
	index := cmn.Pawn

	if color == cmn.Black {
		index += 6
	}

	pawns := b.Pieces[index]

	var left, right cmn.Bitboard

	if color == cmn.White {
		left = (pawns << 7) & cmn.NotHFile
		right = (pawns << 9) & cmn.NotAFile

		leftProm := left & cmn.Rank8 & b.Occupancy[enemy]
		rightProm := right & cmn.Rank8 & b.Occupancy[enemy]

		leftNorm := left & ^cmn.Rank8 & b.Occupancy[enemy]
		rightNorm := right & ^cmn.Rank8 & b.Occupancy[enemy]

		b.iteratePawnCaptures(leftNorm, 7, moves)
		b.iteratePawnCaptures(rightNorm, 9, moves)

		b.iteratePawnPromotions(leftProm, 7, true, moves)
		b.iteratePawnPromotions(rightProm, 9, true, moves)

		if b.EnPassant != -1 {
			epBB := cmn.Bitboard(1) << b.EnPassant

			leftEP := left & epBB
			rightEP := right & epBB

			b.iterateEnPassant(leftEP, 7, moves)
			b.iterateEnPassant(rightEP, 9, moves)
		}

	} else {
		left = (pawns >> 9) & cmn.NotHFile
		right = (pawns >> 7) & cmn.NotAFile

		leftProm := left & cmn.Rank1 & b.Occupancy[enemy]
		rightProm := right & cmn.Rank1 & b.Occupancy[enemy]

		leftNorm := left & ^cmn.Rank1 & b.Occupancy[enemy]
		rightNorm := right & ^cmn.Rank1 & b.Occupancy[enemy]

		b.iteratePawnCaptures(leftNorm, -9, moves)
		b.iteratePawnCaptures(rightNorm, -7, moves)

		b.iteratePawnPromotions(leftProm, -9, true, moves)
		b.iteratePawnPromotions(rightProm, -7, true, moves)

		if b.EnPassant != -1 {
			epBB := cmn.Bitboard(1) << b.EnPassant

			leftEP := left & epBB
			rightEP := right & epBB

			b.iterateEnPassant(leftEP, -9, moves)
			b.iterateEnPassant(rightEP, -7, moves)
		}
	}
}

func (b *Board) generateCapturesByPiece(moves *[]cmn.Move, piece cmn.PieceType) {
	color := b.SideToMove
	enemy := cmn.GetOpposite(color)

	index := piece
	if color == cmn.Black {
		index += 6
	}

	bb := b.Pieces[index]

	for bb != 0 {
		from := bits.TrailingZeros64(uint64(bb))

		var attacks cmn.Bitboard

		switch piece {
		case cmn.Bishop:
			attacks = b.getBishopAttacks(from)
		case cmn.Rook:
			attacks = b.getRookAttacks(from)
		case cmn.Queen:
			attacks = b.getBishopAttacks(from) | b.getRookAttacks(from)
		case cmn.Knight:
			attacks = cmn.KnightAttacks[from]
		}

		attacks &= b.Occupancy[enemy]

		b.iterateMoves(from, attacks, piece, moves)

		bb &= bb - 1
	}
}
