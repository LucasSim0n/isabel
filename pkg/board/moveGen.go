package board

import (
	"math/bits"

	"github.com/LucasSim0n/isabel/pkg/cmn"
)

func (b *Board) GenerateMoves() []cmn.Move {
	moves := make([]cmn.Move, 0, 64)

	b.generateKingMoves(&moves)
	b.generateKnightMoves(&moves)
	b.generatePawnMoves(&moves)
	b.generateBishopMoves(&moves)
	b.generateRookMoves(&moves)
	b.generateQueenMoves(&moves)

	return moves
}

func (b *Board) GenerateLegalMoves() []cmn.Move {
	legal := make([]cmn.Move, 0, 64)

	pseudo := b.GenerateMoves()

	for _, move := range pseudo {
		undo := b.MakeMove(move)

		movedSide := cmn.GetOpposite(b.SideToMove)

		inCheck := b.IsSquareAttacked(int(b.KingSq[movedSide]), b.SideToMove)

		b.UnmakeMove(move, undo)

		if !inCheck {
			legal = append(legal, move)
		}
	}

	return legal
}

func (b *Board) getBishopAttacks(sq int) cmn.Bitboard {
	var attacks cmn.Bitboard

	rank := sq / 8
	file := sq % 8

	for r, f := rank+1, file+1; r <= 7 && f <= 7; r, f = r+1, f+1 {
		sq := r*8 + f
		attacks |= 1 << sq

		if (b.Occupancy[2]>>sq)&1 != 0 {
			break
		}
	}

	for r, f := rank+1, file-1; r <= 7 && f >= 0; r, f = r+1, f-1 {
		sq := r*8 + f
		attacks |= 1 << sq

		if (b.Occupancy[2]>>sq)&1 != 0 {
			break
		}
	}

	for r, f := rank-1, file-1; r >= 0 && f >= 0; r, f = r-1, f-1 {
		sq := r*8 + f
		attacks |= 1 << sq

		if (b.Occupancy[2]>>sq)&1 != 0 {
			break
		}
	}

	for r, f := rank-1, file+1; r >= 0 && f <= 7; r, f = r-1, f+1 {
		sq := r*8 + f
		attacks |= 1 << sq

		if (b.Occupancy[2]>>sq)&1 != 0 {
			break
		}
	}

	return attacks
}

func (b *Board) getRookAttacks(sq int) cmn.Bitboard {
	var attacks cmn.Bitboard

	rank := sq / 8
	file := sq % 8

	for r := rank + 1; r <= 7; r++ {
		sq := r*8 + file
		attacks |= 1 << sq

		if (b.Occupancy[2]>>sq)&1 != 0 {
			break
		}
	}

	for r := rank - 1; r >= 0; r-- {
		sq := r*8 + file
		attacks |= 1 << sq

		if (b.Occupancy[2]>>sq)&1 != 0 {
			break
		}
	}

	for f := file + 1; f <= 7; f++ {
		sq := rank*8 + f
		attacks |= 1 << sq

		if (b.Occupancy[2]>>sq)&1 != 0 {
			break
		}
	}

	for f := file - 1; f >= 0; f-- {
		sq := rank*8 + f
		attacks |= 1 << sq

		if (b.Occupancy[2]>>sq)&1 != 0 {
			break
		}
	}

	return attacks
}

func (b *Board) generateKingMoves(moves *[]cmn.Move) {
	color := b.SideToMove

	from := int(b.KingSq[color])
	attacks := cmn.KingAttacks[from]
	attacks &= ^b.Occupancy[color]

	b.iterateMoves(from, attacks, cmn.King, moves)

	b.generateCastlingMoves(moves)
}

