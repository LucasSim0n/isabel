package board

import (
	"fmt"
	"math/bits"
	"strings"
)

type Board struct {
	Pieces    [12]Bitboard
	Occupancy [3]Bitboard

	SideToMove     Color
	CastlingRights uint8
	EnPassant      int
	HalfmoveClock  int
	FullmoveNumber int

	KingSq [2]uint8

	Hash uint64
}

func (b *Board) Evaluate() int {
	score := 0

	for p := Pawn; p <= Queen; p++ {

		whiteCount := bits.OnesCount64(uint64(b.Pieces[p]))

		blackCount := bits.OnesCount64(uint64(b.Pieces[p+6]))

		score += whiteCount * pieceValues[p]
		score -= blackCount * pieceValues[p]
	}

	if b.SideToMove == Black {
		score = -score
	}

	return score
}

func (b *Board) Search(depth, alpha, beta int) int {
	if depth == 0 {
		return b.Evaluate()
	}

	moves := b.GenerateLegalMoves()

	if len(moves) == 0 {
		inCheck := b.IsSquareAttacked(int(b.KingSq[b.SideToMove]), getOpposite(b.SideToMove))

		if inCheck {
			return -MateScore
		}

		return 0
	}

	bestScore := -Infinity

	for _, move := range moves {
		undo := b.MakeMove(move)

		score := -b.Search(
			depth-1,
			-beta,
			-alpha,
		)

		b.UnmakeMove(move, undo)

		if score > bestScore {
			bestScore = score
		}

		if score > alpha {
			alpha = score
		}

		if alpha >= beta {
			break
		}
	}

	return bestScore
}

func (b *Board) FindBestMove(depth int) Move {
	moves := b.GenerateLegalMoves()

	bestScore := -Infinity
	bestMove := moves[0]

	alpha := -Infinity
	beta := Infinity

	for _, move := range moves {

		undo := b.MakeMove(move)

		score := -b.Search(
			depth-1,
			-beta,
			-alpha,
		)

		b.UnmakeMove(move, undo)

		if score > bestScore {
			bestScore = score
			bestMove = move
		}

		if score > alpha {
			alpha = score
		}
	}

	return bestMove
}

func (b *Board) updateOccupancy() {
	b.Occupancy[White] = 0
	b.Occupancy[Black] = 0

	for i := range 6 {
		b.Occupancy[White] |= b.Pieces[i]
		b.Occupancy[Black] |= b.Pieces[i+6]
	}

	b.Occupancy[2] = b.Occupancy[White] | b.Occupancy[Black]
}

func (b *Board) MakeMove(move Move) Undo {

	undo := Undo{
		CastlingRights: b.CastlingRights,
		EnPassant:      b.EnPassant,
		Halfmove:       b.HalfmoveClock,
		FullmoveNumber: b.FullmoveNumber,
		Hash:           b.Hash,
	}

	color := b.SideToMove
	enemy := getOpposite(b.SideToMove)

	if b.EnPassant != -1 {
		b.Hash ^= epKeys[b.EnPassant%8]
	}

	b.EnPassant = -1

	if move.Flags&FlagCapture != 0 {
		b.removePiece(move.Capture, enemy, move.To)
		b.HalfmoveClock = 0

		switch move.To {
		case 0:
			b.clearCastle(0b1011)
		case 7:
			b.clearCastle(0b0111)
		case 56:
			b.clearCastle(0b1110)
		case 63:
			b.clearCastle(0b1101)
		}
	}

	if move.Flags&FlagEnPassant != 0 {
		b.handleEnPassant(move, color)
		b.HalfmoveClock = 0
	}

	b.movePiece(move.Piece, color, move.From, move.To)

	if move.Flags&FlagCastling != 0 {
		b.handleCastling(move.To, color)
	}

	if move.Piece == King {
		b.KingSq[color] = uint8(move.To)

		if color == White {
			b.clearCastle(0b0011)
		} else {
			b.clearCastle(0b1100)
		}
	}

	if move.Piece == Rook {
		switch move.From {
		case 0:
			b.clearCastle(0b1011)
		case 7:
			b.clearCastle(0b0111)
		case 56:
			b.clearCastle(0b1110)
		case 63:
			b.clearCastle(0b1101)
		}
	}

	if move.Flags&FlagDoublePawnPush != 0 {
		if color == White {
			b.EnPassant = move.To - 8
		} else {
			b.EnPassant = move.To + 8
		}

		b.Hash ^= epKeys[b.EnPassant%8]
	}

	if move.Piece == Pawn {
		b.HalfmoveClock = 0
	} else {
		b.HalfmoveClock++
	}

	if move.Flags&FlagPromotion != 0 {
		b.handlePromotion(move, color)
	}

	if color == Black {
		b.FullmoveNumber++
	}

	b.SideToMove = enemy
	b.Hash ^= sideKey

	return undo
}

