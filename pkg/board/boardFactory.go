package board

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"

	"github.com/LucasSim0n/isabel/pkg/cmn"
)

func NewFromFen(fen string) (*Board, error) {
	b := &Board{
		EnPassant:      -1,
		FullmoveNumber: 1,
	}

	parts := strings.Fields(fen)

	piecePlacement := parts[0]
	rank := 7
	file := 0

	for _, char := range piecePlacement {
		if char == '/' {
			rank--
			file = 0
			continue
		}

		if num, err := strconv.Atoi(string(char)); err == nil {
			file += num

		} else {
			pieceInfo, ok := cmn.PieceCharMap[char]
			if !ok {
				return nil, fmt.Errorf("Unknown FEN char: %s", fen)
			}

			square := rank*8 + file

			index := int(pieceInfo.Type)
			if pieceInfo.Color == cmn.Black {
				index += 6
			}

			b.Pieces[index] |= (1 << square)

			file++
		}
	}

	b.updateOccupancy()

	whiteKingBB := b.Pieces[cmn.King]
	b.KingSq[cmn.White] = uint8(bits.TrailingZeros64(uint64(whiteKingBB)))

	blackKingBB := b.Pieces[cmn.King+6]
	b.KingSq[cmn.Black] = uint8(bits.TrailingZeros64(uint64(blackKingBB)))

	b.SideToMove = cmn.White
	if parts[1] == "b" {
		b.SideToMove = cmn.Black
	}

	b.CastlingRights = 0
	for _, c := range parts[2] {
		switch c {
		case 'K':
			b.CastlingRights |= 0b1000
		case 'Q':
			b.CastlingRights |= 0b0100
		case 'k':
			b.CastlingRights |= 0b0010
		case 'q':
			b.CastlingRights |= 0b0001
		}
	}

	if parts[3] != "-" {
		s, err := cmn.NotationToSquare(parts[3])
		if err != nil {
			return nil, fmt.Errorf("Incorrect fen: %s", fen)
		}
		b.EnPassant = s
	}

	if len(parts) > 4 {
		halfmove, err := strconv.Atoi(parts[4])
		if err != nil {
			return nil, fmt.Errorf("Incorrect fen: %s", fen)
		}
		b.HalfmoveClock = halfmove

		if len(parts) > 5 {
			fullmove, err := strconv.Atoi(parts[5])
			if err != nil {
				return nil, fmt.Errorf("Incorrect fen: %s", fen)
			}
			b.FullmoveNumber = fullmove
		}
	}

	b.Hash = b.ComputeHash()

	return b, nil
}

func NewStartPosGame() *Board {
	return startPosGame.Clone()
}