func (b *Board) generateCastlingMoves(moves *[]cmn.Move) {
	color := b.SideToMove

	if color == cmn.White {
		if !b.IsSquareAttacked(4, cmn.Black) {
			if b.CastlingRights&0b1000 != 0 {
				if ((b.Occupancy[2]>>5)&1 == 0) && ((b.Occupancy[2]>>6)&1 == 0) &&
					!b.IsSquareAttacked(5, cmn.Black) &&
					!b.IsSquareAttacked(6, cmn.Black) {
					*moves = append(*moves, cmn.Move{
						From:  4,
						To:    6,
						Piece: cmn.King,
						Flags: cmn.FlagCastling,
					})
				}
			}

			if b.CastlingRights&0b0100 != 0 {
				if ((b.Occupancy[2]>>3)&1 == 0) &&
					((b.Occupancy[2]>>2)&1 == 0) &&
					((b.Occupancy[2]>>1)&1 == 0) &&
					!b.IsSquareAttacked(3, cmn.Black) &&
					!b.IsSquareAttacked(2, cmn.Black) {
					*moves = append(*moves, cmn.Move{
						From:  4,
						To:    2,
						Piece: cmn.King,
						Flags: cmn.FlagCastling,
					})
				}
			}
		}
	} else if !b.IsSquareAttacked(60, cmn.White) {
		if b.CastlingRights&0b0010 != 0 {
			if ((b.Occupancy[2]>>61)&1 == 0) && ((b.Occupancy[2]>>62)&1 == 0) &&
				!b.IsSquareAttacked(61, cmn.White) &&
				!b.IsSquareAttacked(62, cmn.White) {
				*moves = append(*moves, cmn.Move{
					From:  60,
					To:    62,
					Piece: cmn.King,
					Flags: cmn.FlagCastling,
				})
			}
		}

		if b.CastlingRights&0b0001 != 0 {
			if ((b.Occupancy[2]>>59)&1 == 0) &&
				((b.Occupancy[2]>>58)&1 == 0) &&
				((b.Occupancy[2]>>57)&1 == 0) &&
				!b.IsSquareAttacked(59, cmn.White) &&
				!b.IsSquareAttacked(58, cmn.White) {
				*moves = append(*moves, cmn.Move{
					From:  60,
					To:    58,
					Piece: cmn.King,
					Flags: cmn.FlagCastling,
				})
			}
		}
	}
}

func (b *Board) generateKnightMoves(moves *[]cmn.Move) {
	color := b.SideToMove
	index := int(cmn.Knight)
	if color == cmn.Black {
		index += 6
	}

	bb := b.Pieces[index]
	for bb != 0 {
		from := bits.TrailingZeros64(uint64(bb))
		attacks := cmn.KnightAttacks[from]
		attacks &= ^b.Occupancy[color]

		b.iterateMoves(from, attacks, cmn.Knight, moves)

		bb &= bb - 1
	}
}

func (b *Board) generateBishopMoves(moves *[]cmn.Move) {
	color := b.SideToMove

	index := cmn.Bishop
	if color == cmn.Black {
		index += 6
	}

	bb := b.Pieces[index]

	for bb != 0 {
		from := bits.TrailingZeros64(uint64(bb))

		attacks := b.getBishopAttacks(from)
		attacks &= ^b.Occupancy[color]

		b.iterateMoves(from, attacks, cmn.Bishop, moves)

		bb &= bb - 1
	}
}

func (b *Board) generateRookMoves(moves *[]cmn.Move) {
	color := b.SideToMove

	index := cmn.Rook
	if color == cmn.Black {
		index += 6
	}

	bb := b.Pieces[index]

	for bb != 0 {
		from := bits.TrailingZeros64(uint64(bb))

		attacks := b.getRookAttacks(from)
		attacks &= ^b.Occupancy[color]

		b.iterateMoves(from, attacks, cmn.Rook, moves)

		bb &= bb - 1
	}
}

func (b *Board) generateQueenMoves(moves *[]cmn.Move) {
	color := b.SideToMove

	index := cmn.Queen
	if color == cmn.Black {
		index += 6
	}

	bb := b.Pieces[index]

	for bb != 0 {
		from := bits.TrailingZeros64(uint64(bb))

		attacks := b.getRookAttacks(from) | b.getBishopAttacks(from)
		attacks &= ^b.Occupancy[color]

		b.iterateMoves(from, attacks, cmn.Queen, moves)

		bb &= bb - 1
	}
}