func (b *Board) movePiece(piece PieceType, color Color, from, to int) {
	b.removePiece(piece, color, from)
	b.addPiece(piece, color, to)
}

func (b *Board) addPiece(piece PieceType, color Color, sq int) {
	index := int(piece)
	if color == Black {
		index += 6
	}

	b.Pieces[index] |= (1 << sq)
	b.Occupancy[color] |= (1 << sq)
	b.Occupancy[2] |= (1 << sq)
	b.Hash ^= pieceKeys[index][sq]
}

func (b *Board) removePiece(piece PieceType, color Color, sq int) {
	index := int(piece)
	if color == Black {
		index += 6
	}

	b.Pieces[index] &^= (1 << sq)
	b.Occupancy[color] &^= (1 << sq)
	b.Occupancy[2] &^= (1 << sq)
	b.Hash ^= pieceKeys[index][sq]
}

func (b *Board) handleCastling(sq int, c Color) {
	switch sq {
	case 2:
		b.movePiece(Rook, c, 0, 3)
	case 6:
		b.movePiece(Rook, c, 7, 5)
	case 58:
		b.movePiece(Rook, c, 56, 59)
	case 62:
		b.movePiece(Rook, c, 63, 61)
	}
}

func (b *Board) undoCastling(sq int, c Color) {
	switch sq {
	case 2:
		b.movePiece(Rook, c, 3, 0)
	case 6:
		b.movePiece(Rook, c, 5, 7)
	case 58:
		b.movePiece(Rook, c, 59, 56)
	case 62:
		b.movePiece(Rook, c, 61, 63)
	}
}

func (b *Board) handlePromotion(move Move, color Color) {
	b.removePiece(Pawn, color, move.To)
	b.addPiece(move.Promotion, color, move.To)
}

func (b *Board) undoPromotion(move Move, color Color) {
	b.removePiece(move.Promotion, color, move.To)
	b.addPiece(Pawn, b.SideToMove, move.To)
}

func (b *Board) handleEnPassant(move Move, color Color) {
	var captureSq int

	if color == White {
		captureSq = move.To - 8
	} else {
		captureSq = move.To + 8
	}

	b.removePiece(Pawn, getOpposite(color), captureSq)
}

func (b *Board) undoEnPassant(move Move, enemy Color) {
	var captureSq int

	if b.SideToMove == White {
		captureSq = move.To - 8
	} else {
		captureSq = move.To + 8
	}

	b.addPiece(Pawn, enemy, captureSq)
}

func (b *Board) UnmakeMove(move Move, undo Undo) {
	b.HalfmoveClock = undo.Halfmove
	b.FullmoveNumber = undo.FullmoveNumber

	b.Hash ^= castleKeys[b.CastlingRights]
	b.CastlingRights = undo.CastlingRights
	b.Hash ^= castleKeys[b.CastlingRights]

	if b.EnPassant != -1 {
		b.Hash ^= epKeys[b.EnPassant%8]
	}
	b.EnPassant = undo.EnPassant
	if b.EnPassant != -1 {
		b.Hash ^= epKeys[b.EnPassant%8]
	}

	b.SideToMove = getOpposite(b.SideToMove)
	b.Hash ^= sideKey

	color := b.SideToMove
	enemy := getOpposite(color)

	if move.Flags&FlagPromotion != 0 {
		b.undoPromotion(move, color)
	}

	b.movePiece(move.Piece, color, move.To, move.From)

	if move.Flags&FlagEnPassant != 0 {
		b.undoEnPassant(move, enemy)
	}

	if move.Piece == King {
		b.KingSq[color] = uint8(move.From)
	}

	if move.Flags&FlagCapture != 0 {
		b.addPiece(move.Capture, enemy, move.To)
	}

	if move.Flags&FlagCastling != 0 {
		b.undoCastling(move.To, color)
	}
}

