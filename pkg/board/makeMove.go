package board

import "github.com/LucasSim0n/isabel/pkg/cmn"

func (b *Board) MakeMove(move cmn.Move) cmn.Undo {

	undo := cmn.Undo{
		CastlingRights: b.CastlingRights,
		EnPassant:      b.EnPassant,
		Halfmove:       b.HalfmoveClock,
		FullmoveNumber: b.FullmoveNumber,
		Hash:           b.Hash,
	}

	color := b.SideToMove
	enemy := cmn.GetOpposite(b.SideToMove)

	if b.EnPassant != -1 {
		b.Hash ^= cmn.EpKeys[b.EnPassant%8]
	}

	b.EnPassant = -1

	if move.Flags&cmn.FlagCapture != 0 {
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

	if move.Flags&cmn.FlagEnPassant != 0 {
		b.handleEnPassant(move, color)
		b.HalfmoveClock = 0
	}

	b.movePiece(move.Piece, color, move.From, move.To)

	if move.Flags&cmn.FlagCastling != 0 {
		b.handleCastling(move.To, color)
	}

	if move.Piece == cmn.King {
		b.KingSq[color] = uint8(move.To)

		if color == cmn.White {
			b.clearCastle(0b0011)
		} else {
			b.clearCastle(0b1100)
		}
	}

	if move.Piece == cmn.Rook {
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

	if move.Flags&cmn.FlagDoublePawnPush != 0 {
		if color == cmn.White {
			b.EnPassant = move.To - 8
		} else {
			b.EnPassant = move.To + 8
		}

		b.Hash ^= cmn.EpKeys[b.EnPassant%8]
	}

	if move.Piece == cmn.Pawn {
		b.HalfmoveClock = 0
	} else {
		b.HalfmoveClock++
	}

	if move.Flags&cmn.FlagPromotion != 0 {
		b.handlePromotion(move, color)
	}

	if color == cmn.Black {
		b.FullmoveNumber++
	}

	b.SideToMove = enemy
	b.Hash ^= cmn.SideKey

	return undo
}

func (b *Board) UnmakeMove(move cmn.Move, undo cmn.Undo) {
	b.HalfmoveClock = undo.Halfmove
	b.FullmoveNumber = undo.FullmoveNumber

	b.Hash ^= cmn.CastleKeys[b.CastlingRights]
	b.CastlingRights = undo.CastlingRights
	b.Hash ^= cmn.CastleKeys[b.CastlingRights]

	if b.EnPassant != -1 {
		b.Hash ^= cmn.EpKeys[b.EnPassant%8]
	}
	b.EnPassant = undo.EnPassant
	if b.EnPassant != -1 {
		b.Hash ^= cmn.EpKeys[b.EnPassant%8]
	}

	b.SideToMove = cmn.GetOpposite(b.SideToMove)
	b.Hash ^= cmn.SideKey

	color := b.SideToMove
	enemy := cmn.GetOpposite(color)

	if move.Flags&cmn.FlagPromotion != 0 {
		b.undoPromotion(move, color)
	}

	b.movePiece(move.Piece, color, move.To, move.From)

	if move.Flags&cmn.FlagEnPassant != 0 {
		b.undoEnPassant(move, enemy)
	}

	if move.Piece == cmn.King {
		b.KingSq[color] = uint8(move.From)
	}

	if move.Flags&cmn.FlagCapture != 0 {
		b.addPiece(move.Capture, enemy, move.To)
	}

	if move.Flags&cmn.FlagCastling != 0 {
		b.undoCastling(move.To, color)
	}
}

func (b *Board) movePiece(piece cmn.PieceType, color cmn.Color, from, to int) {
	b.removePiece(piece, color, from)
	b.addPiece(piece, color, to)
}

func (b *Board) addPiece(piece cmn.PieceType, color cmn.Color, sq int) {
	index := int(piece)
	if color == cmn.Black {
		index += 6
	}

	b.Pieces[index] |= (1 << sq)
	b.Occupancy[color] |= (1 << sq)
	b.Occupancy[2] |= (1 << sq)
	b.Hash ^= cmn.PieceKeys[index][sq]
}

func (b *Board) removePiece(piece cmn.PieceType, color cmn.Color, sq int) {
	index := int(piece)
	if color == cmn.Black {
		index += 6
	}

	b.Pieces[index] &^= (1 << sq)
	b.Occupancy[color] &^= (1 << sq)
	b.Occupancy[2] &^= (1 << sq)
	b.Hash ^= cmn.PieceKeys[index][sq]
}

/*** Special Moves ***/

func (b *Board) handleCastling(sq int, c cmn.Color) {
	switch sq {
	case 2:
		b.movePiece(cmn.Rook, c, 0, 3)
	case 6:
		b.movePiece(cmn.Rook, c, 7, 5)
	case 58:
		b.movePiece(cmn.Rook, c, 56, 59)
	case 62:
		b.movePiece(cmn.Rook, c, 63, 61)
	}
}

func (b *Board) undoCastling(sq int, c cmn.Color) {
	switch sq {
	case 2:
		b.movePiece(cmn.Rook, c, 3, 0)
	case 6:
		b.movePiece(cmn.Rook, c, 5, 7)
	case 58:
		b.movePiece(cmn.Rook, c, 59, 56)
	case 62:
		b.movePiece(cmn.Rook, c, 61, 63)
	}
}

func (b *Board) handlePromotion(move cmn.Move, color cmn.Color) {
	b.removePiece(cmn.Pawn, color, move.To)
	b.addPiece(move.Promotion, color, move.To)
}

func (b *Board) undoPromotion(move cmn.Move, color cmn.Color) {
	b.removePiece(move.Promotion, color, move.To)
	b.addPiece(cmn.Pawn, b.SideToMove, move.To)
}

func (b *Board) handleEnPassant(move cmn.Move, color cmn.Color) {
	var captureSq int

	if color == cmn.White {
		captureSq = move.To - 8
	} else {
		captureSq = move.To + 8
	}

	b.removePiece(cmn.Pawn, cmn.GetOpposite(color), captureSq)
}

func (b *Board) undoEnPassant(move cmn.Move, enemy cmn.Color) {
	var captureSq int

	if b.SideToMove == cmn.White {
		captureSq = move.To - 8
	} else {
		captureSq = move.To + 8
	}

	b.addPiece(cmn.Pawn, enemy, captureSq)
}

func (b *Board) clearCastle(mask uint8) {
	b.Hash ^= cmn.CastleKeys[b.CastlingRights]

	b.CastlingRights &= mask

	b.Hash ^= cmn.CastleKeys[b.CastlingRights]
}