func (b *Board) generatePawnMoves(moves *[]cmn.Move) {
	color := b.SideToMove
	enemy := cmn.GetOpposite(color)
	forward := 8
	index := cmn.Pawn

	if color == cmn.Black {
		index += 6
		forward = -8
	}

	pawns := b.Pieces[index]

	empty := ^b.Occupancy[2]

	var singlePush cmn.Bitboard
	if color == cmn.White {
		singlePush = (pawns << 8) & empty

		promotions := singlePush & cmn.Rank8
		normal := singlePush & ^cmn.Rank8
		b.iteratePawnMoves(normal, forward, moves)
		b.iteratePawnPromotions(promotions, forward, false, moves)

	} else {
		singlePush = (pawns >> 8) & empty

		promotions := singlePush & cmn.Rank1
		normal := singlePush & ^cmn.Rank1

		b.iteratePawnMoves(normal, forward, moves)
		b.iteratePawnPromotions(promotions, forward, false, moves)
	}

	if color == cmn.White {
		p := pawns & cmn.Rank2
		s1 := (p << 8) & empty
		doublePush := (s1 << 8) & empty

		b.iteratePawnMoves(doublePush, forward*2, moves)
	} else {
		p := pawns & cmn.Rank7
		s1 := (p >> 8) & empty
		doublePush := (s1 >> 8) & empty

		b.iteratePawnMoves(doublePush, forward*2, moves)
	}

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

func (b *Board) iteratePawnPromotions(bb cmn.Bitboard, offset int, isCapture bool, moves *[]cmn.Move) {
	for bb != 0 {
		to := bits.TrailingZeros64(uint64(bb))
		from := to - offset

		promotions := []cmn.PieceType{cmn.Queen, cmn.Rook, cmn.Bishop, cmn.Knight}

		for _, promo := range promotions {
			move := cmn.NewMove(from, to, cmn.Pawn)
			move.Flags |= cmn.FlagPromotion
			move.Promotion = promo

			if isCapture {
				move.Flags |= cmn.FlagCapture
				move.Capture = b.GetPieceAt(to)
			}

			*moves = append(*moves, move)
		}

		bb &= bb - 1
	}
}

func (b *Board) iteratePawnMoves(bb cmn.Bitboard, offset int, moves *[]cmn.Move) {
	for bb != 0 {
		to := bits.TrailingZeros64(uint64(bb))
		from := to - offset

		move := cmn.NewMove(from, to, cmn.Pawn)

		if offset == 16 || offset == -16 {
			move.Flags |= cmn.FlagDoublePawnPush
		}
		*moves = append(*moves, move)

		bb &= bb - 1
	}

}

func (b *Board) iteratePawnCaptures(bb cmn.Bitboard, offset int, moves *[]cmn.Move) {
	for bb != 0 {
		to := bits.TrailingZeros64(uint64(bb))
		from := to - offset

		move := cmn.NewMove(from, to, cmn.Pawn)
		move.Flags |= cmn.FlagCapture
		move.Capture = b.GetPieceAt(to)

		*moves = append(*moves, move)

		bb &= bb - 1
	}
}

func (b *Board) iterateEnPassant(bb cmn.Bitboard, offset int, moves *[]cmn.Move) {
	for bb != 0 {
		to := bits.TrailingZeros64(uint64(bb))
		from := to - offset

		move := cmn.NewMove(from, to, cmn.Pawn)
		move.Flags |= cmn.FlagEnPassant

		*moves = append(*moves, move)

		bb &= bb - 1
	}
}

func (b *Board) iterateMoves(from int, attacks cmn.Bitboard, piece cmn.PieceType, moves *[]cmn.Move) {

	enemy := cmn.GetOpposite(b.SideToMove)

	for attacks != 0 {
		to := bits.TrailingZeros64(uint64(attacks))

		move := cmn.NewMove(from, to, piece)

		if (b.Occupancy[enemy]>>to)&1 != 0 {
			move.Flags |= cmn.FlagCapture

			move.Capture = b.GetPieceAt(to)
		}

		*moves = append(*moves, move)

		attacks &= attacks - 1
	}
}