func (b *Board) clearCastle(mask uint8) {
	b.Hash ^= castleKeys[b.CastlingRights]

	b.CastlingRights &= mask

	b.Hash ^= castleKeys[b.CastlingRights]
}

func (b *Board) GenerateMoves() []Move {
	moves := make([]Move, 0, 64)

	b.generateKingMoves(&moves)
	b.generateKnightMoves(&moves)
	b.generatePawnMoves(&moves)
	b.generateBishopMoves(&moves)
	b.generateRookMoves(&moves)
	b.generateQueenMoves(&moves)

	return moves
}

func (b *Board) GenerateLegalMoves() []Move {
	legal := make([]Move, 0, 64)

	pseudo := b.GenerateMoves()

	for _, move := range pseudo {
		undo := b.MakeMove(move)

		movedSide := getOpposite(b.SideToMove)

		inCheck := b.IsSquareAttacked(int(b.KingSq[movedSide]), b.SideToMove)

		b.UnmakeMove(move, undo)

		if !inCheck {
			legal = append(legal, move)
		}
	}

	return legal
}

func (b *Board) IsSquareAttacked(sq int, by Color) bool {

	index := 0
	if by == Black {
		index = 6
	}

	if by == White {
		attackers := ((Bitboard(1) << sq) >> 7 & notAFile) |
			((Bitboard(1) << sq) >> 9 & notHFile)

		if attackers&b.Pieces[Pawn] != 0 {
			return true
		}
	} else {
		attackers := ((Bitboard(1) << sq) << 7 & notHFile) | ((Bitboard(1) << sq) << 9 & notAFile)

		if attackers&b.Pieces[int(Pawn)+index] != 0 {
			return true
		}
	}

	knights := b.Pieces[int(Knight)+index]

	if knightAttacks[sq]&knights != 0 {
		return true
	}

	king := b.KingSq[by]
	if kingAttacks[king]&(1<<sq) != 0 {
		return true
	}

	bishops := b.Pieces[int(Bishop)+index] | b.Pieces[int(Queen)+index]
	if b.getBishopAttacks(sq)&bishops != 0 {
		return true
	}

	rooks := b.Pieces[int(Rook)+index] | b.Pieces[int(Queen)+index]
	if b.getRookAttacks(sq)&rooks != 0 {
		return true
	}

	return false
}

func (b *Board) getBishopAttacks(sq int) Bitboard {
	var attacks Bitboard

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

func (b *Board) getRookAttacks(sq int) Bitboard {
	var attacks Bitboard

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

func (b *Board) generateKingMoves(moves *[]Move) {
	color := b.SideToMove

	from := int(b.KingSq[color])
	attacks := kingAttacks[from]
	attacks &= ^b.Occupancy[color]

	b.iterateMoves(from, attacks, King, moves)

	b.generateCastlingMoves(moves)
}

func (b *Board) generateCastlingMoves(moves *[]Move) {
	color := b.SideToMove

	if color == White {
		if !b.IsSquareAttacked(4, Black) {
			if b.CastlingRights&0b1000 != 0 {
				if ((b.Occupancy[2]>>5)&1 == 0) && ((b.Occupancy[2]>>6)&1 == 0) &&
					!b.IsSquareAttacked(5, Black) &&
					!b.IsSquareAttacked(6, Black) {
					*moves = append(*moves, Move{
						From:  4,
						To:    6,
						Piece: King,
						Flags: FlagCastling,
					})
				}
			}

			if b.CastlingRights&0b0100 != 0 {
				if ((b.Occupancy[2]>>3)&1 == 0) &&
					((b.Occupancy[2]>>2)&1 == 0) &&
					((b.Occupancy[2]>>1)&1 == 0) &&
					!b.IsSquareAttacked(3, Black) &&
					!b.IsSquareAttacked(2, Black) {
					*moves = append(*moves, Move{
						From:  4,
						To:    2,
						Piece: King,
						Flags: FlagCastling,
					})
				}
			}
		}
	} else if !b.IsSquareAttacked(60, White) {
		if b.CastlingRights&0b0010 != 0 {
			if ((b.Occupancy[2]>>61)&1 == 0) && ((b.Occupancy[2]>>62)&1 == 0) &&
				!b.IsSquareAttacked(61, White) &&
				!b.IsSquareAttacked(62, White) {
				*moves = append(*moves, Move{
					From:  60,
					To:    62,
					Piece: King,
					Flags: FlagCastling,
				})
			}
		}

		if b.CastlingRights&0b0001 != 0 {
			if ((b.Occupancy[2]>>59)&1 == 0) &&
				((b.Occupancy[2]>>58)&1 == 0) &&
				((b.Occupancy[2]>>57)&1 == 0) &&
				!b.IsSquareAttacked(59, White) &&
				!b.IsSquareAttacked(58, White) {
				*moves = append(*moves, Move{
					From:  60,
					To:    58,
					Piece: King,
					Flags: FlagCastling,
				})
			}
		}
	}
}

func (b *Board) generateKnightMoves(moves *[]Move) {
	color := b.SideToMove
	index := int(Knight)
	if color == Black {
		index += 6
	}

	bb := b.Pieces[index]
	for bb != 0 {
		from := bits.TrailingZeros64(uint64(bb))
		attacks := knightAttacks[from]
		attacks &= ^b.Occupancy[color]

		b.iterateMoves(from, attacks, Knight, moves)

		bb &= bb - 1
	}
}

func (b *Board) generateBishopMoves(moves *[]Move) {
	color := b.SideToMove

	index := Bishop
	if color == Black {
		index += 6
	}

	bb := b.Pieces[index]

	for bb != 0 {
		from := bits.TrailingZeros64(uint64(bb))

		attacks := b.getBishopAttacks(from)
		attacks &= ^b.Occupancy[color]

		b.iterateMoves(from, attacks, Bishop, moves)

		bb &= bb - 1
	}
}

func (b *Board) generateRookMoves(moves *[]Move) {
	color := b.SideToMove

	index := Rook
	if color == Black {
		index += 6
	}

	bb := b.Pieces[index]

	for bb != 0 {
		from := bits.TrailingZeros64(uint64(bb))

		attacks := b.getRookAttacks(from)
		attacks &= ^b.Occupancy[color]

		b.iterateMoves(from, attacks, Rook, moves)

		bb &= bb - 1
	}
}

func (b *Board) generateQueenMoves(moves *[]Move) {
	color := b.SideToMove

	index := Queen
	if color == Black {
		index += 6
	}

	bb := b.Pieces[index]

	for bb != 0 {
		from := bits.TrailingZeros64(uint64(bb))

		attacks := b.getRookAttacks(from) | b.getBishopAttacks(from)
		attacks &= ^b.Occupancy[color]

		b.iterateMoves(from, attacks, Queen, moves)

		bb &= bb - 1
	}
}

func (b *Board) generatePawnMoves(moves *[]Move) {
	color := b.SideToMove
	enemy := getOpposite(color)
	forward := 8
	index := Pawn

	if color == Black {
		index += 6
		forward = -8
	}

	pawns := b.Pieces[index]

	empty := ^b.Occupancy[2]

	var singlePush Bitboard
	if color == White {
		singlePush = (pawns << 8) & empty

		promotions := singlePush & rank8
		normal := singlePush & ^rank8
		b.iteratePawnMoves(normal, forward, moves)
		b.iteratePawnPromotions(promotions, forward, false, moves)

	} else {
		singlePush = (pawns >> 8) & empty

		promotions := singlePush & rank1
		normal := singlePush & ^rank1

		b.iteratePawnMoves(normal, forward, moves)
		b.iteratePawnPromotions(promotions, forward, false, moves)
	}

	if color == White {
		p := pawns & rank2
		s1 := (p << 8) & empty
		doublePush := (s1 << 8) & empty

		b.iteratePawnMoves(doublePush, forward*2, moves)
	} else {
		p := pawns & rank7
		s1 := (p >> 8) & empty
		doublePush := (s1 >> 8) & empty

		b.iteratePawnMoves(doublePush, forward*2, moves)
	}

	var left, right Bitboard

	if color == White {
		left = (pawns << 7) & notHFile
		right = (pawns << 9) & notAFile

		leftProm := left & rank8 & b.Occupancy[enemy]
		rightProm := right & rank8 & b.Occupancy[enemy]

		leftNorm := left & ^rank8 & b.Occupancy[enemy]
		rightNorm := right & ^rank8 & b.Occupancy[enemy]

		b.iteratePawnCaptures(leftNorm, 7, moves)
		b.iteratePawnCaptures(rightNorm, 9, moves)

		b.iteratePawnPromotions(leftProm, 7, true, moves)
		b.iteratePawnPromotions(rightProm, 9, true, moves)

		if b.EnPassant != -1 {
			epBB := Bitboard(1) << b.EnPassant

			leftEP := left & epBB
			rightEP := right & epBB

			b.iterateEnPassant(leftEP, 7, moves)
			b.iterateEnPassant(rightEP, 9, moves)
		}

	} else {
		left = (pawns >> 9) & notHFile
		right = (pawns >> 7) & notAFile

		leftProm := left & rank1 & b.Occupancy[enemy]
		rightProm := right & rank1 & b.Occupancy[enemy]

		leftNorm := left & ^rank1 & b.Occupancy[enemy]
		rightNorm := right & ^rank1 & b.Occupancy[enemy]

		b.iteratePawnCaptures(leftNorm, -9, moves)
		b.iteratePawnCaptures(rightNorm, -7, moves)

		b.iteratePawnPromotions(leftProm, -9, true, moves)
		b.iteratePawnPromotions(rightProm, -7, true, moves)

		if b.EnPassant != -1 {
			epBB := Bitboard(1) << b.EnPassant

			leftEP := left & epBB
			rightEP := right & epBB

			b.iterateEnPassant(leftEP, -9, moves)
			b.iterateEnPassant(rightEP, -7, moves)
		}
	}
}

func (b *Board) iteratePawnPromotions(bb Bitboard, offset int, isCapture bool, moves *[]Move) {
	for bb != 0 {
		to := bits.TrailingZeros64(uint64(bb))
		from := to - offset

		promotions := []PieceType{Queen, Rook, Bishop, Knight}

		for _, promo := range promotions {
			move := NewMove(from, to, Pawn)
			move.Flags |= FlagPromotion
			move.Promotion = promo

			if isCapture {
				move.Flags |= FlagCapture
				move.Capture = b.getPieceAt(to)
			}

			*moves = append(*moves, move)
		}

		bb &= bb - 1
	}
}

func (b *Board) iteratePawnMoves(bb Bitboard, offset int, moves *[]Move) {
	for bb != 0 {
		to := bits.TrailingZeros64(uint64(bb))
		from := to - offset

		move := NewMove(from, to, Pawn)

		if offset == 16 || offset == -16 {
			move.Flags |= FlagDoublePawnPush
		}
		*moves = append(*moves, move)

		bb &= bb - 1
	}

}

func (b *Board) iteratePawnCaptures(bb Bitboard, offset int, moves *[]Move) {
	for bb != 0 {
		to := bits.TrailingZeros64(uint64(bb))
		from := to - offset

		move := NewMove(from, to, Pawn)
		move.Flags |= FlagCapture
		move.Capture = b.getPieceAt(to)

		*moves = append(*moves, move)

		bb &= bb - 1
	}
}

func (b *Board) iterateEnPassant(bb Bitboard, offset int, moves *[]Move) {
	for bb != 0 {
		to := bits.TrailingZeros64(uint64(bb))
		from := to - offset

		move := NewMove(from, to, Pawn)
		move.Flags |= FlagEnPassant

		*moves = append(*moves, move)

		bb &= bb - 1
	}
}

func (b *Board) iterateMoves(from int, attacks Bitboard, piece PieceType, moves *[]Move) {

	enemy := getOpposite(b.SideToMove)

	for attacks != 0 {
		to := bits.TrailingZeros64(uint64(attacks))

		move := NewMove(from, to, piece)

		if (b.Occupancy[enemy]>>to)&1 != 0 {
			move.Flags |= FlagCapture

			move.Capture = b.getPieceAt(to)
		}

		*moves = append(*moves, move)

		attacks &= attacks - 1
	}
}

func (b *Board) getPieceAt(sq int) PieceType {
	for i := range 12 {
		if (b.Pieces[i]>>sq)&1 != 0 {
			return PieceType(i % 6)
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

			h ^= pieceKeys[p][sq]

			bb &= bb - 1
		}
	}

	if b.SideToMove == Black {
		h ^= sideKey
	}

	h ^= castleKeys[b.CastlingRights]

	if b.EnPassant != -1 {
		file := b.EnPassant % 8
		h ^= epKeys[file]
	}

	return h
}

func (b *Board) Perft(depth int) uint64 {
	if depth == 0 {
		return 1
	}

	var nodes uint64

	moves := b.GenerateLegalMoves()

	for _, move := range moves {

		before := b.Clone()

		undo := b.MakeMove(move)

		nodes += b.Perft(depth - 1)

		b.UnmakeMove(move, undo)

		if !before.Equals(b) {
			fmt.Printf("Corruption after: %s depth=%d\n", move, depth)

			fmt.Println("BEFORE:")
			fmt.Println(before)

			fmt.Println("AFTER:")
			fmt.Println(b)
		}
	}

	return nodes
}

func (b *Board) PerftDivide(depth int) uint64 {
	if depth == 0 {
		return 1
	}

	var total uint64

	moves := b.GenerateLegalMoves()

	for _, move := range moves {
		undo := b.MakeMove(move)

		nodes := b.Perft(depth - 1)

		b.UnmakeMove(move, undo)

		fmt.Printf("%s: %d\n", move.String(), nodes)

		total += nodes
	}

	return total
}

func (b *Board) Clone() *Board {
	c := *b
	return &c
}

func (b *Board) Equals(other *Board) bool {

	for i := range 12 {
		if b.Pieces[i] != other.Pieces[i] {
			fmt.Printf("Pieces[%d] mismatch\n", i)
			return false
		}
	}

	for i := range 3 {
		if b.Occupancy[i] != other.Occupancy[i] {
			fmt.Printf("Occupancy[%d] mismatch\n", i)
			return false
		}
	}

	if b.SideToMove != other.SideToMove {
		fmt.Println("SideToMove mismatch")
		return false
	}

	if b.CastlingRights != other.CastlingRights {
		fmt.Println("CastlingRights mismatch")
		return false
	}

	if b.EnPassant != other.EnPassant {
		fmt.Println("EnPassant mismatch")
		return false
	}

	if b.HalfmoveClock != other.HalfmoveClock {
		fmt.Println("HalfmoveClock mismatch")
		return false
	}

	if b.FullmoveNumber != other.FullmoveNumber {
		fmt.Println("FullmoveNumber mismatch")
		return false
	}

	if b.KingSq != other.KingSq {
		fmt.Println("KingSq mismatch")
		return false
	}

	if b.Hash != other.Hash {
		fmt.Println("Hash mismatch")
		return false
	}

	return true
}

/*  TODO:
// 	func (b *Board) ToFEN() string
*/

func (b *Board) String() string {
	var sb strings.Builder
	pieceChars := []rune{'P', 'N', 'B', 'R', 'Q', 'K', 'p', 'n', 'b', 'r', 'q', 'k'}

	sb.WriteString("  +-----------------+\n")
	for rank := 7; rank >= 0; rank-- {
		fmt.Fprintf(&sb, "%d | ", rank+1)
		for file := range 8 {
			square := rank*8 + file
			found := false
			for pieceTypeIndex := range 12 {
				if (b.Pieces[pieceTypeIndex]>>square)&1 != 0 {
					sb.WriteRune(pieceChars[pieceTypeIndex])
					found = true
					break
				}
			}
			if !found {
				sb.WriteRune('.')
			}
			sb.WriteRune(' ')
		}
		sb.WriteString("|\n")
	}
	sb.WriteString("  +-----------------+\n")
	sb.WriteString("    a b c d e f g h\n")

	fmt.Fprintf(&sb, "Side to move: %s\n", func() string {
		if b.SideToMove == White {
			return "White"
		} else {
			return "Black"
		}
	}())

	fmt.Fprintf(&sb, "Castling rights: %04b\n", b.CastlingRights)
	if b.EnPassant != -1 {
		fmt.Fprintf(&sb, "En Passant square: %s\n", squareToNotation(b.EnPassant))
	} else {
		sb.WriteString("En Passant square: -\n")
	}
	fmt.Fprintf(&sb, "Halfmove clock: %d\n", b.HalfmoveClock)
	fmt.Fprintf(&sb, "Fullmove number: %d\n", b.FullmoveNumber)
	fmt.Fprintf(&sb, "White King: %s, Black King: %s\n", squareToNotation(int(b.KingSq[White])), squareToNotation(int(b.KingSq[Black])))
	fmt.Fprintf(&sb, "White King: %d, Black King: %d\n", b.KingSq[White], b.KingSq[Black])

	return sb.String()
}
